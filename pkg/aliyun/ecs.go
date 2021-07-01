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
	"fmt"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
)

func (a *OpenAPI) DescribeInstances(ctx context.Context, vpcID string) ([]ecs.Instance, error) {
	var result []ecs.Instance
	for i := 1; ; {
		req := ecs.CreateDescribeInstancesRequest()
		req.PageSize = requests.NewInteger(100)
		req.VpcId = vpcID
		req.RegionId = RegionFromContext(ctx)

		a.ReadOnlyRateLimiter.Accept()
		start := time.Now()
		resp, err := a.client.ECS().DescribeInstances(req)
		OpenAPILatency.WithLabelValues("DescribeInstances", fmt.Sprint(err != nil)).Observe(MsSince(start))
		if err != nil {
			return nil, err
		}
		result = append(result, resp.Instances.Instance...)

		if resp.TotalCount < resp.PageNumber*resp.PageSize {
			break
		}
		i++
	}
	return result, nil
}

func (a *OpenAPI) DescribeNetworkInterface(ctx context.Context, networkInterfaceID, nicType string) ([]ecs.NetworkInterfaceSet, error) {
	var result []ecs.NetworkInterfaceSet
	for i := 1; ; {
		req := ecs.CreateDescribeNetworkInterfacesRequest()
		req.PageSize = requests.NewInteger(100)
		req.NetworkInterfaceId = &[]string{networkInterfaceID}
		req.RegionId = RegionFromContext(ctx)
		req.Type = nicType

		a.ReadOnlyRateLimiter.Accept()
		start := time.Now()
		resp, err := a.client.ECS().DescribeNetworkInterfaces(req)
		OpenAPILatency.WithLabelValues("DescribeNetworkInterfaces", fmt.Sprint(err != nil)).Observe(MsSince(start))
		if err != nil {
			return nil, err
		}
		result = append(result, resp.NetworkInterfaceSets.NetworkInterfaceSet...)

		if resp.TotalCount < resp.PageNumber*resp.PageSize {
			break
		}
		i++
	}
	return result, nil
}
