package kongclient

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/codemk8/apihub/pkg/k8s"
	"github.com/stretchr/testify/assert"
)

func createKongClient() *KongK8sClient {
	params := KongParams{}
	kong := NewKongK8sClient(params)
	return kong
}

func Test_parseNsAndSvc(t *testing.T) {
	ns, svc, ok := parseNsAndSvc("ns1:service")
	assert.True(t, ok)
	assert.Equal(t, "ns1", ns)
	assert.Equal(t, "service", svc)

	ns, svc, ok = parseNsAndSvc("service")
	assert.True(t, ok)
	assert.Equal(t, "default", ns)
	assert.Equal(t, "service", svc)

	ns, svc, ok = parseNsAndSvc(":service")
	assert.True(t, ok)
	assert.Equal(t, "", ns)
	assert.Equal(t, "service", svc)

	ns, svc, ok = parseNsAndSvc(":")
	assert.False(t, ok)

	ns, svc, ok = parseNsAndSvc("a:")
	assert.False(t, ok)
}

func TestGetKubernetesCluterIP(t *testing.T) {
	kong := createKongClient()
	ns, svc, ok := parseNsAndSvc("kubernetes")
	assert.True(t, ok)
	clusterIPs := k8s.GetServiceClusterIPPort(ns, svc, kong.K8sCs)
	assert.Equal(t, 1, len(clusterIPs))
	// default:kubernetes always has a 443 clusterIP
	assert.Equal(t, int32(443), clusterIPs[0])
}

func TestNewKongClient(t *testing.T) {
	kong := createKongClient()
	assert.Equal(t, "localhost", kong.KongSvcHost)
	total := kong.SmokeTestKong()
	assert.True(t, total >= 0)
	if total >= 0 {
		fmt.Printf("Found %d existing APIs in Kong\n", total)
	}
}

func TestRegisterServiceToKong(t *testing.T) {
	//kong := createKongClient()
	//kong.RegisterServiceToKong("default:http-echoserver", nil)
}

func TestPutAPIToKong(t *testing.T) {
	// This tests assumes we have a http-echoserver service running
	// in k8s already
	kong := createKongClient()
	name := "echo" // http-echoserver
	req := KongPutAPISpec{
		Name: name, //"http-echoserver",
		// we can also use a public upstream for testing, e.g: http://httpbin.org/
		UpstreamURL: "http://http-echoserver:80",
		URIs:        "/http-echoserver",
		StripURI:    true,
	}
	code, err := kong.PutNewAPI(&req)
	// or http.StatusOK(200)
	assert.Equal(t, http.StatusCreated, code)
	assert.Equal(t, nil, err)
	//time.Sleep(5 * time.Millisecond)
	code, err = kong.DeleteAPI(name)
	assert.Equal(t, http.StatusNoContent, code)
	assert.Equal(t, nil, err)
}

func TestMakePutParams(t *testing.T) {
	deployP := DeployParams{
		Uris:     "myapi/v1/",
		Force:    false,
		SvcRoot:  "",
		Name:     "abcd",
		StripURI: false,
	}
	service := "test_service"
	putParams := makePutParams(&deployP, service, 8080)
	assert.Equal(t, "abcd", putParams.Name)
	assert.Equal(t, "http://test_service:8080", putParams.UpstreamURL)
	assert.Equal(t, "/myapi/v1", putParams.URIs)

	// Use service as the Name if deployParam does not have an nonempty name
	deployP.Name = ""
	putParams = makePutParams(&deployP, service, 8080)
	assert.Equal(t, putParams.Name, service)

}
