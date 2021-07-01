// Copyright 2021 l1b0k
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

package cluster

import (
	"fmt"

	"github.com/l1b0k/cs/cmd"
	"github.com/l1b0k/cs/pkg/aliyun"
	"github.com/spf13/cobra"
)

// vars
var (
	openAPI *aliyun.OpenAPI
)

func init() {
	cmd.RootCmd.AddCommand(clusterCmd)
}

var clusterCmd = &cobra.Command{
	Use:     "cluster",
	Short:   "show cluster info",
	Long:    "show cluster info",
	Aliases: []string{"cs"},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		var err error
		openAPI, err = aliyun.NewOpenAPI()
		if err != nil {
			cobra.CheckErr(fmt.Errorf("failed init aliyun client, %w", err))
		}
	},
}
