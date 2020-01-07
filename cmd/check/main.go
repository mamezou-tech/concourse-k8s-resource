package main

import (
	"encoding/json"
	"github.com/kudoh/concourse-k8s-resource/pkg/k8s"
	"github.com/kudoh/concourse-k8s-resource/pkg/models"
	"github.com/kudoh/concourse-k8s-resource/pkg/utils"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"log"
	"os"
)

func main() {
	var request models.CheckRequest
	if err := json.NewDecoder(os.Stdin).Decode(&request); err != nil {
		log.Fatalln("Illegal input format:", err)
	}
	utils.Debug(&request.Source, "request: ", request)

	clientset, _ := k8s.BuildClientSet(&request.Source)
	response := createCheckResponse(&request, clientset)

	utils.Debug(&request.Source, "response:", response)

	if err := json.NewEncoder(os.Stdout).Encode(response); err != nil {
		log.Fatalln("Output Failure:", err)
	}
}

func createCheckResponse(request *models.CheckRequest, clientset *kubernetes.Clientset) []models.Version {
	newVersion, err := k8s.GetCurrentVersion(&request.Source, clientset)
	if errors.IsNotFound(err) {
		// resource not found
		log.Printf("%v not found", request.Source.WatchResources)
		return []models.Version{request.Version}
	} else if err != nil {
		panic(err.Error())
	}

	var response []models.Version
	if request.Version.Revision != "" {
		response = append(response, request.Version)
		if request.Version.Revision != newVersion.Revision {
			response = append(response, *newVersion)
		}
	} else {
		response = append(response, *newVersion)
	}
	return response
}
