package k8s

import (
	"log"
	"testing"
)

func TestGetServiceNodePort(t *testing.T) {
	cs, _ := NewK8sClient()
	log.Printf("Node port %v", GetServiceNodePort("default", "apihub-kong-kong-admin", cs))
	log.Printf("cluster IP port %v", GetServiceClusterIPPort("default", "apihub-kong-kong-admin", cs))

	log.Printf("Node port %v", GetServiceNodePort("default", "http-echoserver", cs))
	log.Printf("cluster IP port %v", GetServiceClusterIPPort("default", "http-echoserver", cs))
}
