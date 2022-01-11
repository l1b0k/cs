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
	"net"
	"strings"

	"github.com/l1b0k/cs/pkg/aliyun"
	"github.com/l1b0k/cs/pkg/utils"
	"github.com/l1b0k/cs/pkg/views"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func init() {
	clusterCmd.AddCommand(clusterInspectCmd)
}

var clusterInspectCmd = &cobra.Command{
	Use:   "inspect",
	Short: "inspect cluster",
	Long:  "inspect cluster",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		clusterInspect(args[0])
	},
}

func clusterInspect(clusterID string) {
	ctx := context.TODO()

	pterm.DefaultSection.Println("ACK Cluster")

	// 1. cluster info
	cluster, err := openAPI.DescribeClusterDetail(ctx, clusterID)
	cobra.CheckErr(err)

	t := pterm.DefaultTable.WithHasHeader(true)
	var data [][]string
	data = append(data, []string{"Name", "ClusterID", "Region", "State", "ClusterType", "CurrentVersion", "VPC", "SecurityGroupID"},
		[]string{cluster.Name, cluster.ClusterID, cluster.RegionID, cluster.State, cluster.ClusterType, cluster.CurrentVersion, cluster.VPCID, cluster.SecurityGroupID})
	cobra.CheckErr(t.WithData(data).Render())

	kubeconfig, err := openAPI.DescribeClusterUserKubeconfig(context.TODO(), clusterID, viper.GetBool("private"))
	cobra.CheckErr(err)

	k8sClient, err := utils.ClientFromKubeConfig(kubeconfig)
	cobra.CheckErr(err)

	pterm.DefaultSection.WithLevel(2).Println("Nodes")

	// 2. node info
	nodes, err := GetNodes(k8sClient, nil)
	cobra.CheckErr(err)

	data = [][]string{
		{"Name", "InstanceID", "InternalIP", "Zone", "InstanceType"},
	}
	for _, node := range nodes {
		data = append(data, []string{node.Name, node.InstanceID, node.InternalIP.String(), node.Zone, node.InstanceType})
	}
	cobra.CheckErr(t.WithData(data).Render())

	pterm.DefaultSection.Println("VPC")

	ctx = aliyun.WithRegion(ctx, cluster.RegionID)
	// 3. show vpc info
	vsws, err := openAPI.DescribeVSwitches(ctx, cluster.VPCID)
	cobra.CheckErr(err)

	data = [][]string{
		{"vSwitch", "IPv4", "Count", "Zone"},
	}
	for _, vsw := range vsws {
		data = append(data, []string{vsw.VSwitchId, vsw.CidrBlock, fmt.Sprintf("%d", vsw.AvailableIpAddressCount), vsw.ZoneId})
	}
	cobra.CheckErr(t.WithData(data).Render())

	pterm.DefaultSection.WithLevel(2).Println("ECS instance info")
	// 4. show ecs info
	instances, err := openAPI.DescribeInstances(ctx, cluster.VPCID)
	cobra.CheckErr(err)

	for _, node := range nodes {
		var ins *ecs.Instance
		for _, i := range instances {
			if i.InstanceId == node.InstanceID {
				ins = &i
				break
			}
		}
		if ins == nil {
			continue
		}
		var publicIP string
		if len(ins.PublicIpAddress.IpAddress) > 0 {
			publicIP = ins.PublicIpAddress.IpAddress[0]
		}
		if ins.EipAddress.IpAddress != "" {
			publicIP = ins.EipAddress.IpAddress
		}

		t := pterm.DefaultTable.WithHasHeader(true)
		var data [][]string
		data = append(data, []string{"Name", "InstanceID", "InternalIP", "Zone", "PublicIP"},
			[]string{node.Name, node.InstanceID, node.InternalIP.String(), node.Zone, publicIP})
		_ = t.WithData(data).Render()

		pterm.DefaultSection.WithLevel(3).Println("ENI info")

		eniList := pterm.LeveledList{}
		for _, ins := range instances {
			if ins.InstanceId != node.InstanceID {
				continue
			}

			for _, eniInfo := range ins.NetworkInterfaces.NetworkInterface {
				eniSet, err := openAPI.DescribeNetworkInterface(ctx, eniInfo.NetworkInterfaceId, "")
				if err != nil {
					cobra.CheckErr(err)
				}
				if len(eniSet) == 0 {
					continue
				}
				eni := eniSet[0]
				eniList = append(eniList, pterm.LeveledListItem{Level: 0, Text: fmt.Sprintf("%s %s", eni.NetworkInterfaceId, views.ENIColor(eni.Type))})
				eniList = append(eniList, pterm.LeveledListItem{Level: 1, Text: fmt.Sprintf("%s %s", eni.VSwitchId, strings.Join(eni.SecurityGroupIds.SecurityGroupId, ","))})
				eniList = append(eniList, pterm.LeveledListItem{Level: 1, Text: fmt.Sprintf("%s %s", eni.MacAddress, views.IPColor(eni.PrivateIpAddress))})
				for _, addr := range eni.PrivateIpSets.PrivateIpSet {
					eniList = append(eniList, pterm.LeveledListItem{Level: 2, Text: views.IPColor(addr.PrivateIpAddress)})
				}
				for _, addr := range eni.Ipv6Sets.Ipv6Set {
					eniList = append(eniList, pterm.LeveledListItem{Level: 2, Text: views.IPColor(addr.Ipv6Address)})
				}
				if eni.Type == "Trunk" {
					memberENISet, err := openAPI.DescribeNetworkInterface(ctx, "", "Member")
					if err != nil {
						cobra.CheckErr(err)
					}
					for _, memberENI := range memberENISet {
						if memberENI.Attachment.TrunkNetworkInterfaceId != eni.NetworkInterfaceId {
							continue
						}
						eniList = append(eniList, pterm.LeveledListItem{Level: 1, Text: fmt.Sprintf("%s %s", memberENI.NetworkInterfaceId, views.ENIColor(memberENI.Type))})
						eniList = append(eniList, pterm.LeveledListItem{Level: 2, Text: fmt.Sprintf("%s %s", memberENI.VSwitchId, strings.Join(memberENI.SecurityGroupIds.SecurityGroupId, ","))})
						eniList = append(eniList, pterm.LeveledListItem{Level: 2, Text: fmt.Sprintf("%s %s", memberENI.MacAddress, views.IPColor(memberENI.PrivateIpAddress))})
					}
				}
			}
		}
		root := pterm.NewTreeFromLeveledList(eniList)
		_ = pterm.DefaultTree.WithRoot(root).Render()
	}
}

// Node is the k8s node struct
type Node struct {
	Name         string
	InstanceID   string
	InternalIP   net.IP
	Region       string
	Zone         string
	InstanceType string
}

func GetNodes(client kubernetes.Interface, filterName []string) ([]Node, error) {
	nodes := []Node{}
	k8sNodes, err := client.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, k8sNode := range k8sNodes.Items {
		found := false
		if len(filterName) == 0 {
			found = true
		}
		for _, n := range filterName {
			if n == k8sNode.Name {
				found = true
				break
			}
		}
		if !found {
			continue
		}
		node := Node{
			Name:         k8sNode.Name,
			InstanceID:   utils.ParseInstanceID(&k8sNode),
			InternalIP:   utils.ParseInternalIP(&k8sNode),
			Region:       utils.ParseRegion(&k8sNode),
			Zone:         utils.ParseZone(&k8sNode),
			InstanceType: utils.ParseInstanceType(&k8sNode),
		}
		nodes = append(nodes, node)
	}
	return nodes, nil
}
