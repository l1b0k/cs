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

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials/provider"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cs"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/spf13/viper"
	"k8s.io/client-go/util/flowcontrol"
)

type ClientInterface interface {
	VPC() *vpc.Client
	ECS() *ecs.Client
	CS() *cs.Client
}

func NewSimpleClient() (ClientInterface, error) {
	region := viper.GetString("region")
	c := &SimpleClient{}
	vC, err := vpc.NewClientWithProvider(region, provider.DefaultChain)
	if err != nil {
		return nil, err
	}
	c.vpc = vC
	eC, err := ecs.NewClientWithProvider(region, provider.DefaultChain)
	if err != nil {
		return nil, err
	}
	c.ecs = eC
	cC, err := cs.NewClientWithProvider(region, provider.DefaultChain)
	if err != nil {
		return nil, err
	}
	c.cs = cC

	if viper.GetBool("private") {
		c.vpc.SetEndpointRules(c.ecs.EndpointMap, "regional", "vpc")
		c.ecs.SetEndpointRules(c.ecs.EndpointMap, "regional", "vpc")
		c.cs.SetEndpointRules(c.ecs.EndpointMap, "regional", "vpc")
	}
	return c, nil
}

type SimpleClient struct {
	vpc *vpc.Client
	ecs *ecs.Client
	cs  *cs.Client
}

func (c *SimpleClient) VPC() *vpc.Client {
	return c.vpc
}

func (c *SimpleClient) ECS() *ecs.Client {
	return c.ecs
}

func (c *SimpleClient) CS() *cs.Client {
	return c.cs
}

type OpenAPI struct {
	client ClientInterface

	ReadOnlyRateLimiter flowcontrol.RateLimiter
}

func NewOpenAPI() (*OpenAPI, error) {
	c, err := NewSimpleClient()
	if err != nil {
		return nil, err
	}
	return &OpenAPI{
		client:              c,
		ReadOnlyRateLimiter: flowcontrol.NewTokenBucketRateLimiter(8, 10),
	}, nil
}

type regionContextKey struct{}

func RegionFromContext(ctx context.Context) string {
	if v, ok := ctx.Value(regionContextKey{}).(string); ok {
		return v
	}
	return ""
}

func WithRegion(ctx context.Context, region string) context.Context {
	return context.WithValue(ctx, regionContextKey{}, region)
}
