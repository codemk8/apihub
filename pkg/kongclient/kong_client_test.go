package kongclient

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/codemk8/apihub/pkg/k8s"
	"github.com/go-resty/resty"
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
	initResty()
	response, err := resty.R().Get(kong.constructKongAPIUrl("/"))
	if err != nil {
		fmt.Printf("Wrong kong configuration %v", err)
	} else {
		resp := KongGetResp{}
		err := json.Unmarshal(response.Body(), &resp)
		if err != nil {
			fmt.Printf("Error parsing returned json %v", err)
		} else {
			fmt.Println(resp)
		}
	}
}

func TestRegisterServiceToKong(t *testing.T) {
	kong := createKongClient()
	kong.RegisterServiceToKong("default:http-echoserver")
}
