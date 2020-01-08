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
	var request models.InRequest
	if err := json.NewDecoder(os.Stdin).Decode(&request); err != nil {
		log.Fatalln("Illegal input format", err)
	}
	utils.Debug(&request.Source, "request: ", request)
	utils.ChangeWorkingDir()

	clientset, _ := k8s.BuildClientSet(&request.Source)
	response := createResponse(request, clientset)
	utils.WriteFile("version", response.Version.Revision)

	utils.Debug(&request.Source, "response: ", *response)
	if err := json.NewEncoder(os.Stdout).Encode(response); err != nil {
		log.Fatalln("Output Failure", err)
	}
}

func createResponse(request models.InRequest, clientset *kubernetes.Clientset) *models.InResponse {

	metadatas, err := k8s.GenerateMetadatas(&request.Source, clientset)
	if errors.IsNotFound(err) {
		log.Println("not found", err)
		return &models.InResponse{Version: request.Version}
	} else if err != nil {
		panic(err.Error())
	}

	response := models.InResponse{Version: request.Version, Metadata: metadatas}
	return &response
}
