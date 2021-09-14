package wirejacket

import (
	"path/filepath"
	"strings"
)

func getFileDir(path string) string {
	return filepath.Dir(path)
}

func getFileName(path string) string {
	return filepath.Base(path)
}

func getFileNameWithoutExtension(path string) string {
	file := filepath.Base(path)
	splited := strings.Split(file, ".")
	return splited[0]
}

func getFileExtension(path string) string {
	ext := filepath.Ext(path)
	if len(ext) > 0 {
		return ext[1:]
	}
	return ext
}

func isContain(list []string, key string) bool {
	for _, s := range list {
		if s == key {
			return true
		}
	}
	return false
}

func removeElement(slice []string, s string) []string {
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
