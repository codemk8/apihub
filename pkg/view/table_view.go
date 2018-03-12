package view

import (
	"os"
	"strconv"
	"strings"

	"github.com/codemk8/apihub/pkg/kongclient"
	"github.com/olekukonko/tablewriter"
)

// View interface
type View interface {
	show(*kongclient.KongGetResp)
}

// TableView shows the API list nicely
type TableView struct {
	Table *tablewriter.Table
}

// NewTableView creates the Table
func NewTableView() *TableView {
	return &TableView{
		Table: tablewriter.NewWriter(os.Stdout),
	}
}

// View implements the View interface
func (v TableView) View(apis *kongclient.KongGetResp) {
	data := [][]string{}
	for _, api := range apis.Data {
		data = append(data, []string{api.Name, api.UpstreamURL,
			strconv.FormatBool(api.StripURI), strings.Join(api.Uris, " ")})
	}

	v.Table.SetHeader([]string{"Name", "UpstreamURL", "StripURI", "Uris"})
	v.Table.SetFooter([]string{"", "", "Total", strconv.Itoa(apis.Total)})
	v.Table.SetBorder(true)
	v.Table.AppendBulk(data)
	v.Table.Render()
}
