package k8s

import (
	"fmt"
	"github.com/mamezou-tech/concourse-k8s-resource/pkg/models"
	"k8s.io/client-go/kubernetes"
	"log"
	"strings"
)

func GetCurrentVersion(source *models.Source, clientset kubernetes.Interface) (*models.Version, error) {
	builder := newResourceVersionBuilder()
	for _, resource := range source.WatchResources {
		switch {
		case IsDeployment(resource.Kind):
			r, err := NewDeploymentReader(clientset, source.Namespace, resource.Name)
			if err != nil {
				return nil, err
			}
			builder.addReader(r)
		case IsStatefulSet(resource.Kind):
			r, err := NewStatefulSetReader(clientset, source.Namespace, resource.Name)
			if err != nil {
				return nil, err
			}
			builder.addReader(r)
		default:
			log.Fatalln("unsupported resource kind", resource.Kind)
		}
	}
	return builder.build(), nil
}

type resourceVersionBuilder struct {
	metadataReaders []MetadataReader
}

func newResourceVersionBuilder() *resourceVersionBuilder {
	return &resourceVersionBuilder{}
}

func (builder *resourceVersionBuilder) addReader(reader MetadataReader) {
	builder.metadataReaders = append(builder.metadataReaders, reader)
}

func (builder *resourceVersionBuilder) build() *models.Version {
	var versions []string
	for _, reader := range builder.metadataReaders {
		objectMeta := reader.GetObjectMeta()
		rev, err := reader.GetRevision()
		if err != nil {
			rev = -1
		}
		versions = append(versions, fmt.Sprintf("%s:%s:%d", objectMeta.Namespace,
			objectMeta.Name, rev))
	}
	return &models.Version{Revision: strings.Join(versions, "+")}
}
