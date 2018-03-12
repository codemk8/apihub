package kongclient

import (
	"log"
	"net/http"
)

// Deploy implements the "deploy" command
func Deploy(service string, deployParams *DeployParams) bool {
	// TODO set params from cached values
	// use default now
	params := KongParams{}
	kong := NewKongK8sClient(params)

	if kong == nil {
		log.Println("Error init API gateway client, check if kong is a valid service in k8s.")
		return false
	}
	return kong.RegisterServiceToKong(service, deployParams)
}

// Remove implements the "remove" command
func Remove(names []string) bool {
	params := KongParams{}
	kong := NewKongK8sClient(params)

	if kong == nil {
		log.Println("Error init API gateway client, check if kong is a valid service in k8s.")
		return false
	}
	ok := true
	for _, name := range names {
		code, _ := kong.DeleteAPI(name)
		if code != http.StatusNoContent {
			log.Printf("Error delete API %s, err code %d\n", name, code)
			ok = false
		}
	}
	return ok
}

// List implements the "list" command
func List() *KongGetResp {
	params := KongParams{}
	kong := NewKongK8sClient(params)

	if kong == nil {
		log.Println("Error init API gateway client, check if kong is a valid service in k8s.")
		return nil
	}
	return kong.ListAPIs()
}
