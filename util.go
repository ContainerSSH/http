package http

import (
	"io/ioutil"
	"strings"
)

func loadPem(spec string) ([]byte, error) {
	if !strings.HasPrefix(strings.TrimSpace(spec), "-----") {
		return ioutil.ReadFile(spec)
	}
	return []byte(spec), nil
}

