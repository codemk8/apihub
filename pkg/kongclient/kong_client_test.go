package kongclient

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
}

func TestNewKongClient(t *testing.T) {
	params := KongParams{}

	kong := NewKongK8sClient(params)
	assert.Equal(t, "localhost", kong.KongSvcHost)
}

func TestRegisterServiceToKong(t *testing.T) {
	params := KongParams{}
	kong := NewKongK8sClient(params)
	kong.RegisterServiceToKong("default:http-echoserver")
}
