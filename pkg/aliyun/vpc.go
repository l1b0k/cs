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
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
)

func (a *OpenAPI) DescribeVSwitches(ctx context.Context, vpcID string) ([]vpc.VSwitch, error) {
	var result []vpc.VSwitch
	for i := 1; ; {
		req := vpc.CreateDescribeVSwitchesRequest()
		req.PageSize = requests.NewInteger(50)
		req.VpcId = vpcID
		req.RegionId = RegionFromContext(ctx)

		a.ReadOnlyRateLimiter.Accept()
		start := time.Now()
		resp, err := a.client.VPC().DescribeVSwitches(req)
		OpenAPILatency.WithLabelValues("DescribeVSwitches", fmt.Sprint(err != nil)).Observe(MsSince(start))
		if err != nil {
			return nil, err
		}
		result = append(result, resp.VSwitches.VSwitch...)

		if resp.TotalCount < resp.PageNumber*resp.PageSize {
			break
		}
		i++
	}
	return result, nil
}
