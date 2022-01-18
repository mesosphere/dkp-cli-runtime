package symlinkexechandler_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"
)

type HandlerTestSuite struct {
	suite.Suite
	tempDir            string
	testExecutableName string
	beforePATH         string
}

func (suite *HandlerTestSuite) SetupSuite() {
	suite.tempDir = suite.T().TempDir()
	suite.testExecutableName = "simple"
	cmd := exec.Command( //nolint:gosec // Building test binary into temporary directory.
		"go", "build",
		"-o", filepath.Join(suite.tempDir, suite.testExecutableName),
		"testdata/simple_main.go")
	output, err := cmd.CombinedOutput()
	suite.Require().NoError(err)
	suite.T().Log(string(output))
	suite.beforePATH = os.Getenv("PATH")
	os.Setenv("PATH", suite.tempDir)
}

func (suite *HandlerTestSuite) TearDownSuite() {
	os.Setenv("PATH", suite.beforePATH)
}

func (suite *HandlerTestSuite) TestDirectInvocation() {
	cmd := exec.Command( //nolint:gosec // Building test binary into temporary directory.
		suite.testExecutableName,
	)
	output, err := cmd.CombinedOutput()
	suite.Require().NoError(err)
	suite.Assert().Equal(`"simple"`, string(output))
}

func (suite *HandlerTestSuite) TestDirectInvocationWithFlags() {
	cmd := exec.Command( //nolint:gosec // Building test binary into temporary directory.
		suite.testExecutableName,
		"--flag1",
		"--flag2",
		"with spaces",
	)
	output, err := cmd.CombinedOutput()
	suite.Require().NoError(err)
	suite.Assert().Equal(`"simple" "--flag1" "--flag2" "with spaces"`, string(output))
}

func (suite *HandlerTestSuite) TestSymlinkInvocationWithFlags() {
	symlinkTempDir := suite.T().TempDir()
	symlinkPath := filepath.Join(symlinkTempDir, "simple_symlink")
	os.Symlink(filepath.Join(suite.tempDir, suite.testExecutableName), symlinkPath)
	cmd := exec.Command(
		symlinkPath,
		"--flag1",
		"--flag2",
		"with spaces",
	)
	output, err := cmd.CombinedOutput()
	suite.Require().NoError(err)
	suite.Assert().Equal(
		fmt.Sprintf(`%q "--flag1" "--flag2" "with spaces"`, symlinkPath),
		string(output),
	)
}

func (suite *HandlerTestSuite) TestSymlinkInvocationWithSubcommands() {
	symlinkTempDir := suite.T().TempDir()
	symlinkPath := filepath.Join(symlinkTempDir, "simple-subcommand")
	os.Symlink(filepath.Join(suite.tempDir, suite.testExecutableName), symlinkPath)
	cmd := exec.Command(
		symlinkPath,
		"--flag1",
		"--flag2",
		"with spaces",
	)
	output, err := cmd.CombinedOutput()
	suite.Require().NoError(err)
	suite.Assert().Equal(
		`"simple" "subcommand" "--flag1" "--flag2" "with spaces"`,
		string(output),
	)
}

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}
