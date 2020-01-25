package k8s

import (
	"github.com/mamezou-tech/concourse-k8s-resource/pkg/models"
	"k8s.io/client-go/kubernetes"
	"log"
	"strconv"
)

func GenerateMetadatas(source *models.Source, clientset kubernetes.Interface) ([]models.Metadata, error) {

	var metas []MetadataReader
	for _, resource := range source.WatchResources {
		switch {
		case IsDeployment(resource.Kind):
			reader, err := NewDeploymentReader(clientset, source.Namespace, resource.Name)
			if err != nil {
				return nil, err
			}
			metas = append(metas, reader)
		case IsStatefulSet(resource.Kind):
			reader, err := NewStatefulSetReader(clientset, source.Namespace, resource.Name)
			if err != nil {
				return nil, err
			}
			metas = append(metas, reader)
		default:
			log.Fatalln("unsupported resource kind", resource.Kind)
		}
	}
	return toMetadatas(metas), nil
}

func toMetadatas(readers []MetadataReader) []models.Metadata {
	var resp []models.Metadata
	for _, reader := range readers {
		meta := reader.GetObjectMeta()
		revision, err := reader.GetRevision()
		if err != nil {
			revision = -1
		}
		converted := []models.Metadata{
			{Name: "name", Value: meta.Name},
			{Name: "namespace", Value: meta.Namespace},
			{Name: "resource_version", Value: meta.ResourceVersion},
			{Name: "uid", Value: string(meta.UID)},
			{Name: "creation_timestamp", Value: meta.CreationTimestamp.String()},
			{Name: "revision", Value: strconv.FormatInt(revision, 10)},
		}
		resp = append(resp, converted...)
	}
	return resp
}
