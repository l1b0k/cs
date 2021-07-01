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

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func init() {
	clusterCmd.AddCommand(clusterLsCmd)
}

var clusterLsCmd = &cobra.Command{
	Use:     "list",
	Short:   "list cluster",
	Long:    "list cluster",
	Aliases: []string{"ls"},
	Run: func(cmd *cobra.Command, args []string) {
		clusterLs()
	},
}

func clusterLs() {
	clusters, err := openAPI.DescribeClustersV1Request(context.TODO())
	cobra.CheckErr(err)

	table := pterm.DefaultTable.WithHasHeader(true)
	var data [][]string
	data = append(data, []string{"Name", "ClusterID", "Region", "State", "ClusterType", "CurrentVersion", "VPC", "SecurityGroupID"})

	for _, cluster := range clusters {
		data = append(data, []string{cluster.Name, cluster.ClusterID, cluster.RegionID, pterm.LightGreen(cluster.State), cluster.ClusterType, cluster.CurrentVersion, cluster.VPCID, cluster.SecurityGroupID})
	}

	_ = table.WithData(data).Render()
}
