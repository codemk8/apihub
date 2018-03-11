// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"

	"github.com/codemk8/apihub/pkg/kongclient"
	"github.com/spf13/cobra"
)

var (
	name     string
	uris     string
	force    bool
	stripURI bool
	subPath  string
)

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("Please specify one service name to deploy.")
			return
		}
		params := &kongclient.DeployParams{
			Uris:     uris,
			Force:    force,
			Name:     name,
			StripURI: stripURI,
		}
		ok := kongclient.Deploy(args[0], params)
		if ok {
			fmt.Printf("Deploy successfully")
		}
	},
}

func init() {
	rootCmd.AddCommand(deployCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deployCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deployCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	deployCmd.Flags().StringVarP(&uris, "uris", "u", "/", "path to the API gateway")
	deployCmd.Flags().StringVarP(&name, "name", "n", "", "API name")
	deployCmd.Flags().StringVarP(&subPath, "subpath", "p", "", "Subpath after the service name, e.g. http://myapi/api/v1 needs to specify subpath as \"api/v1\"")
	deployCmd.Flags().BoolVarP(&force, "force", "f", false, "Force to delete existing ones")
	deployCmd.Flags().BoolVarP(&stripURI, "strip", "s", true, "Strip the matching prefix from the upstream URI")
}
