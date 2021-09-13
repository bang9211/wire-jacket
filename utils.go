package wirejacket

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-ps"
)

// GetOrSetHomePath returns the existing value for env 'OSSICONES_SRC_HOME' if present.
// Otherwise, it sets and returns the given value.
func GetOrSetHomePath() (string, error) {
	homePath := os.Getenv("OSSICONES_SRC_HOME")
	if homePath == "" {
		cmdOut, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
		if err != nil {
			log.Fatal(err.Error(), string(cmdOut))
			return "", err
		}
		homePath = strings.TrimSpace(string(cmdOut))
		os.Setenv("OSSICONES_SRC_HOME", homePath)
	}
	return homePath, nil
}

// IsRunning checks if the same process is already running.
// It returns error, if failed to get list of processes because of the unsupported OS.
func IsRunning() (bool, error) {
	processList, err := ps.Processes()
	if err != nil {
		return false, err
	}

	if checkProcessIsAlreadyRunningByName(os.Args[0], processList) != nil {
		return true, nil
	}

	return false, nil
}

func checkProcessIsAlreadyRunningByName(processName string, processList []ps.Process) ps.Process {
	for _, process := range processList {
		if processName == process.Executable() && os.Getpid() != process.Pid() {
			return process
		}
	}

	return nil
}

func HandleError(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func GetFileDir(path string) string {
	return filepath.Dir(path)
}

func GetFileName(path string) string {
	return filepath.Base(path)
}

func GetFileNameWithoutExtension(path string) string {
	file := filepath.Base(path)
	splited := strings.Split(file, ".")
	return splited[0]
}

func GetFileExtension(path string) string {
	ext := filepath.Ext(path)
	if len(ext) > 0 {
		return ext[1:]
	}
	return ext
}

func IsContain(list []string, key string) bool {
	for _, s := range list {
		if s == key {
			return true
		}
	}
	return false
}

func RemoveElement(slice []string, s string) []string {
	index := -1
	for i, k := range slice {
		if k == s {
			index = i
		}
	}
	if index == -1 {
		return slice
	}
	return append(slice[:index], slice[index+1:]...)
}
