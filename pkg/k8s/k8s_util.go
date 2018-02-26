package k8s

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
	api "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	tillerNs      = "kube-system"
	tillerSvcName = "tiller-deploy"
	// persistent volume
	namespace = "default"
	pvName    = "apihub-infra"
)

const (
	pvCapacity = "1Gi"
)

// MakeHostDirForPv makes a directory for hostPath PV
// For single-node cluster only
func MakeHostDirForPv(pvName string) (string, error) {
	home, _ := homedir.Dir()
	path := filepath.Join(home, ".apihub", "pvs")
	err := os.MkdirAll(path, 0755)
	return path, err
}

// K8sClient is reused in every session
type K8sClient struct {
	clientset *kubernetes.Clientset
}

// NewK8sClient create the K8sClient
func NewK8sClient() (*K8sClient, error) {
	home, err := homedir.Dir()
	if err != nil {
		log.Print(err.Error())
		return nil, err
	}
	kubeconfig := filepath.Join(home, ".kube", "config")

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Print(err.Error())
		return nil, err
	}

	var ks K8sClient
	// create the clientset
	cs, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Print(err.Error())
		return nil, err
	}
	ks.clientset = cs
	return &ks, nil
}

// CheckK8s does basic checkings
func CheckK8s() error {
	k8sClient, err := NewK8sClient()
	if err != nil {
		log.Print(err.Error())
		return err
	}
	pods, err := k8sClient.clientset.CoreV1().Pods("").List(metav1.ListOptions{})
	if err != nil {
		log.Print(err.Error())
		return err
	}
	_, err = k8sClient.clientset.CoreV1().Services("").Get(tillerSvcName, metav1.GetOptions{})
	if err != nil {
		log.Printf("Could not find service %s, check helm installation.\n", tillerSvcName)
		log.Print(err.Error())
		return err
	}

	log.Printf("Found %d pods in the cluster, checking OK.\n", len(pods.Items))
	return nil
}

// func getNamespace(ns: string) {
//
//}
//
func AddPV() error {
	k8sClient, err := NewK8sClient()
	if err != nil {
		return err
	}
	pvs, err := k8sClient.clientset.CoreV1().PersistentVolumes().List(metav1.ListOptions{})
	storageQuantity1Gi, _ := resource.ParseQuantity("1Gi")
	hostPath, err := MakeHostDirForPv(pvName)
	if err != nil {
		log.Printf("Error creating directory %s for PV, error %v", hostPath, err)
		return err
	}
	//	namespace, err := k8sClient.clientset.CoreV1().Namespaces().Get(namespace, metav1.GetOptions{})
	k8sexpVolume := &api.PersistentVolume{
		TypeMeta: metav1.TypeMeta{Kind: "PersistentVolume",
			APIVersion: "v1"},
		ObjectMeta: metav1.ObjectMeta{
			Name:   pvName,
			Labels: map[string]string{"type": "local"},
		},
		Spec: api.PersistentVolumeSpec{
			Capacity: api.ResourceList{
				api.ResourceStorage: storageQuantity1Gi,
			},
			PersistentVolumeSource: api.PersistentVolumeSource{
				HostPath: &api.HostPathVolumeSource{
					Path: hostPath,
				},
			},
			AccessModes: []api.PersistentVolumeAccessMode{
				api.PersistentVolumeAccessMode("ReadWriteOnce"),
			},
			PersistentVolumeReclaimPolicy: api.PersistentVolumeReclaimRecycle,
		},
		Status: api.PersistentVolumeStatus{},
	}
	_, err = k8sClient.clientset.CoreV1().PersistentVolumes().Create(k8sexpVolume)
	if err != nil {
		log.Print(err.Error())
	}
	log.Printf("Found %d pvs", len(pvs.Items))
	for _, pv := range pvs.Items {
		for key, label := range pv.Labels {
			fmt.Printf("%s: %s\n", key, label)
		}
		print(pv.Spec.HostPath.Path)
	}
	return nil
}
