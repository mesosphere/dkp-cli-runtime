package version_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mesosphere/dkp-cli-runtime/runtime_core/cmd/version"
)

func TestVersionCommand(t *testing.T) {
	testVersion := version.Version{
		Major:      "1",
		Minor:      "2",
		GitVersion: "v1.2",
		BuildDate:  "2021-12-13 12:52:12 UTC",
		GoVersion:  "go1.17.5",
		Compiler:   "gc",
		Platform:   "linux/amd64",
	}

	testVersion2 := version.Version{
		Major:      "2",
		Minor:      "3",
		GitVersion: "v2.3",
		BuildDate:  "2021-12-14 12:52:12 UTC",
		GoVersion:  "go1.17.7",
		Compiler:   "gcc-go",
		Platform:   "linux/arm64",
	}

	t.Run("with single version", func(t *testing.T) {
		//nolint:unparam // error is needed to be passed to NewOptions
		versionGetter := func() (version.Versions, error) {
			return version.Versions{
				"": testVersion,
			}, nil
		}

		assertOutput(t, "default", versionGetter,
			[]string{},
			`version.Version{Major:"1", Minor:"2", GitVersion:"v1.2", GitCommit:"", GitTreeState:"", `+
				`BuildDate:"2021-12-13 12:52:12 UTC", GoVersion:"go1.17.5", Compiler:"gc", Platform:"linux/amd64"}`+"\n",
		)

		assertOutput(t, "short", versionGetter,
			[]string{"--short"},
			"v1.2\n",
		)

		assertOutput(t, "JSON", versionGetter,
			[]string{"-o", "json"},
			jsonOutputSingle,
		)

		assertOutput(t, "YAML", versionGetter,
			[]string{"-o", "yaml"},
			yamlOutputSingle,
		)

		t.Run("invalid format", func(t *testing.T) {
			outBuf := bytes.Buffer{}
			versionCmd := version.NewCommandWithVersionGetter(&outBuf, versionGetter)
			require.NotNil(t, versionCmd)

			versionCmd.SetArgs([]string{"-o", "invalid"})
			err := versionCmd.Execute()
			assert.Error(t, err)
		})
	})

	t.Run("with multiple versions", func(t *testing.T) {
		versionGetter := func() (version.Versions, error) {
			return version.Versions{
				"A": testVersion,
				"B": testVersion2,
			}, nil
		}

		assertOutput(t, "default", versionGetter,
			[]string{},
			`A: version.Version{Major:"1", Minor:"2", GitVersion:"v1.2", GitCommit:"", GitTreeState:"", `+
				`BuildDate:"2021-12-13 12:52:12 UTC", GoVersion:"go1.17.5", Compiler:"gc", Platform:"linux/amd64"}`+"\n"+
				`B: version.Version{Major:"2", Minor:"3", GitVersion:"v2.3", GitCommit:"", GitTreeState:"", `+
				`BuildDate:"2021-12-14 12:52:12 UTC", GoVersion:"go1.17.7", Compiler:"gcc-go", Platform:"linux/arm64"}`+"\n",
		)

		assertOutput(t, "short", versionGetter,
			[]string{"--short"},
			"A: v1.2\nB: v2.3\n",
		)

		assertOutput(t, "JSON", versionGetter,
			[]string{"-o", "json"},
			jsonOutputMulti,
		)

		assertOutput(t, "YAML", versionGetter,
			[]string{"-o", "yaml"},
			yamlOutputMulti,
		)
	})
}

func assertOutput(
	t *testing.T, testName string,
	versionGetter func() (version.Versions, error),
	params []string, expectedOutput string,
) {
	t.Run(testName, func(t *testing.T) {
		outBuf := bytes.Buffer{}
		versionCmd := version.NewCommandWithVersionGetter(&outBuf, versionGetter)
		require.NotNil(t, versionCmd)

		versionCmd.SetArgs(params)
		err := versionCmd.Execute()
		require.NoError(t, err)
		assert.Equal(t, expectedOutput, outBuf.String())
	})
}

const jsonOutputSingle = `{
  "major": "1",
  "minor": "2",
  "gitVersion": "v1.2",
  "gitCommit": "",
  "gitTreeState": "",
  "buildDate": "2021-12-13 12:52:12 UTC",
  "goVersion": "go1.17.5",
  "compiler": "gc",
  "platform": "linux/amd64"
}
`

const yamlOutputSingle = `buildDate: 2021-12-13 12:52:12 UTC
compiler: gc
gitCommit: ""
gitTreeState: ""
gitVersion: v1.2
goVersion: go1.17.5
major: "1"
minor: "2"
platform: linux/amd64

`

const jsonOutputMulti = `{
  "A": {
    "major": "1",
    "minor": "2",
    "gitVersion": "v1.2",
    "gitCommit": "",
    "gitTreeState": "",
    "buildDate": "2021-12-13 12:52:12 UTC",
    "goVersion": "go1.17.5",
    "compiler": "gc",
    "platform": "linux/amd64"
  },
  "B": {
    "major": "2",
    "minor": "3",
    "gitVersion": "v2.3",
    "gitCommit": "",
    "gitTreeState": "",
    "buildDate": "2021-12-14 12:52:12 UTC",
    "goVersion": "go1.17.7",
    "compiler": "gcc-go",
    "platform": "linux/arm64"
  }
}
`

const yamlOutputMulti = `A:
  buildDate: 2021-12-13 12:52:12 UTC
  compiler: gc
  gitCommit: ""
  gitTreeState: ""
  gitVersion: v1.2
  goVersion: go1.17.5
  major: "1"
  minor: "2"
  platform: linux/amd64
B:
  buildDate: 2021-12-14 12:52:12 UTC
  compiler: gcc-go
  gitCommit: ""
  gitTreeState: ""
  gitVersion: v2.3
  goVersion: go1.17.7
  major: "2"
  minor: "3"
  platform: linux/arm64

`
