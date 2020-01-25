package k8s

import (
	"fmt"
	"github.com/mamezou-tech/concourse-k8s-resource/pkg/models"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"log"
	"os"
	"os/signal"
	"time"
)

type statusChecker struct {
	complete  chan error
	timeout   <-chan time.Time
	interrupt chan os.Signal
	clientset kubernetes.Interface
	resource  *models.WatchResource
	namespace string
}

func CheckResourceStatus(clientset kubernetes.Interface, namespace string, resources []models.WatchResource, waitseconds int32) bool {

	var timeout time.Duration
	if waitseconds == 0 {
		timeout = 5 * time.Minute
	} else {
		timeout = time.Duration(waitseconds) * time.Second
	}

	var checkers []statusChecker
	for idx := range resources {
		checker := statusChecker{
			complete:  make(chan error),
			timeout:   time.After(timeout),
			interrupt: make(chan os.Signal, 1),
			clientset: clientset,
			resource:  &resources[idx],
			namespace: namespace,
		}
		checkers = append(checkers, checker)
		signal.Notify(checker.interrupt, os.Interrupt)
		// run check parallel
		go func() {
			checker.complete <- checker.check()
		}()
	}

	for _, c := range checkers {
		// block until completed
		select {
		case err := <-c.complete:
			if err != nil {
				log.Println("status check", "error!", "->", c.resource.Name, err)
				return false
			} else {
				log.Println("status check", "ok!", "->", c.resource.Name)
			}
		case <-c.timeout:
			log.Println("status check", "timeout!", "->", c.resource.Name)
			return false
		}
	}
	return true
}

func (c *statusChecker) check() error {

	for i := 0; ; i++ {
		select {
		case <-c.interrupt:
			signal.Stop(c.interrupt)
			return fmt.Errorf("interrupted")
		default:
		}
		switch {
		case IsDeployment(c.resource.Kind):
			d, err := c.clientset.AppsV1().Deployments(c.namespace).Get(c.resource.Name, metav1.GetOptions{})
			if err != nil {
				return err
			}
			if d.Status.ReadyReplicas == *d.Spec.Replicas {
				return nil
			}
		case IsStatefulSet(c.resource.Kind):
			sts, err := c.clientset.AppsV1().StatefulSets(c.namespace).Get(c.resource.Name, metav1.GetOptions{})
			if err != nil {
				return err
			}
			if sts.Status.ReadyReplicas == *sts.Spec.Replicas {
				return nil
			}
		default:
			log.Fatalln("unsupported resource kind", c.resource.Kind)
		}
		if i%10 == 0 {
			log.Println("waiting", c.resource.Name)
		}
		time.Sleep(1 * time.Second)
	}
}
