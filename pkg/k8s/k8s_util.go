package k8s

import (
	"log"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	tillerSvcName = "tiller-deploy"
)

func CheckK8s() bool {
	home, err := homedir.Dir()
	if err != nil {
		log.Print(err.Error())
		return false
	}
	kubeconfig := filepath.Join(home, ".kube", "config")

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Print(err.Error())
		return false
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Print(err.Error())
		return false
	}
	pods, err := clientset.CoreV1().Pods("").List(metav1.ListOptions{})
	if err != nil {
		log.Print(err.Error())
		return false
	}
	svcs, err := clientset.CoreV1().Services("").List(metav1.ListOptions{})
	if err != nil {
		log.Print(err.Error())
		return false
	}

	log.Printf("Found %d pods, %d services in the cluster, checking OK.\n", len(pods.Items), len(svcs.Items))
	for _, svc := range svcs.Items {
		if svc.Name == tillerSvcName {
			return true
		}
	}
	log.Printf("Could not find service %s, check helm installation.\n", tillerSvcName)
	return false
}
