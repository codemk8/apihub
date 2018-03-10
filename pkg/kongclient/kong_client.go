package kongclient

import (
	"crypto/tls"
	"encoding/json"
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
func (kc *KongK8sClient) constructKongAPIUrl(path string) string {
	return "https://" + kc.KongSvcHost + ":" + kc.KongSvcPort + "/apis" + path
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

// PostNewAPI tries to add a new API to kong
func (kc *KongK8sClient) PutNewAPI() (int, error) {
	req := KongPutAPISpec{
		Name:        "http-echoserver",
		UpstreamURL: "http://http-echoserver:80",
		URIs:        "/http-echoserver",
		StripURI:    false,
	}
	body, _ := json.Marshal(&req)
	response, err := resty.R().SetHeader("Content-Type", "application/json").
		SetBody(body).Put(kc.constructKongAPIUrl("/"))
	if err != nil {
		fmt.Printf("Error Put API to kong: %+v", err)
		return 0, err
	}
	APISpec := KongAPISpec{}
	err = json.Unmarshal(response.Body(), &APISpec)
	if err != nil {
		fmt.Printf("Error unmarshalling response: %v", err)
		return 0, err
	}
	// Expecting http.StatusOK(200) or http.StatusCreated(201) code
	// Or http.StatusConflict(409) if there is conflict
	fmt.Printf("Return code %d\n", response.StatusCode())
	fmt.Printf("Response: %+v\n", APISpec)
	return response.StatusCode(), err
}

// SmokeTestKong calls a simple API on kong admin
func (kc *KongK8sClient) SmokeTestKong() (int, error) {
	response, err := resty.R().Get(kc.constructKongAPIUrl("/"))
	if err != nil {
		fmt.Printf("Error calling GET on Kong admin: %v", err)
		return 0, err
	}
	APIResult := KongGetResp{}
	err = json.Unmarshal(response.Body(), &APIResult)
	if err != nil {
		fmt.Printf("Error unmarshalling response: %v", err)
		return 0, err
	}
	return APIResult.Total, nil
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
	for _, service := range args {
		kong.RegisterServiceToKong(service)
	}
	return true
}
