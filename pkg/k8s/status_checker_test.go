package k8s

import (
	"github.com/mamezou-tech/concourse-k8s-resource/pkg/models"
	"github.com/stretchr/testify/assert"
	appv1 "k8s.io/api/apps/v1"
	"k8s.io/api/core/v1"
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
			UID:       "uid",
		},
		Spec: appv1.DeploymentSpec{
			Replicas: replicas(3),
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{},
			},
		},
	}
	app1rs := appv1.ReplicaSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "app1",
			Namespace: "test",
			OwnerReferences: []metav1.OwnerReference{{
				Kind:       "Deployment",
				Name:       "app1",
				UID:        "uid",
				Controller: control(true),
			}},
		},
		Spec: appv1.ReplicaSetSpec{
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{},
			},
		},
		Status: appv1.ReplicaSetStatus{
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
	clientset := fake.NewSimpleClientset(&app1, &app2, &app1rs)
	resources := []models.WatchResource{
		{Name: "app1", Kind: "deployment"},
		{Name: "app2", Kind: "statefulset"},
	}
	time.AfterFunc(1*time.Second, func() {
		t.Log("ready for pod...")
		app1rs.Status.ReadyReplicas = 3
		clientset.AppsV1().ReplicaSets("test").Update(&app1rs)
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
			UID:       "uid",
		},
		Spec: appv1.DeploymentSpec{
			Replicas: replicas(3),
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{},
			},
		},
	}
	app1rs := appv1.ReplicaSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "app1",
			Namespace: "test",
			OwnerReferences: []metav1.OwnerReference{{
				Kind:       "Deployment",
				Name:       "app1",
				UID:        "uid",
				Controller: control(true),
			}},
		},
		Spec: appv1.ReplicaSetSpec{
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{},
			},
		},
		Status: appv1.ReplicaSetStatus{
			ReadyReplicas: 1,
		},
	}
	clientset := fake.NewSimpleClientset(&app1, &app1rs)
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

func control(b bool) *bool {
	return &b
}
