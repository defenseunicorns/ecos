package archive

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mikevanhemert/ecos/src/pkg/utils"
)

func (a *Archiver) Destroy(archiveName string) error {
	if err := a.LoadArchive(archiveName); err != nil {
		return err
	}

	for componentName := range a.config.Spec.Components {
		fmt.Printf("\nCOMPONENT %s\n\n", strings.ToUpper(componentName))

		originalDir, _ := os.Getwd()
		envVars := []string{}

		// terraform destoy
		if err := os.Chdir(filepath.Join(a.config.TempPaths.Base, "components", componentName)); err != nil {
			return fmt.Errorf("Unable to access directory '%s': %s", componentName, err)
		}

		envVars = a.HandleVariables(componentName)

		if err := utils.ExecCommand("terraform", envVars, "destroy", "-auto-approve"); err != nil {
			return fmt.Errorf("Unable to destroy Terraform: %w", err)
		}

		os.Chdir(originalDir)
	}

	// Write the updated ecos.yaml
	if err := utils.WriteYaml(filepath.Join(a.config.TempPaths.Base, "ecos.yaml"), a.config.Spec, 0644); err != nil {
		return fmt.Errorf("Unable to save updated ecos.yaml spec: %w", err)
	}

	utcTime := time.Now().UTC()
	updatedArchvie := "ecos-" + a.config.Spec.Metadata.Name + "-" + a.config.Spec.Metadata.Version + "-" + strings.ReplaceAll(utcTime.Format(time.RFC3339Nano), ":", ".") + ".tar.zst"

	return utils.Archive(a.config.TempPaths.Base, updatedArchvie)
}
