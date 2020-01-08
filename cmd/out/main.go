package main

import (
	"encoding/json"
	"github.com/kudoh/concourse-k8s-resource/pkg/k8s"
	"github.com/kudoh/concourse-k8s-resource/pkg/k8s/kubectl"
	"github.com/kudoh/concourse-k8s-resource/pkg/models"
	"github.com/kudoh/concourse-k8s-resource/pkg/utils"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"log"
	"os"
)

var streams = genericclioptions.IOStreams{
	In:     os.Stdin,
	Out:    os.Stderr, // concourse console
	ErrOut: os.Stderr,
}

func main() {

	var request models.OutRequest
	if err := json.NewDecoder(os.Stdin).Decode(&request); err != nil {
		log.Fatalln("Illegal input format", err)
	}

	utils.Debug(&request.Source, "request: ", request)
	utils.ChangeWorkingDir()

	clientset, clientConfig := k8s.BuildClientSet(&request.Source)
	if request.Source.Namespace == "" {
		request.Source.Namespace = "default"
	}

	factory := kubectl.NewCommandFactory(&request.Params)
	commandConfig := &kubectl.CommandConfig{
		Clientset:    clientset,
		ClientConfig: clientConfig,
		Streams:      streams,
		Namespace:    request.Source.Namespace,
		Resources:    request.Source.WatchResources,
		Params:       &request.Params,
	}
	if err := kubectl.RunCommand(factory, commandConfig); err != nil {
		log.Fatalln("cannot run kubectl command", err)
	}
	if !request.Params.Delete {
		log.Println("check status for", request.Source.WatchResources)
		if ok := k8s.CheckResourceStatus(clientset, request.Source.Namespace, request.Source.WatchResources, request.Params.StatusCheckTimeout); !ok {
			log.Fatalln("resource is not running...")
		}
	}

	response := createResponse(request, clientset)

	utils.Debug(&request.Source, "response: ", *response)
	if err := json.NewEncoder(os.Stdout).Encode(response); err != nil {
		log.Fatalln("Output Failure", err)
	}
}

func createResponse(request models.OutRequest, clientset *kubernetes.Clientset) *models.OutResponse {

	if request.Params.Delete {
		// resources is deleted, so just return empty response
		return &models.OutResponse{
			Version:  models.Version{},
			Metadata: nil,
		}
	}

	// apply or undo
	version, err := k8s.GetCurrentVersion(&request.Source, clientset)
	if err != nil {
		log.Fatalln(err)
	}
	metadatas, err := k8s.GenerateMetadatas(&request.Source, clientset)
	if err != nil {
		log.Fatalln(err)
	}

	response := models.OutResponse{
		Version:  *version,
		Metadata: metadatas,
	}
	return &response
}
