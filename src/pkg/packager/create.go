package packager

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mikevanhemert/ecos/src/pkg/message"
	"github.com/mikevanhemert/ecos/src/pkg/utils"
)

func (p *Packager) Create(baseDir string) error {

	if err := utils.ReadYaml(filepath.Join(baseDir, "ecos.yaml"), p.config); err != nil {
		return fmt.Errorf("Unable to read the ecos.yaml file: %s", err.Error())
	}

	// Change working directory to match baseDir
	if baseDir != "" {
		if err := os.Chdir(baseDir); err != nil {
			return fmt.Errorf("Unable to access directory %s: %s", baseDir, err.Error())
		}
		message.Note("Using build directory %s", baseDir)
	}
	return nil
}
