package envdir

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

var envDirTest = []struct {
	testName, dir string
	err           bool
	files         map[string]string
	bigFiles      []string
}{
	{
		testName: "nonexistent dir",
		dir:      "/fakedir",
		err:      true,
	},
	{
		testName: "not regular file",
		dir:      "/dev",
		err:      true,
	},
	{
		testName: "big file",
		bigFiles: []string{"TEST"},
		err:      true,
	},
	{
		testName: "env file with =",
		files: map[string]string{
			"ASD=": "",
		},
		err: true,
	},
	{
		testName: "env file start with number",
		files: map[string]string{
			"2XXX": "",
		},
		err: true,
	},
	{
		testName: "correct data",
		files: map[string]string{
			"TEST": "TEST",
			"LESS": "23",
		},
		err: false,
	},
}

var envCmdTest = []struct {
	testName string
	cmd      []string
	envMap   map[string]string
	err      bool
}{
	{
		testName: "unexistent binary",
		cmd:      []string{"/bin/ASDASSDASDASD"},
		err:      true,
	},
	{
		testName: "invalid args",
		cmd:      []string{"/bin/ls", "-ZZZ", "23"},
		err:      true,
	},
	{
		testName: "correct cmd",
		cmd:      []string{"/bin/ls"},
		err:      false,
	},
	{
		testName: "correct cmd with args",
		cmd:      []string{"/bin/ls", "-al"},
		err:      false,
	},
	{
		testName: "correct cmd with args and envs",
		cmd:      []string{"/bin/ls", "-al"},
		envMap: map[string]string{
			"TEST": "TEST",
			"PATH": "",
		},
		err: false,
	},
}

func TestEnvDir(t *testing.T) {
	for _, testData := range envDirTest {
		t.Run(testData.testName, func(t *testing.T) {

			var err error

			if len(testData.bigFiles) > 0 || len(testData.files) > 0 {

				testData.dir, err = ioutil.TempDir("/tmp", "autotest")
				require.NoError(t, err)
				defer os.RemoveAll(testData.dir)

				for _, fileName := range testData.bigFiles {
					file, err := os.Create(path.Join(testData.dir, fileName))
					require.NoError(t, err)
					err = file.Truncate(1e7)
					require.NoError(t, err)
				}

				for fileName, fileContent := range testData.files {
					err = ioutil.WriteFile(path.Join(testData.dir, fileName), []byte(fileContent), 0600)
					require.NoError(t, err)
				}

			}

			_, err = ReadDir(testData.dir)

			if testData.err {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

		})
	}
}

func TestRunCmd(t *testing.T) {
	for _, testData := range envCmdTest {
		t.Run(testData.testName, func(t *testing.T) {

			rc := RunCmd(testData.cmd, testData.envMap)

			if testData.err {
				assert.NotEqual(t, rc, 0)
			} else {
				assert.Equal(t, rc, 0)
			}

		})
	}
}
