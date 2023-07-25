package archive

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mikevanhemert/ecos/src/pkg/utils"
)

func (a *Archiver) Collect() error {
	var (
		err     error
		envVars []string
	)

	// Create a temp dir
	if a.config.TempPaths.Base, err = utils.MakeTempDir(""); err != nil {
		return fmt.Errorf("Unable to create temp path for collecting resources: %w", err)
	}

	// Read in the Ecos Spec
	if err := utils.ReadYaml("ecos.yaml", &a.config.Spec); err != nil {
		return fmt.Errorf("Unable to read the Ecos spec: %w", err)
	}

	// Write out the Ecos Spec
	if err := utils.WriteYaml(filepath.Join(a.config.TempPaths.Base, "ecos.yaml"), a.config.Spec, 0766); err != nil {
		return fmt.Errorf("Unable to write the Ecos spec: %w", err)
	}

	for componentName := range a.config.Spec.Components {
		fmt.Printf("\nCOMPONENT %s\n\n", strings.ToUpper(componentName))

		componentDir := filepath.Join(a.config.TempPaths.Base, "components", componentName)
		originalDir, _ := os.Getwd()

		// copy dir to temp location
		if err := utils.Copy(componentName, componentDir); err != nil {
			return fmt.Errorf("Unable to copy component %s: %w", componentName, err)
		}

		// Providers: terraform providers mirror
		if err := os.Chdir(componentDir); err != nil {
			return fmt.Errorf("Unable to access directory '%s': %w", componentName, err)
		}

		if _, err := utils.ExecCommand("terraform", envVars, "providers", "mirror", "providers"); err != nil {
			fmt.Printf("Warn: Unable to mirror Terraform providers: %s", err.Error())
		}

		// Modules: terraform get
		if _, err := utils.ExecCommand("terraform", envVars, "get"); err != nil {
			return fmt.Errorf("Unable to get terraform modules: %w", err)
		}

		os.Chdir(originalDir)
	}

	archiveName := "ecos-" + a.config.Spec.Metadata.Name + "-" + a.config.Spec.Metadata.Version + ".tar.zst"

	return utils.Archive(a.config.TempPaths.Base, archiveName)
}
