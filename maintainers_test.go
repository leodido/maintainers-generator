package main

import (
	"testing"

	"gotest.tools/assert"
)

type maintainersEncodeTestCase struct {
	title  string
	input  Maintainers
	result string
}

var maintainersEncodeTestCases = []maintainersEncodeTestCase{
	{
		title: "a",
		input: Maintainers{
			Maintainer{
				Name:    "Leonardo Di Donato",
				GitHub:  "https://github.com/leodido",
				Company: "Sysdig",
				Projects: []string{
					"https://github.com/falcosecurity/falco",
					"https://github.com/falcosecurity/client-go",
					"https://github.com/falcosecurity/client-rs",
				},
			},
			Maintainer{
				Name:    "Another Contributor",
				GitHub:  "https://github.com/another",
				Company: "Another company",
				Projects: []string{
					"https://github.com/falcosecurity/client-go",
					"https://github.com/falcosecurity/falcosidekick",
				},
			},
		},
		result: `- name: Leonardo Di Donato
  github: https://github.com/leodido
  company: Sysdig
  projects:
  - https://github.com/falcosecurity/falco
  - https://github.com/falcosecurity/client-go
  - https://github.com/falcosecurity/client-rs
- name: Another Contributor
  github: https://github.com/another
  company: Another company
  projects:
  - https://github.com/falcosecurity/client-go
  - https://github.com/falcosecurity/falcosidekick
`,
	},
}

func TestMaintainersEncode(t *testing.T) {
	t.Helper()
	for _, tc := range maintainersEncodeTestCases {
		tc := tc
		t.Run(tc.title, func(t *testing.T) {
			res, err := tc.input.Encode()
			assert.NilError(t, err)
			assert.Equal(t, res, tc.result)
		})
	}
}
