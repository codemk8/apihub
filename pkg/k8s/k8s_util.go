package k8s

import (
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
	pvNamespace = "default"
	pvName      = "apihub-infra"
)

const (
	pvCapacity = "2Gi"
)

// MakeHostDirForPv makes a directory for hostPath PV
// For single-node cluster only
func MakeHostDirForPv(pvName string) (string, error) {
	home, _ := homedir.Dir()
	path := filepath.Join(home, ".apihub", "pvs")
	err := os.MkdirAll(path, 0755)
	return path, err
}

// NewK8sClient create the K8sClient
func NewK8sClient() (*kubernetes.Clientset, error) {
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

	// create the clientset
	cs, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Print(err.Error())
	}
	return cs, err
}

// CheckK8s does basic checkings
func CheckK8s() error {
	k8sClient, err := NewK8sClient()
	if err != nil {
		log.Print(err.Error())
		return err
	}
	pods, err := k8sClient.CoreV1().Pods("").List(metav1.ListOptions{})
	if err != nil {
		log.Print(err.Error())
		return err
	}
	_, err = k8sClient.CoreV1().Services(tillerNs).Get(tillerSvcName, metav1.GetOptions{})
	if err != nil {
		log.Printf("Could not find service %s, check helm installation, error %v\n", tillerSvcName, err)
		return err
	}

	log.Printf("Found %d pods in the cluster, checking OK.\n", len(pods.Items))
	return nil
}

// AddPV adds necessary persistent volumes for apihub
func AddPV() error {
	k8sClient, err := NewK8sClient()
	if err != nil {
		return err
	}

	pv, err := k8sClient.CoreV1().PersistentVolumes().Get(pvName, metav1.GetOptions{})
	if err != nil {
		log.Printf("Error getting PV %v", err)
	}
	// if found, no need to create the PV
	if pv != nil {
		return nil
	}
	storageQuantity, _ := resource.ParseQuantity(pvCapacity)
	hostPath, err := MakeHostDirForPv(pvName)
	if err != nil {
		log.Printf("Error creating directory %s for PV, error %v", hostPath, err)
		return err
	}

	k8sexpVolume := &api.PersistentVolume{
		TypeMeta: metav1.TypeMeta{Kind: "PersistentVolume",
			APIVersion: "v1"},
		ObjectMeta: metav1.ObjectMeta{
			Name:   pvName,
			Labels: map[string]string{"tag": "apihub"},
		},
		Spec: api.PersistentVolumeSpec{
			Capacity: api.ResourceList{
				api.ResourceStorage: storageQuantity,
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
	_, err = k8sClient.CoreV1().PersistentVolumes().Create(k8sexpVolume)
	if err != nil {
		log.Print(err.Error())
	}
	return nil
}
