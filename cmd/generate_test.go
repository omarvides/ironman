package cmd

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/spf13/cobra"

	"github.com/ironman-project/ironman/pkg/ironman"

	"github.com/ironman-project/ironman/pkg/testutils"
)

type cmdTestCase struct {
	name     string
	args     []string
	flags    []string
	expected string
	err      bool
}

type testCmdFactory func(ironman *ironman.Ironman, out io.Writer) *cobra.Command

//Pre-installs a template for running tests
func setUpGenerateCmd(t *testing.T, client *ironman.Ironman, testCase cmdTestCase) {
	installCmd := newInstallCommand(client, ioutil.Discard)
	//equivalente to ironman install https://github.com/ironman-project/template-example.git
	args := []string{"https://github.com/ironman-project/template-example.git"}
	runTestCmd(installCmd, t, args, nil)
}

func TestGenerateCmd(t *testing.T) {
	tempGenerateDir := testutils.CreateTempDir("temp-generate", t)
	defer func() {
		_ = os.RemoveAll(tempGenerateDir)
	}()
	tests := []cmdTestCase{
		{
			"successful generate",
			[]string{"template-example", filepath.Join(tempGenerateDir, "test-gen")},
			[]string{""},
			"Running template generator app\nDone\n",
			false,
		},
	}
	runCmdTests(t, tests, func(client *ironman.Ironman, out io.Writer) *cobra.Command {
		return newGenerateCommand(client, out)
	}, setUpGenerateCmd, nil)

}

type cmdTestCaseSetUpTearDown func(*testing.T, *ironman.Ironman, cmdTestCase)

func runCmdTests(t *testing.T, tests []cmdTestCase, cmdFactory testCmdFactory, setUp cmdTestCaseSetUpTearDown, tearDown cmdTestCaseSetUpTearDown) {
	var buf bytes.Buffer
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempHome := testutils.CreateTempDir("ihome", t)
			client := ironman.New(tempHome)
			defer func() {
				_ = os.RemoveAll(tempHome)
			}()

			if setUp != nil {
				setUp(t, client, tt)
			}
			if tearDown != nil {
				defer tearDown(t, client, tt)
			}

			cmd := cmdFactory(client, &buf)
			err := runTestCmd(cmd, t, tt.args, tt.flags)
			cmd.ParseFlags(tt.flags)

			if (err != nil) != tt.err {
				t.Errorf("expected error, got '%v'", err)
			}
			re := regexp.MustCompile(tt.expected)
			if !re.Match(buf.Bytes()) {
				t.Errorf("expected\n%q\ngot\n%q", tt.expected, buf.String())
			}
			buf.Reset()
		})
	}
}

func runTestCmd(cmd *cobra.Command, t *testing.T, args []string, flags []string) error {
	cmd.ParseFlags(flags)
	err := cmd.RunE(cmd, args)
	return err
}
