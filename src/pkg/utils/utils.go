package utils

import (
	"os"

	goyaml "github.com/goccy/go-yaml"
	"github.com/mikevanhemert/ecos/src/pkg/message"
)

func ReadYaml(path string, config any) error {
	message.Debugf("Loading ecos config %s", path)

	file, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return goyaml.Unmarshal(file, config)
}
