package envdir

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"
)

const envFileMaxSize = 1024 * 1024

var envFileNameRegex *regexp.Regexp

func init() {
	// regular expression for match env file name
	envFileNameRegex = regexp.MustCompile(`^[^0-9]{1}[^=]*$`)
}

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

		if !envFileNameRegex.MatchString(file.Name()) {
			return nil, fmt.Errorf("invalid file name: %s", file.Name())
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
	command.Env = os.Environ()
	for envName, envValue := range env {
		if envValue == "" {
			// removing environment variable
			for index, value := range command.Env {
				if strings.HasPrefix(value, envName+"=") {
					command.Env = append(command.Env[:index], command.Env[index+1:]...)
					break
				}
			}
		} else {
			command.Env = append(command.Env,
				fmt.Sprintf("%s=%s", envName, envValue))
		}
	}

	// redirect streams
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	// run command
	if err := command.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return exitError.ExitCode()
		}
		return -1
	}

	return 0
}
