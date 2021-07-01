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

package runtime

import (
	"context"
	"fmt"
	"time"

	util2 "github.com/l1b0k/cs/pkg/runtime/util"
	"google.golang.org/grpc"
	runtimeapi "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
)

// RemoteRuntimeService implements Interface
var _ Interface = &RemoteRuntimeService{}

// RemoteRuntimeService is a gRPC implementation of internalapi.RuntimeService.
type RemoteRuntimeService struct {
	timeout       time.Duration
	runtimeClient runtimeapi.RuntimeServiceClient
}

// NewRemoteRuntimeService creates a new internalapi.RuntimeService.
func NewRemoteRuntimeService(endpoint string, connectionTimeout time.Duration) (*RemoteRuntimeService, error) {
	addr, dialer, err := util2.GetAddressAndDialer(endpoint)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), connectionTimeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, addr, grpc.WithInsecure(), grpc.WithContextDialer(dialer), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(maxMsgSize)))
	if err != nil {
		return nil, err
	}

	return &RemoteRuntimeService{
		timeout:       connectionTimeout,
		runtimeClient: runtimeapi.NewRuntimeServiceClient(conn),
	}, nil
}

// ListPodSandbox returns a list of PodSandboxes.
func (r *RemoteRuntimeService) ListPodSandbox() ([]Container, error) {
	ctx, cancel := getContextWithTimeout(r.timeout)
	defer cancel()

	var result []Container
	resp, err := r.runtimeClient.ListPodSandbox(ctx, &runtimeapi.ListPodSandboxRequest{})
	if err != nil {
		return nil, err
	}
	for _, c := range resp.Items {
		statusResp, err := r.runtimeClient.PodSandboxStatus(ctx, &runtimeapi.PodSandboxStatusRequest{
			PodSandboxId: c.Id,
			Verbose:      true,
		})
		if err != nil {
			return nil, err
		}
		for k, v := range statusResp.GetInfo() {
			fmt.Printf("%s %s\n", k, v)
		}
		result = append(result, Container{
			ID:           c.Id,
			Pid:          statusResp.GetInfo()[""],
			PodNamespace: c.GetMetadata().Namespace,
			PodName:      c.GetMetadata().Name,
			Type:         "",
		})

	}

	return result, nil
}
