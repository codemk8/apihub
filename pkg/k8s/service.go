package k8s

import (
	"log"

	"k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// GetServiceSpec gets the spec of a service in k8s
func GetServiceSpec(namespace string, svcName string, clientSet *kubernetes.Clientset) *v1.ServiceSpec {
	svc, err := clientSet.CoreV1().Services(namespace).Get(svcName, meta_v1.GetOptions{})
	if err != nil {
		log.Printf("Error getting service %s:%s: %v", namespace, svcName, err)
		return nil
	}
	if svc == nil {
		log.Printf("Could not find the service %s:%s", namespace, svcName)
		return nil
	}
	return &svc.Spec
}

// GetServiceNodePort gets the node port number of a service in k8s
func GetServiceNodePort(namespace string, svcName string, clientSet *kubernetes.Clientset) []int32 {
	var ports []int32
	spec := GetServiceSpec(namespace, svcName, clientSet)
	if spec != nil {
		if spec.Type == "NodePort" {
			for i := range spec.Ports {
				ports = append(ports, spec.Ports[i].NodePort)
			}
		}
	}
	return ports
}

// GetServiceClusterIPPort gets the node port number of a service in k8s
func GetServiceClusterIPPort(namespace string, svcName string, clientSet *kubernetes.Clientset) []int32 {
	var ports []int32
	spec := GetServiceSpec(namespace, svcName, clientSet)
	if spec != nil {
		if spec.Type == "ClusterIP" {
			for i := range spec.Ports {
				ports = append(ports, spec.Ports[i].Port)
			}
		}
	}
	return ports
}
