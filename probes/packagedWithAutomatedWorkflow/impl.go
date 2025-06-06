// Copyright 2023 OpenSSF Scorecard Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//nolint:stylecheck
package packagedWithAutomatedWorkflow

import (
	"embed"
	"fmt"

	"github.com/ossf/scorecard/v5/checker"
	"github.com/ossf/scorecard/v5/finding"
	"github.com/ossf/scorecard/v5/internal/checknames"
	"github.com/ossf/scorecard/v5/internal/probes"
	"github.com/ossf/scorecard/v5/probes/internal/utils/uerror"
)

func init() {
	probes.MustRegister(Probe, Run, []checknames.CheckName{checknames.Packaging})
}

//go:embed *.yml
var fs embed.FS

const Probe = "packagedWithAutomatedWorkflow"

func Run(raw *checker.RawResults) ([]finding.Finding, string, error) {
	if raw == nil {
		return nil, "", fmt.Errorf("%w: raw", uerror.ErrNil)
	}

	r := raw.PackagingResults
	var findings []finding.Finding
	for _, p := range r.Packages {
		if p.Msg != nil {
			continue
		}
		// Presence of a single non-debug message means the
		// check passes.
		f, err := finding.NewWith(fs, Probe,
			"Project packages its releases by way of GitHub Actions.", nil,
			finding.OutcomeTrue)
		if err != nil {
			return nil, Probe, fmt.Errorf("create finding: %w", err)
		}
		loc := &finding.Location{}
		if p.File != nil {
			loc.Path = p.File.Path
			loc.Type = p.File.Type
			loc.LineStart = &p.File.Offset
		}
		f = f.WithLocation(loc)
		findings = append(findings, *f)
	}

	if len(findings) > 0 {
		return findings, Probe, nil
	}

	f, err := finding.NewWith(fs, Probe,
		"no GitHub/GitLab publishing workflow detected.", nil,
		finding.OutcomeFalse)
	if err != nil {
		return nil, Probe, fmt.Errorf("create finding: %w", err)
	}
	return []finding.Finding{*f}, Probe, nil
}
