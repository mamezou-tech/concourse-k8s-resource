package k8s

import (
	"github.com/mamezou-tech/concourse-k8s-resource/pkg/models"
	"github.com/stretchr/testify/assert"
	appv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	"testing"
	"time"
)

func TestCheckResourceStatus(t *testing.T) {
	assert := assert.New(t)

	app1 := appv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "test",
			Name:      "app1",
		},
		Spec: appv1.DeploymentSpec{
			Replicas: replicas(3),
		},
		Status: appv1.DeploymentStatus{
			ReadyReplicas: 1,
		},
	}
	app2 := appv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "test",
			Name:      "app2",
		},
		Spec: appv1.StatefulSetSpec{
			Replicas: replicas(2),
		},
		Status: appv1.StatefulSetStatus{
			ReadyReplicas: 0,
		},
	}
	clientset := fake.NewSimpleClientset(&app1, &app2)
	resources := []models.WatchResource{
		{Name: "app1", Kind: "deployment"},
		{Name: "app2", Kind: "statefulset"},
	}
	time.AfterFunc(1*time.Second, func() {
		t.Log("ready for pod...")
		app1.Status.ReadyReplicas = 3
		clientset.AppsV1().Deployments("test").Update(&app1)
		app2.Status.ReadyReplicas = 2
		clientset.AppsV1().StatefulSets("test").Update(&app2)
	})

	ok := CheckResourceStatus(clientset, "test", resources, 5)

	assert.True(ok)
}

func TestTimeout(t *testing.T) {
	app1 := appv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "test",
			Name:      "app1",
		},
		Spec: appv1.DeploymentSpec{
			Replicas: replicas(1),
		},
		Status: appv1.DeploymentStatus{
			ReadyReplicas: 0, // not updated
		},
	}
	clientset := fake.NewSimpleClientset(&app1)
	resources := []models.WatchResource{
		{Name: "app1", Kind: "deployment"},
	}

	ok := CheckResourceStatus(clientset, "test", resources, 2)

	assert.False(t, ok)
}

func replicas(num int) *int32 {
	n := int32(num)
	return &n
}
