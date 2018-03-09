package kongclient

import (
	"crypto/tls"
	"fmt"
	"strings"

	"github.com/codemk8/apihub/pkg/k8s"
	"github.com/go-resty/resty"
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

// https://${external_url}:${kong_admin_port}/apis
func (kc *KongK8sClient) constructKongAPIUrl() string {
	return "https://" + kc.KongSvcHost + ":" + kc.KongSvcPort + "/apis"
}

// RegisterServiceToKong adds an API endpoint to kong API
func (kc *KongK8sClient) RegisterServiceToKong(svc string) bool {
	ns, svcName, ok := parseNsAndSvc(svc)
	if !ok {
		fmt.Printf("Error parsing service name :%s\n", svc)
		return ok
	}
	clusterIPs := k8s.GetServiceClusterIPPort(ns, svcName, kc.K8sCs)
	if len(clusterIPs) == 0 {
		fmt.Printf("Could not find any clusterIPs for service %s", svc)
		return false
	}
	// Only register the first port now

	return true
}

func initResty() {
	resty.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
}

// Deploy implements the "deploy" command
func Deploy(args []string, deployParams *DeployParams) bool {
	// TODO set params from cached values
	// use default now
	params := KongParams{}
	kong := NewKongK8sClient(params)

	if kong == nil {
		fmt.Println("Error init API gateway client, check if kong is a valid service in k8s.")
		return false
	}
	initResty()

	for _, service := range args {
		kong.RegisterServiceToKong(service)
	}
	return true
}
