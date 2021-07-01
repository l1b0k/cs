# cs

a CLI for ACK cluster

## commands

- cluster
    - ls: list clusters
    - inspect {cluster_id}: get cluster detail
    - switch {cluster_id}: set kubeconfig to this cluster
    - ns ls : list netns info

## usage

```sh
export ALIBABA_CLOUD_REGION_ID=cn-hangzhou
export ALIBABA_CLOUD_ACCESS_KEY_ID=
export ALIBABA_CLOUD_ACCESS_KEY_SECRET=

❯ cs cluster ls
Name          | ClusterID                         | Region      | State   | ClusterType       | CurrentVersion  | VPC                       | SecurityGroupID
cluster-dev | cf6c67398726b4a17******** | cn-hangzhou | running | ManagedKubernetes | 1.18.8-aliyun.1 | vpc-bp11nakc4ux0iwl40mi3e | sg-bp18sxd8cxak178otfm8

❯ cs cluster switch cf6c67398726b4a17********
Current | Name          | Cluster                                   | API
*       | cluster-dev | cluster-cf6c67398726b4a17******** | https://114.*.*.*:6443

❯ cs cluster inspect cf6c67398726b4a17********

# ACK Cluster

Name          | ClusterID                         | Region      | State   | ClusterType       | CurrentVersion  | VPC                       | SecurityGroupID
cluster-dev | cf6c67398726b4a17******** | cn-hangzhou | running | ManagedKubernetes | 1.18.8-aliyun.1 | vpc-bp11nakc4ux0iwl40mi3e | sg-bp18sxd8cxak178otfm8

## Nodes

+---------------------------+------------------------+---------------+---------------+
|           NAME            |       INSTANCEID       |  INTERNALIP   |     ZONE      |
+---------------------------+------------------------+---------------+---------------+
| cn-hangzhou.192.168.12.42 | i-bp1aa3aakcsrhcuz2tkz | 192.168.12.42 | cn-hangzhou-i |
| cn-hangzhou.192.168.12.43 | i-bp1aa3aakcsrhcuz2tl0 | 192.168.12.43 | cn-hangzhou-i |
+---------------------------+------------------------+---------------+---------------+

# VPC

+---------------------------+-----------------+-------+---------------+
|          VSWITCH          |      IPV4       | COUNT |     ZONE      |
+---------------------------+-----------------+-------+---------------+
| vsw-bp154vjnrydhsg16xus7o | 192.168.0.0/19  |  8182 | cn-hangzhou-i |
| vsw-bp1k4vsi1dnuldcr2ddnb | 192.168.32.0/19 |  8180 | cn-hangzhou-i |
| vsw-bp14ijorh5upbgzkonlzk | 192.168.96.0/19 |  8185 | cn-hangzhou-i |
| vsw-bp1o7r5x86br1mgjxr9xk | 192.168.64.0/19 |  8188 | cn-hangzhou-i |
+---------------------------+-----------------+-------+---------------+

## ECS instance info

Name                      | InstanceID             | InternalIP    | Zone
cn-hangzhou.192.168.12.42 | i-bp1aa3aakcsrhcuz2tkz | 192.168.12.42 | cn-hangzhou-i

### ENI info

├─┬eni-bp11vypf0x3t6sf75yqt Primary
│ ├──
│ └─┬00:16:3e:17:c6:4b 192.168.12.42
│   └──192.168.12.42
└─┬eni-bp1cgcw2x546axiwcxoy Secondary
  ├──
  └─┬00:16:3e:0f:0c:c0 192.168.41.124
    └──192.168.41.124

Name                      | InstanceID             | InternalIP    | Zone
cn-hangzhou.192.168.12.43 | i-bp1aa3aakcsrhcuz2tl0 | 192.168.12.43 | cn-hangzhou-i

### ENI info

├─┬eni-bp1hflnr86s5x0533re2 Primary
│ ├──
│ └─┬00:16:3e:15:a9:4d 192.168.12.43
│   └──192.168.12.43
└─┬eni-bp136qodaanft39r5xwe Secondary
  ├──
  └─┬00:16:3e:0e:54:aa 192.168.121.118
    └──192.168.121.118

```

    