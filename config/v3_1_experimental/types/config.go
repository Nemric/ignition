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

package types

import (
	"github.com/coreos/go-semver/semver"

	"github.com/coreos/ignition/v2/config/validate/report"
)

var (
	MaxVersion = semver.Version{
		Major:      3,
		Minor:      1,
		PreRelease: "experimental",
	}
)

func (c Config) Validate() report.Report {
	r := report.Report{}
	rules := []rule{}

	for _, rule := range rules {
		rule(c, &r)
	}
	return r
}

type rule func(cfg Config, report *report.Report)