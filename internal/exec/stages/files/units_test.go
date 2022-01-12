// Copyright 2020 CoreOS, Inc.
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

package files

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/coreos/ignition/v2/config/shared/errors"
	cfgutil "github.com/coreos/ignition/v2/config/util"
	"github.com/coreos/ignition/v2/config/v3_4_experimental"
	"github.com/coreos/ignition/v2/config/v3_4_experimental/types"
	"github.com/coreos/ignition/v2/internal/exec/util"
	"github.com/coreos/ignition/v2/internal/log"
)

func TestParseInstanceUnit(t *testing.T) {
	type in struct {
		unit types.Unit
	}
	type out struct {
		unitName string
		instance string
		parseErr error
	}
	tests := []struct {
		in  in
		out out
	}{
		{in: in{types.Unit{Name: "echo@bar.service"}},
			out: out{unitName: "echo@.service", instance: "bar",
				parseErr: nil},
		},

		{in: in{types.Unit{Name: "echo@foo.service"}},
			out: out{unitName: "echo@.service", instance: "foo",
				parseErr: nil},
		},
		{in: in{types.Unit{Name: "echo.service"}},
			out: out{unitName: "", instance: "",
				parseErr: errors.ErrInvalidInstantiatedUnit},
		},
		{in: in{types.Unit{Name: "echo@fooservice"}},
			out: out{unitName: "", instance: "",
				parseErr: errors.ErrNoSystemdExt},
		},
		{in: in{types.Unit{Name: "echo@.service"}},
			out: out{unitName: "echo@.service", instance: "",
				parseErr: nil},
		},
		{in: in{types.Unit{Name: "postgresql@9.3-main.service"}},
			out: out{unitName: "postgresql@.service", instance: "9.3-main",
				parseErr: nil},
		},
	}
	for i, test := range tests {
		unitName, instance, err := parseInstanceUnit(test.in.unit)
		if test.out.parseErr != err {
			t.Errorf("#%d: bad error: want %v, got %v", i, test.out.parseErr, err)
		}
		if !reflect.DeepEqual(test.out.unitName, unitName) {
			t.Errorf("#%d: bad unitName: want %v, got %v", i, test.out.unitName, unitName)
		}
		if !reflect.DeepEqual(test.out.instance, instance) {
			t.Errorf("#%d: bad instance: want %v, got %v", i, test.out.instance, instance)
		}
	}
}

func TestSystemdUnitPath(t *testing.T) {

	tests := []struct {
		in  types.Unit
		out string
	}{
		{
			types.Unit{Name: "test.service", Scope: cfgutil.StrToPtr("system")},
			"etc/systemd/system",
		},
		{
			types.Unit{Name: "test.service"},
			"etc/systemd/system",
		},
		{
			types.Unit{Name: "test.service", Scope: cfgutil.StrToPtr("user")},
			"etc/systemd/user",
		},
	}

	for i, test := range tests {
		path := util.SystemdUnitsPath(test.in)
		if path != test.out {
			t.Errorf("#%d: bad error: want %v, got %v", i, test.out, path)
		}
	}
}

func TestSystemdDropinsPath(t *testing.T) {

	tests := []struct {
		in  types.Unit
		out string
	}{
		{
			types.Unit{Name: "test.service", Scope: cfgutil.StrToPtr("system")},
			"etc/systemd/system/test.service.d",
		},
		{
			types.Unit{Name: "test.service"},
			"etc/systemd/system/test.service.d",
		},
		{
			types.Unit{Name: "test.service", Scope: cfgutil.StrToPtr("user")},
			"etc/systemd/user/test.service.d",
		},
	}

	for i, test := range tests {
		path := util.SystemdDropinsPath(test.in)
		if path != test.out {
			t.Errorf("#%d: bad error: want %v, got %v", i, test.out, path)
		}
	}
}

func TestSystemdPresetPath(t *testing.T) {

	tests := []struct {
		in  types.Unit
		out string
	}{
		{
			types.Unit{Name: "test.service", Scope: cfgutil.StrToPtr("system")},
			"etc/systemd/system-preset/20-ignition.preset",
		},
		{
			types.Unit{Name: "test.service"},
			"etc/systemd/system-preset/20-ignition.preset",
		},
		{
			types.Unit{Name: "test.service", Scope: cfgutil.StrToPtr("user")},
			"etc/systemd/user-preset/20-ignition.preset",
		},
	}

	for i, test := range tests {
		path := util.SystemdPresetPath(test.in)
		if path != test.out {
			t.Errorf("#%d: bad error: want %v, got %v", i, test.out, path)
		}
	}
}

func TestCreateUnits(t *testing.T) {

	config, report, err := v3_4_experimental.Parse([]byte(`{"ignition":{"version":"3.4.0-experimental"},"systemd":{"units":[{"contents":"[Unit]\nDescription=Prometheus node exporter\n[Install]\nWantedBy=multi-user.target\n","enabled":true,"name":"exporter.service"},{"contents":"[Unit]\nDescription=promtail.service\n[Install]\nWantedBy=multi-user.target default.target","enabled":true,"name":"promtail.service","scope":"user"},{"contents":"[Unit]\nDescription=promtail.service\n[Install]\nWantedBy=multi-user.target default.target","enabled":true,"name":"grafana.service","scope":"system"}]}}`))

	if err != nil {
		print(report.Entries, err.Error())
	}

	fmt.Printf("config: %v\n", config)

	tests := []struct {
		in  types.Config
		out error
	}{
		{
			config,
			nil,
		},
	}

	for i, test := range tests {
		var logg log.Logger = log.New(true)
		var st stage
		st.Logger = &logg
		test.out = st.createUnits(test.in)
		if test.out != nil {
			t.Errorf("#%d: error occured: %v", i, test.out)
		}
	}
}
