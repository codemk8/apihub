package kongclient

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/codemk8/apihub/pkg/k8s"
	"github.com/go-resty/resty"
)

func initResty() {
	resty.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
}

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
func (kc *KongK8sClient) makeKongAPIURL(path string) string {
	return "https://" + kc.KongSvcHost + ":" + kc.KongSvcPort + path
}

// RegisterServiceToKong adds an API endpoint to kong API
func (kc *KongK8sClient) RegisterServiceToKong(svc string, deployParams *DeployParams) bool {
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
	if len(clusterIPs) > 1 {
		fmt.Printf("Found multiple clusterIPs for service %s, trying to add the first port %d.", svc, clusterIPs[0])
	}
	// Only register the first port now
	putSpec := makePutParams(deployParams, svc, clusterIPs[0])
	code, err := kc.PutNewAPI(putSpec)
	if err != nil {
		log.Printf("Error put new API %v", err)
	}
	if code == http.StatusOK || code == http.StatusCreated {
		return true
	}
	log.Printf("Put API receive error code %d", code)
	return false
}

// PutNewAPI tries to add a new API to kong
func (kc *KongK8sClient) PutNewAPI(req *KongPutAPISpec) (int, error) {
	body, _ := json.Marshal(req)
	response, err := resty.R().SetHeader("Content-Type", "application/json").
		SetBody(body).Put(kc.makeKongAPIURL("/apis/"))
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
	return response.StatusCode(), err
}

// DeleteAPI tries to add a new API to kong
func (kc *KongK8sClient) DeleteAPI(name string) (int, error) {
	response, err := resty.R().Delete(kc.makeKongAPIURL("/apis/" + name))
	// Expecting http.StatsNoContent (204)
	return response.StatusCode(), err
}

// SmokeTestKong calls a simple API on kong admin
func (kc *KongK8sClient) SmokeTestKong() (int, error) {
	response, err := resty.R().Get(kc.makeKongAPIURL("/"))
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
func Deploy(service string, deployParams *DeployParams) bool {
	// TODO set params from cached values
	// use default now
	params := KongParams{}
	kong := NewKongK8sClient(params)

	if kong == nil {
		fmt.Println("Error init API gateway client, check if kong is a valid service in k8s.")
		return false
	}
	return kong.RegisterServiceToKong(service, deployParams)
}

func makeUpstreamURL(serverName string, port int32, subpath string) string {
	// Hardcode to http for now
	if subpath == "" {
		return "http://" + serverName + ":" + strconv.Itoa(int(port))
	}
	if subpath[0] != '/' {
		return "http://" + serverName + ":" + strconv.Itoa(int(port)) + "/" + subpath
	}
	return "http://" + serverName + ":" + strconv.Itoa(int(port)) + subpath
}

// make sure string are all ascii, which is important for uri
func isASCII(s string) bool {
	for _, c := range s {
		if c > 127 {
			return false
		}
	}
	return true
}

func makePutParams(deploy *DeployParams, service string, port int32) *KongPutAPISpec {
	putSpec := &KongPutAPISpec{
		Name:        deploy.Name,
		StripURI:    deploy.StripURI,
		UpstreamURL: makeUpstreamURL(service, port, deploy.SvcRoot),
		URIs:        deploy.Uris,
		Retries:     5,
	}
	if deploy.Name == "" {
		putSpec.Name = service
	}
	// correct bad UIR input
	// e.g. "myapi/v1/" to "myapi/v1"
	if putSpec.URIs[0] != '/' {
		putSpec.URIs = "/" + putSpec.URIs
	}
	if putSpec.URIs[len(putSpec.URIs)-1] == '/' {
		putSpec.URIs = putSpec.URIs[0 : len(putSpec.URIs)-1]
	}

	return putSpec
}
