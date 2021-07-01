// +build freebsd linux

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

package ns

import (
	"fmt"
	"os"

	"github.com/l1b0k/cs/pkg/netlink"
	"github.com/l1b0k/cs/pkg/views"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func init() {
	nsCmd.AddCommand(nsLsCmd)
}

var nsLsCmd = &cobra.Command{
	Use:     "list",
	Short:   "namespace list command",
	Long:    "namespace list command",
	Aliases: []string{"ls"},
	Run: func(cmd *cobra.Command, args []string) {
		nsLs()
	},
}

func nsLs() {
	pauses, err := runtimeClient.ListPodSandbox()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	table := pterm.DefaultTable.WithHasHeader(true)
	var data [][]string
	data = append(data, []string{"NS", "ID", "Pod", "State"})

	for _, pause := range pauses {
		data = append(data, []string{fmt.Sprintf("/proc/%s/ns", pause.Pid),
			pause.ID,
			fmt.Sprintf("%s/%s", pause.PodNamespace, pause.PodName),
			views.ContainerColor(pause.State),
		})
	}
	table.WithData(data).Render()

	pterm.DefaultSection.Println("netns")

	for _, pause := range pauses {
		pterm.DefaultSection.WithLevel(2).Println(fmt.Sprintf("%s %s/%s", pause.Pid, pause.PodNamespace, pause.PodName))
		err = netlink.GetNetInfo(pause.Pid)
		if err != nil {
			println(err.Error())
		}
	}
}
