package k8s

import (
	"github.com/kudoh/concourse-k8s-resource/pkg/models"
	"github.com/stretchr/testify/assert"
	appv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	"testing"
	"time"
)

func TestGenerateMetadatas(t *testing.T) {
	assert := assert.New(t)

	app1 := appv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Namespace:         "test",
			Name:              "app1",
			UID:               "uid1",
			ResourceVersion:   "1",
			CreationTimestamp: metav1.Date(2020, 1, 2, 12, 13, 14, 1000000000, time.UTC),
			Annotations: map[string]string{
				"deployment.kubernetes.io/revision": "100",
			},
		},
	}
	app2 := appv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Namespace:         "test",
			Name:              "app2",
			UID:               "uid2",
			ResourceVersion:   "2",
			CreationTimestamp: metav1.Date(2020, 1, 3, 13, 14, 15, 1000000000, time.UTC),
		},
		Status: appv1.StatefulSetStatus{
			CurrentRevision: "app2-rev",
		},
	}
	rev := appv1.ControllerRevision{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "test",
			Name:      "app2-rev",
		},
		Revision: 200,
	}

	clientset := fake.NewSimpleClientset(&app1, &app2, &rev)
	clientset.AppsV1()
	source := models.Source{
		Namespace: "test",
		WatchResources: []models.WatchResource{
			{Kind: "Deployment", Name: "app1"},
			{Kind: "StatefulSet", Name: "app2"},
		},
	}

	metadatas, err := GenerateMetadatas(&source, clientset)
	if err != nil {
		assert.FailNow("error", err)
	}
	expected := []models.Metadata{
		{Name: "name", Value: "app1"},
		{Name: "namespace", Value: "test"},
		{Name: "resource_version", Value: "1"},
		{Name: "uid", Value: "uid1"},
		{Name: "creation_timestamp", Value: "2020-01-02 12:13:15 +0000 UTC"},
		{Name: "revision", Value: "100"},
		{Name: "name", Value: "app2"},
		{Name: "namespace", Value: "test"},
		{Name: "resource_version", Value: "2"},
		{Name: "uid", Value: "uid2"},
		{Name: "creation_timestamp", Value: "2020-01-03 13:14:16 +0000 UTC"},
		{Name: "revision", Value: "200"},
	}
	assert.EqualValues(expected, metadatas)
}
