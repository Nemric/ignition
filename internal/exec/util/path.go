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

func SystemdUnitsPath(unit types.Unit) string {
	return unit.GetBasePath()
}

func SystemdPresetPath(unit types.Unit) string {
	switch unit.GetScope() {
	case types.UserUnit:
		return filepath.Join("etc", "systemd", "user-preset", "20-ignition.preset")
	case types.SystemUnit:
		return filepath.Join("etc", "systemd", "system-preset", "20-ignition.preset")
	default:
		return filepath.Join("etc", "systemd", "system-preset", "20-ignition.preset")
	}
}

func SystemdWantsPath(unit types.Unit) string {
	return filepath.Join(SystemdUnitsPath(unit), unit.Name+".wants")
}

func SystemdDropinsPath(unit types.Unit) string {
	return filepath.Join(SystemdUnitsPath(unit), unit.Name+".d")
}
