package k8s

import (
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/kubectl/pkg/util/deployment"
)

type MetadataReader interface {
	GetObjectMeta() *metav1.ObjectMeta
	GetRevision() (int64, error)
}

type DeploymentReader struct {
	resource *appsv1.Deployment
}

var _ MetadataReader = &DeploymentReader{}

func NewDeploymentReader(clientset *kubernetes.Clientset, namespace string, name string) (*DeploymentReader, error) {
	d, err := clientset.AppsV1().Deployments(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return &DeploymentReader{resource: d}, nil
}

func (r *DeploymentReader) GetRevision() (int64, error) {
	return deployment.Revision(r.resource)
}

func (r *DeploymentReader) GetObjectMeta() *metav1.ObjectMeta {
	return &r.resource.ObjectMeta
}

type StatefulSetReader struct {
	resource *appsv1.StatefulSet
	revision *appsv1.ControllerRevision
}

var _ MetadataReader = &StatefulSetReader{}

func NewStatefulSetReader(clientset *kubernetes.Clientset, namespace string, name string) (*StatefulSetReader, error) {
	sts, err := clientset.AppsV1().StatefulSets(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	rev, err := clientset.AppsV1().ControllerRevisions(namespace).Get(sts.Status.CurrentRevision, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return &StatefulSetReader{resource: sts, revision: rev}, nil
}

func (r *StatefulSetReader) GetRevision() (int64, error) {
	return r.revision.Revision, nil
}

func (r *StatefulSetReader) GetObjectMeta() *metav1.ObjectMeta {
	return &r.resource.ObjectMeta
}
