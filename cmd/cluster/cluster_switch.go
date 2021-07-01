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
	"context"
	"fmt"
	"os"

	"github.com/l1b0k/cs/pkg/utils"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	clusterCmd.AddCommand(clusterSwitchCmd)
}

var clusterSwitchCmd = &cobra.Command{
	Use:   "switch",
	Short: "switch kubernetes ctx to this cluster",
	Long:  "switch kubernetes ctx to this cluster",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		clusterSwitch(args[0])
	},
}

func clusterSwitch(clusterID string) {
	cluster, err := openAPI.DescribeClusterDetail(context.TODO(), clusterID)
	cobra.CheckErr(err)

	kubeconfig, err := openAPI.DescribeClusterUserKubeconfig(context.TODO(), clusterID, viper.GetBool("private"))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	apiCfg, err := utils.SetKubeConfigName(kubeconfig, cluster.ClusterID, cluster.Name)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	mergedConf, err := utils.MergeApiConfig(apiCfg, true)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	table := pterm.DefaultTable.WithHasHeader(true)
	var data [][]string
	data = append(data, []string{"Current", "Name", "Cluster", "API"})
	for name, confCtx := range mergedConf.Contexts {
		line := []string{}
		if mergedConf.CurrentContext == name {
			line = append(line, "*")
		} else {
			line = append(line, "")
		}
		line = append(line, name)
		line = append(line, confCtx.Cluster)

		cluster, ok := mergedConf.Clusters[confCtx.Cluster]
		if !ok {
			line = append(line, "")
		} else {
			line = append(line, cluster.Server)
		}
		data = append(data, line)
	}
	_ = table.WithData(data).Render()
}
