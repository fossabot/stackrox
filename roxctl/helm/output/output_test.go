package output

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stackrox/rox/pkg/errox"
	"github.com/stackrox/rox/roxctl/common/environment"
	"github.com/stackrox/rox/roxctl/common/printer"
	"github.com/stackrox/rox/roxctl/helm/internal/common"
	"github.com/stretchr/testify/suite"
)

func TestHelmOutputCommand(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(helmOutputTestSuite))
}

type helmOutputTestSuite struct {
	suite.Suite
	helmOutputCommand helmOutputCommand
	errOur            *bytes.Buffer
}

func (suite *helmOutputTestSuite) SetupTest() {
	testIO, _, _, errOut := environment.TestIO()
	suite.helmOutputCommand = helmOutputCommand{}
	suite.helmOutputCommand.env = environment.NewCLIEnvironment(testIO, printer.DefaultColorPrinter())
	suite.errOur = errOut
}

func (suite *helmOutputTestSuite) TestInvalidCommandArgs() {
	cases := map[string]struct {
		args       []string
		shouldFail bool
		errStdout  string
	}{
		"should not return an error if valid number of arguments given with a correct chartName": {
			args: []string{common.ChartCentralServices},
		},
		"should return an error if no arguments given": {
			args:       []string{},
			shouldFail: true,
			errStdout:  "Error: accepts 1 arg(s), received 0\n",
		},
		"should return an error if too many arguments given": {
			args:       []string{"test_chartName", "another_arg"},
			shouldFail: true,
			errStdout:  "Error: accepts 1 arg(s), received 2\n",
		},
		"should return an error if invalid chart name given": {
			args:       []string{"invalid_chartName"},
			shouldFail: true,
			errStdout: `Error: invalid argument "invalid_chartName" for "output"
`,
		},
	}

	for name, c := range cases {
		suite.Run(name, func() {
			helmCmd := suite.helmOutputCommand
			cmd := Command(helmCmd.env)

			cmd.SetArgs(c.args)
			// Ignore an executing flow of the command
			cmd.RunE = func(cmd *cobra.Command, args []string) error {
				return nil
			}
			// Redirect stdErr
			errOut := bytes.NewBufferString("")
			cmd.SetErr(errOut)

			err := cmd.Execute()
			if c.shouldFail {
				suite.Require().Error(err)
				suite.Assert().Equal(c.errStdout, errOut.String())
			} else {
				suite.Require().NoError(err)
			}
		})
	}
}

func (suite *helmOutputTestSuite) TestConstruct() {
	cmd := &cobra.Command{Use: "test"}
	chartName := "test_chartName"

	helmOutputCmd := suite.helmOutputCommand
	helmOutputCmd.Construct(chartName, cmd)
	suite.Assert().Equal(chartName, helmOutputCmd.chartName)
}

func (suite *helmOutputTestSuite) TestValidate() {
	cases := map[string]struct {
		chartName       string
		outputDir       string
		createOutputDir bool
		removeOutputDir bool
		errOutRegexp    string
		shouldFail      bool
		error           error
		errorRegexp     string
	}{
		"should not fail for valid chartName and provided outputDir": {
			chartName: common.ChartCentralServices,
			outputDir: "test_output_dir",
		},
		"should not fail for valid chartName and non provided outputDir": {
			chartName: common.ChartCentralServices,
			errOutRegexp: `WARN:	No output directory specified, using default directory "./stackrox-central-services-chart"`,
		},
		"should not fail for valid chartName and existed outputDir": {
			chartName:       common.ChartCentralServices,
			createOutputDir: true,
			removeOutputDir: true,
			errOutRegexp: "WARN:	Removed output directory .*",
		},
		"should fail for already existed output directory": {
			chartName:       common.ChartCentralServices,
			createOutputDir: true,
			removeOutputDir: false,
			errOutRegexp: "ERROR:	Directory .* already exists, use --remove or select a different directory with --output-dir.",
			shouldFail:  true,
			error:       errox.AlreadyExists,
			errorRegexp: `directory ".*" already exists`,
		},
	}

	for name, c := range cases {
		suite.Run(name, func() {
			helmOutputCmd := suite.helmOutputCommand
			helmOutputCmd.chartName = c.chartName
			helmOutputCmd.removeOutputDir = c.removeOutputDir
			helmOutputCmd.outputDir = c.outputDir
			if c.createOutputDir {
				helmOutputCmd.outputDir = suite.T().TempDir()
			}

			err := helmOutputCmd.Validate()
			suite.Assert().Regexp(c.errOutRegexp, suite.errOur.String())
			if c.shouldFail {
				suite.Require().Error(err)
				suite.Assert().ErrorIs(err, c.error)
				suite.Assert().Regexp(c.errorRegexp, err.Error())
			} else {
				suite.Require().NoError(err)
			}
		})
	}
}
