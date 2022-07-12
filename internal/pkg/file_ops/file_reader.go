package file_ops

import (
	"io/ioutil"
)

func Read(filepath string) (string, error) {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
