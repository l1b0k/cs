// Copyright 2020 l1b0k
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package views

import (
	"strings"

	"github.com/pterm/pterm"
)

func ContainerColor(key string) string {
	switch strings.ToLower(key) {
	case "running":
		return pterm.LightGreen(key)
	}
	return pterm.LightBlue(key)
}

func K8SColor(key string) string {
	switch key {
	case "Running":
		return pterm.LightGreen(key)
	case "Pending", "ContainerCreating":
		return pterm.LightYellow(key)
	case "CrashLoopBackoff":
		return pterm.LightRed(key)
	case "Completed":
		return pterm.LightWhite(key)
	}
	return pterm.LightBlue(key)
}

func ENIColor(key string) string {
	switch strings.ToLower(key) {
	case "secondary", "member", "trunk":
		return pterm.LightYellow(key)
	case "primary":
		return pterm.LightBlue(key)
	}
	return key
}

func IPColor(key string) string {
	return pterm.LightWhite(key)
}
