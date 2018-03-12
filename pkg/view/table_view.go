package view

import (
	"github.com/codemk8/apihub/pkg/kongclient"
)

// View interface
type View interface {
	show(*kongclient.KongGetResp)
}

// TableView shows the API list nicely
type TableView struct {
}

// View implements the View interface
func (v TableView) View(*kongclient.KongGetResp) {

}
