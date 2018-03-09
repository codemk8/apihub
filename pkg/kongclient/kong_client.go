package kongclient

import (
	"fmt"
	"strings"

	"github.com/codemk8/apihub/pkg/k8s"
)

// parse "default:service_name" into pair of (default(namespace), service_name)
func parseNsAndSvc(svc string) (namespace string, service string, ok bool) {
	if svc == "" {
		return "", "", false
	}
	if strings.Contains(svc, ":") {
		pairs := strings.Split(svc, ":")
		if pairs[1] == "" {
			return "", "", false
		}
		return pairs[0], pairs[1], true
	}
	return "default", svc, true
}

// RegisterServiceToKong adds an API endpoint to kong API
func (kc *KongK8sClient) RegisterServiceToKong(svc string) bool {
	// TODO
	clusterIPs := k8s.GetServiceClusterIPPort("default", svc, kc.K8sCs)
	fmt.Println(clusterIPs)
	return true
}

// Deploy implements the "deploy" command
func Deploy(args []string) {
	// TODO set params from cached values
	// use default now
	params := KongParams{}
	kong := NewKongK8sClient(params)
	if kong == nil {
		fmt.Println("Error init API gateway client, check if kong is a valid service in k8s.")
		return
	}
	for _, service := range args {
		kong.RegisterServiceToKong(service)
	}
}
