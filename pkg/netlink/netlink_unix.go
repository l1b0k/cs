// +build freebsd linux darwin

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

package netlink

import (
	"fmt"
	"strconv"

	"github.com/containernetworking/plugins/pkg/ns"
	"github.com/pterm/pterm"
	"github.com/vishvananda/netlink"
)

func GetNetInfo(pid string) error {
	netNS, err := ns.GetNS(fmt.Sprintf("/proc/%s/ns/net", pid))
	if err != nil {
		return err
	}
	defer netNS.Close()

	err = netNS.Do(func(netNS ns.NetNS) error {
		links, err := netlink.LinkList()
		if err != nil {
			return err
		}
		for _, link := range links {
			pterm.DefaultSection.WithLevel(3).Println("link %s", link.Attrs().Name)

			var data [][]string
			data = append(data, []string{link.Attrs().Name,
				strconv.Itoa(link.Attrs().Index),
				link.Type(),
				link.Attrs().HardwareAddr.String(),
				link.Attrs().Flags.String(),
				//link.Attrs().Slave.SlaveType(),
			})
			pterm.DefaultTable.WithData(data).Render()

			qdiscs, err := netlink.QdiscList(link)
			if err != nil {
				return err
			}
			for _, qdisc := range qdiscs {
				pterm.Printf("%s %s\n", qdisc.Type(), qdisc.Attrs().String())
			}
			pterm.DefaultSection.WithLevel(4).Println("filter")

			// ingress
			ingFilters, err := netlink.FilterList(link, netlink.HANDLE_MIN_INGRESS)
			if err != nil {
				return err
			}
			for _, filter := range ingFilters {
				pterm.Printf("%s %s\n", filter.Type(), filter.Attrs().String())
			}

			egressFilters, err := netlink.FilterList(link, netlink.HANDLE_MIN_EGRESS)
			if err != nil {
				return err
			}
			for _, filter := range egressFilters {
				pterm.Printf("%s %s\n", filter.Type(), filter.Attrs().String())
			}
		}

		return nil
	})

	return nil
}
