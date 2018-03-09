package kongclient

import (
	"testing"

	"github.com/magiconair/properties/assert"
)

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
