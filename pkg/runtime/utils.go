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
	"os"
	"time"
)

// maxMsgSize use 16MB as the default message size limit.
// grpc library default is 4MB
const maxMsgSize = 1024 * 1024 * 16

// getContextWithTimeout returns a context with timeout.
func getContextWithTimeout(timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), timeout)
}

// getContextWithCancel returns a context with cancel.
func getContextWithCancel() (context.Context, context.CancelFunc) {
	return context.WithCancel(context.Background())
}

type Create func() Interface

type register struct {
	Path   string
	Client Create
}

func GetRuntime() Interface {
	eps := []register{
		{Path: "/var/run/docker.sock", Client: func() Interface {
			c, err := NewDockerClient(30 * time.Second)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			return c
		}},
		{Path: "/var/run/containerd/containerd.sock", Client: func() Interface {
			c, err := NewRemoteRuntimeService("/var/run/containerd/containerd.sock", 30*time.Second)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			return c
		}},
	}
	for _, ep := range eps {
		_, err := os.Stat(ep.Path)
		if err == nil {
			return ep.Client()
		}
		if os.IsNotExist(err) {
			continue
		}
		continue
	}
	fmt.Fprintln(os.Stderr, "no runtime client found")
	os.Exit(1)
	return nil
}
