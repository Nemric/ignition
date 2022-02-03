// Copyright 2015 CoreOS, Inc.
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

package util

import (
	"path/filepath"

	"github.com/coreos/ignition/v2/config/v3_4_experimental/types"
)

func (u Util) SystemdUnitPaths(unit types.Unit) []string {
	var paths []string
	switch GetUnitScope(unit) {
	case UserUnit:
		for _, user := range unit.Users {
			home, err := u.GetUserHomeDirByName(string(user))
			if err != nil {
				print(home, err)
			}
			paths = append(paths, filepath.Join(home, ".config", "systemd", "user"))
		}
	case SystemUnit:
		paths = append(paths, filepath.Join("etc", "systemd", "system"))
	case GlobalUnit:
		paths = append(paths, filepath.Join("etc", "systemd", "user"))
	default:
		paths = append(paths, filepath.Join("etc", "systemd", "system"))
	}
	return paths
}

func (u Util) SystemdPresetPath(scope UnitScope) string {
	switch scope {
	case UserUnit:
		return filepath.Join("etc", "systemd", "user-preset", "21-ignition-user.preset")
	case SystemUnit:
		return filepath.Join("etc", "systemd", "system-preset", "20-ignition.preset")
	case GlobalUnit:
		return filepath.Join("etc", "systemd", "user-preset", "20-ignition-global.preset")
	default:
		return filepath.Join("etc", "systemd", "system-preset", "20-ignition.preset")
	}
}

func (u Util) SystemdDropinsPaths(unit types.Unit) []string {
	var paths []string
	for _, path := range u.SystemdUnitPaths(unit) {
		paths = append(paths, filepath.Join(path, unit.Name+".d"))
	}
	return paths
}
