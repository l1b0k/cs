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
	"strconv"
	"time"

	"github.com/docker/docker/api/types"
	docker "github.com/docker/docker/client"
)

// DockerService implements Interface
var _ Interface = &DockerService{}

type DockerService struct {
	timeout time.Duration
	client  *docker.Client
}

func NewDockerClient(timeout time.Duration) (*DockerService, error) {
	client, err := docker.NewClientWithOpts(
		docker.WithVersion("v1.21"),
		docker.WithTimeout(timeout),
	)
	if err != nil {
		return nil, err
	}

	return &DockerService{timeout: timeout, client: client}, nil
}

func (c *DockerService) ListContainer() ([]Container, error) {
	ctx, cancel := getContextWithTimeout(c.timeout)
	defer cancel()

	list, err := c.client.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		return nil, err
	}
	var containers []Container
	for _, s := range list {
		state, err := c.client.ContainerInspect(ctx, s.ID)
		if err != nil {
			continue
		}
		containers = append(containers, Container{
			ID:           s.ID,
			Pid:          strconv.Itoa(state.State.Pid),
			PodNamespace: s.Labels["io.kubernetes.pod.namespace"],
			PodName:      s.Labels["io.kubernetes.pod.name"],
			Type:         s.Labels["io.kubernetes.docker.type"],
		})
	}
	return containers, nil
}

func (c *DockerService) ListPodSandbox() ([]Container, error) {
	ctx, cancel := getContextWithTimeout(c.timeout)
	defer cancel()

	list, err := c.client.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		return nil, err
	}
	var containers []Container
	for _, s := range list {
		state, err := c.client.ContainerInspect(ctx, s.ID)
		if err != nil {
			continue
		}
		if s.Labels["io.kubernetes.docker.type"] != "podsandbox" {
			continue
		}
		containers = append(containers, Container{
			ID:           s.ID,
			Pid:          strconv.Itoa(state.State.Pid),
			PodNamespace: s.Labels["io.kubernetes.pod.namespace"],
			PodName:      s.Labels["io.kubernetes.pod.name"],
			Type:         s.Labels["io.kubernetes.docker.type"],
			State:        s.State,
		})
	}
	return containers, nil
}
