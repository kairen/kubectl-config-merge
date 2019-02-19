package util

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
)

func WriteFile(path string, data []byte, perm os.FileMode) error {
	if index := strings.LastIndex(path, "/"); index != -1 {
		dir := path[:index+1]
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			if err := os.MkdirAll(dir, 0744); err != nil {
				return err
			}
		}
	}
	return ioutil.WriteFile(path, data, perm)
}

func PrettyJson(b []byte) ([]byte, error) {
	var out bytes.Buffer
	err := json.Indent(&out, b, "", "  ")
	return out.Bytes(), err
}
