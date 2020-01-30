package envdir

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"syscall"
)

const envFileMaxSize = 1024 * 1024

// ReadDir scan catalog and return env list
func ReadDir(dir string) (map[string]string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("could not read files in directory: %s", err)
	}

	envMap := map[string]string{}

	for _, file := range files {

		if file.IsDir() {
			return nil, fmt.Errorf("file %s is directory", file.Name())
		}

		if !file.Mode().IsRegular() {
			return nil, fmt.Errorf("file %s is not regular file", file.Name())
		}

		if file.Size() > envFileMaxSize {
			return nil, fmt.Errorf("file size of %s is greater than %d bytes",
				file.Name(), envFileMaxSize)
		}

		if strings.Contains(file.Name(), "=") {
			return nil, fmt.Errorf("invalid file name (contains '='): %s", file.Name())
		}

		for i := 0; i < 10; i++ {
			if strings.HasPrefix(file.Name(), strconv.Itoa(i)) {
				return nil, fmt.Errorf("invalid file name (contains '%d'): %s", i, file.Name())
			}
		}

		fileContent, err := ioutil.ReadFile(path.Join(dir, file.Name()))
		if err != nil {
			return nil, fmt.Errorf("could not read file: %s, cause: %s",
				file.Name(), err)
		}

		envMap[file.Name()] = strings.Trim(string(fileContent), "\r\n\t ")

	}

	return envMap, nil
}

// RunCmd run command with arguments with enviroment list
func RunCmd(cmd []string, env map[string]string) int {
	cmdArgs := []string{}
	if len(cmd) > 1 {
		cmdArgs = cmd[1:]
	}

	command := exec.Command(cmd[0], cmdArgs...)

	// make env list
	for envName, envValue := range env {
		if envValue == "" {
			os.Unsetenv(envName)
		} else {
			os.Setenv(envName, envValue)
		}
	}
	command.Env = os.Environ()

	// redirect streams
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	// run command
	if err := command.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if status, ok := exitError.Sys().(syscall.WaitStatus); ok {
				return status.ExitStatus()
			}
			return -1
		}
		return -2
	}

	return 0
}
