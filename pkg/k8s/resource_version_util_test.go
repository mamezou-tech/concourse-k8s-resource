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

func TestGetCurrentVersion(t *testing.T) {
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
	source := models.Source{
		Namespace: "test",
		WatchResources: []models.WatchResource{
			{Name: "app1", Kind: "deploy"},
			{Name: "app2", Kind: "sts"},
		},
	}
	version, err := GetCurrentVersion(&source, clientset)
	if err != nil {
		assert.FailNow("error", err)
	}
	assert.EqualValues(version, &models.Version{Revision: "test:app1:100+test:app2:200"})
}
