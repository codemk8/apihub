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
	total, err := kong.SmokeTestKong()
	fmt.Printf("Found %d existing APIs in Kong\n", total)
	assert.Equal(t, nil, err)
}

func TestRegisterServiceToKong(t *testing.T) {
	kong := createKongClient()
	kong.RegisterServiceToKong("default:http-echoserver", nil)
}

func TestPutAPIToKong(t *testing.T) {
	kong := createKongClient()
	req := KongPutAPISpec{
		Name:        "http-echoserver",
		UpstreamURL: "http://http-echoserver:80",
		URIs:        "/http-echoserver",
		StripURI:    true,
	}
	code, err := kong.PutNewAPI(&req)
	// or http.StatusOK(200)
	assert.Equal(t, http.StatusCreated, code)
	assert.Equal(t, nil, err)
	//time.Sleep(5 * time.Millisecond)
	/*
		code, err = kong.DeleteAPI("http-echoserver")
		assert.Equal(t, http.StatusNoContent, code)
		assert.Equal(t, nil, err)
	*/
}
