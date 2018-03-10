package kongclient

import (
	"log"
	"reflect"
	"strconv"

	"github.com/codemk8/apihub/pkg/k8s"
	"k8s.io/client-go/kubernetes"
)

// KongParams is used to create KongK8sClient
type KongParams struct {
	GwNamespace string `default:"default"`
	GwName      string `default:"apihub-kong-kong-admin"`
	GwHost      string `default:"localhost"`
	GwPort      string
}

// DeployParams receives flags from command line
type DeployParams struct {
	Uris  string
	Force bool
	Name  string
}

// KongK8sClient sends API request to kong service in k8s
type KongK8sClient struct {
	K8sCs       *kubernetes.Clientset
	KongSvcNs   string
	KongSvcName string
	KongSvcHost string
	KongSvcPort string
}

// NewKongK8sClient creates a new KongK8sClient
func NewKongK8sClient(kongParams KongParams) *KongK8sClient {
	initResty()
	client, err := k8s.NewK8sClient()
	if err != nil {
		log.Printf("Error init k8s client: %v", err)
		return nil
	}
	kongClient := &KongK8sClient{
		K8sCs: client,
	}
	typ := reflect.TypeOf(kongParams)
	if kongParams.GwNamespace == "" {
		f, _ := typ.FieldByName("GwNamespace")
		kongClient.KongSvcNs = f.Tag.Get("default")
	}
	if kongParams.GwName == "" {
		f, _ := typ.FieldByName("GwName")
		kongClient.KongSvcName = f.Tag.Get("default")
	}
	if kongParams.GwHost == "" {
		f, _ := typ.FieldByName("GwHost")
		kongClient.KongSvcHost = f.Tag.Get("default")
	}

	port := k8s.GetServiceNodePort(kongClient.KongSvcNs, kongClient.KongSvcName, client)
	if len(port) != 1 {
		log.Printf("Error getting kong admin node port.")
		return nil
	}

	kongClient.KongSvcPort = strconv.Itoa(int(port[0]))
	return kongClient
}
