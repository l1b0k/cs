package cluster

import (
	"os"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials/provider"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cs"
	"github.com/spf13/cobra"
)

func init() {
	clusterCmd.AddCommand(clusterCreateCmd)
}

var clusterCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "create cluster",
	Long:  "create cluster",
	Example: `
refer: https://help.aliyun.com/document_detail/93084.htm

Example config file:

{
  "name": "cs_hangzhou_mip_ipvlan",
  "cluster_type": "ManagedKubernetes",
  "disable_rollback": true,
  "timeout_mins": 60,
  "kubernetes_version": "1.20.11-aliyun.1",
  "region_id": "cn-hangzhou",
  "snat_entry": true,
  "cloud_monitor_flags": false,
  "endpoint_public_access": true,
  "deletion_protection": false,
  "proxy_mode": "ipvs",
  "tags": [],
  "timezone": "Asia/Shanghai",
  "addons": [
    {
      "name": "terway-eniip",
      "config": "{\"IPVlan\":\"true\",\"NetworkPolicy\":\"false\"}"
    },
    {
      "name": "csi-plugin"
    },
    {
      "name": "csi-provisioner"
    },
    {
      "name": "nginx-ingress-controller",
      "disabled": true
    }
  ],
  "cluster_spec": "ack.pro.small",
  "load_balancer_spec": "slb.s2.small",
  "os_type": "Linux",
  "platform": "AliyunLinux",
  "image_type": "AliyunLinux",
  "pod_vswitch_ids": [
    "POD_VSWITCHES"
  ],
  "runtime": {
    "name": "docker",
    "version": "19.03.15"
  },
  "worker_instance_types": [
    "ecs.g5ne.2xlarge"
  ],
  "num_of_nodes": 2,
  "worker_system_disk_category": "cloud_ssd",
  "worker_system_disk_size": 120,
  "worker_system_disk_performance_level": "",
  "charge_type": "PostPaid",
  "vpcid": "YOUR_VPC",
  "service_cidr": "172.16.0.0/16",
  "vswitch_ids": [
    "NODE_VSWITCHES"
  ],
  "key_pair": "SSH_KEY",
  "cpu_policy": "none",
  "is_enterprise_security_group": true
}
`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		clusterCreate(args[0])
	},
}

func clusterCreate(configPath string) {
	client, err := cs.NewClientWithProvider("cn-hangzhou", provider.DefaultChain)
	if err != nil {
		panic(err)
	}
	req := cs.CreateCreateClusterRequest()
	req.Method = "POST"
	req.Scheme = "https"
	req.Domain = "cs.cn-hangzhou.aliyuncs.com"
	req.Headers["Content-Type"] = "application/json"
	req.RegionId = "cn-hangzhou"

	c, err := os.ReadFile(configPath)
	if err != nil {
		panic(err)
	}
	req.Content = c
	resp, err := client.CreateCluster(req)
	if err != nil {
		panic(err)
	}
	println(resp.String())
}
