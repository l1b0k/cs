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

package aliyun

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cs"
)

type Clusters struct {
	Clusters []Cluster `json:"clusters,omitempty"`
	PageInfo PageInfo  `json:"page_info"`
}

type Cluster struct {
	Name                   string    `json:"name"`
	ClusterID              string    `json:"cluster_id"`
	RegionID               string    `json:"region_id"`
	Size                   int       `json:"size"`
	State                  string    `json:"state"`
	ClusterType            string    `json:"cluster_type"`
	Created                time.Time `json:"created"`
	Updated                time.Time `json:"updated"`
	CurrentVersion         string    `json:"current_version"`
	ResourceGroupID        string    `json:"resource_group_id"`
	VPCID                  string    `json:"vpc_id"`
	VSwitchID              string    `json:"vswitch_id"`
	SecurityGroupID        string    `json:"security_group_id"`
	ZoneID                 string    `json:"zone_id"`
	ExternalLoadbalancerID string    `json:"external_loadbalancer_id"`
}

type PageInfo struct {
	TotalCount int `json:"total_count"`
	PageNumber int `json:"page_number"`
	PageSize   int `json:"page_size"`
}

type DescribeClusterUserKubeconfigResp struct {
	Config string `json:"config"`
}

func (a *OpenAPI) DescribeClustersV1Request(ctx context.Context) ([]Cluster, error) {
	var result []Cluster
	for i := 1; ; {
		req := cs.CreateDescribeClustersV1Request()
		req.PageSize = requests.NewInteger(50)
		req.RegionId = RegionFromContext(ctx)

		a.ReadOnlyRateLimiter.Accept()
		start := time.Now()
		resp, err := a.client.CS().DescribeClustersV1(req)
		OpenAPILatency.WithLabelValues("DescribeClustersV1", fmt.Sprint(err != nil)).Observe(MsSince(start))
		if err != nil {
			return nil, err
		}

		clusters := Clusters{}
		err = json.Unmarshal(resp.BaseResponse.GetHttpContentBytes(), &clusters)
		if err != nil {
			return nil, err
		}

		result = append(result, clusters.Clusters...)

		if clusters.PageInfo.TotalCount < clusters.PageInfo.PageNumber*clusters.PageInfo.PageSize {
			break
		}
		i++
	}
	return result, nil
}

func (a *OpenAPI) DescribeClusterDetail(ctx context.Context, clusterID string) (*Cluster, error) {
	req := cs.CreateDescribeClusterDetailRequest()
	req.ClusterId = clusterID
	req.RegionId = RegionFromContext(ctx)

	a.ReadOnlyRateLimiter.Accept()
	start := time.Now()
	resp, err := a.client.CS().DescribeClusterDetail(req)
	OpenAPILatency.WithLabelValues("DescribeClusterDetail", fmt.Sprint(err != nil)).Observe(MsSince(start))
	if err != nil {
		return nil, err
	}
	cluster := &Cluster{}
	err = json.Unmarshal(resp.BaseResponse.GetHttpContentBytes(), cluster)
	if err != nil {
		return nil, err
	}

	return cluster, nil
}

func (a *OpenAPI) DescribeClusterUserKubeconfig(ctx context.Context, clusterID string, privateAddr bool) ([]byte, error) {
	req := cs.CreateDescribeClusterUserKubeconfigRequest()
	req.ClusterId = clusterID
	req.PrivateIpAddress = requests.NewBoolean(privateAddr)
	req.RegionId = RegionFromContext(ctx)

	a.ReadOnlyRateLimiter.Accept()
	start := time.Now()
	resp, err := a.client.CS().DescribeClusterUserKubeconfig(req)
	OpenAPILatency.WithLabelValues("DescribeClusterUserKubeconfig", fmt.Sprint(err != nil)).Observe(MsSince(start))
	if err != nil {
		return nil, err
	}

	cfg := &DescribeClusterUserKubeconfigResp{}
	err = json.Unmarshal(resp.BaseResponse.GetHttpContentBytes(), cfg)
	if err != nil {
		return nil, err
	}

	return []byte(cfg.Config), nil
}
