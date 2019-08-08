/*
 * @Author: kidd
 * @Date: 1/19/19 10:24 PM
 */

package gopath

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// SelfPath gets compiled executable file absolute path
func SelfPath() string {
	path, _ := filepath.Abs(os.Args[0])
	return path
}

// SelfDir gets compiled executable file directory
func SelfDir() string {
	return filepath.Dir(SelfPath())
}

func CurrentPath() (string, error) {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}
	path, err := filepath.Abs(file)
	if err != nil {
		return "", err
	}
	i := strings.LastIndex(path, "/")
	if i < 0 {
		i = strings.LastIndex(path, "\\")
	}
	if i < 0 {
		return "", errors.New(`error: Can't find "/" or "\".`)
	}
	return string(path[0:i]), nil
}

func CreatePath(path string) error {
	exists := Exists(path)
	if !exists {
		if err := os.Mkdir(path, os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}

// FileExists reports whether the named file or directory exists.
func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func RemovePath(path string) error {
	return os.RemoveAll(path)
}
