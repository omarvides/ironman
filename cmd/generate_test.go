package cmd

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"

	"github.com/ironman-project/ironman/pkg/ironman"

	testhelpers "github.com/ironman-project/ironman/cmd/testing"
	"github.com/ironman-project/ironman/pkg/testutils"
)

//Pre-installs a template for running tests
func setUpGenerateCmd(t *testing.T, client *ironman.Ironman, testCase testhelpers.CmdTestCase) {
	installCmd := newInstallCommand(client, ioutil.Discard)
	//equivalent to "ironman install https://github.com/ironman-project/template-example.git"
	args := []string{"https://github.com/ironman-project/template-example.git"}
	testhelpers.RunTestCmd(installCmd, t, args, nil)
}

func TestGenerateCmd(t *testing.T) {
	tempGenerateDir := testutils.CreateTempDir("temp-generate", t)
	defer func() {
		_ = os.RemoveAll(tempGenerateDir)
	}()
	tests := []testhelpers.CmdTestCase{
		{
			"successful generate",
			[]string{"template-example", filepath.Join(tempGenerateDir, "test-gen")},
			[]string{""},
			"Running template generator app\nDone\n",
			false,
		},
		{
			"successful generate with parameters",
			[]string{"template-example", filepath.Join(tempGenerateDir, "test-gen-with-parameters")},
			[]string{"--set", "key=value"},
			"Running template generator app\nDone\n",
			false,
		},
		{
			"template id required",
			[]string{},
			[]string{},
			"",
			true,
		},
	}
	testhelpers.RunCmdTests(t, tests, func(client *ironman.Ironman, out io.Writer) *cobra.Command {
		return newGenerateCommand(client, out)
	}, setUpGenerateCmd, nil)

}
