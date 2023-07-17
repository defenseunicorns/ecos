package archive

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mikevanhemert/ecos/src/pkg/utils"
)

func (a *Archiver) Apply(archiveName string) error {
	if err := a.LoadArchive(archiveName); err != nil {
		return err
	}

	for componentName := range a.config.Spec.Components {
		fmt.Printf("\nCOMPONENT %s\n\n", strings.ToUpper(componentName))

		componentDir := filepath.Join(a.config.TempPaths.Base, "components", componentName)
		originalDir, _ := os.Getwd()
		envVars := []string{}

		// terraform init
		if err := os.Chdir(componentDir); err != nil {
			return fmt.Errorf("Unable to access directory '%s': %w", componentName, err)
		}

		if err := utils.ExecCommand("terraform", envVars, "init", "-plugin-dir", "providers", "-get=false"); err != nil {
			return fmt.Errorf("Unable to initialize Terraform: %w", err)
		}

		_ = os.RemoveAll(filepath.Join(componentDir, "providers"))

		// terraform apply
		envVars = a.HandleVariables(componentName)

		if err := utils.ExecCommand("terraform", envVars, "apply", "-auto-approve", "-input=false"); err != nil {
			return fmt.Errorf("Unable to apply Terraform: %w", err)
		}

		// TODO extract transitives (spec.component[name].transitives[name]) and save in variables map (a.config.PackageVariables.VariableMap)
		// 1. Extract outputs spec.components[].transitives[] to a.config.PackageVariables.VariableMap

		// TODO write out ecos Template files to originalDir/out/[component]/templates/
		// 1. Extract outputs spec.components[].templates[].variables[] to map[string]interface{}{}
		// ... tfname -> name
		// 2. Read the template file as bytes[] from spec.components[].templates.template
		// 3. Flesh out the template and write the file
		// ... tmpl, _ := template.New(/1. name/).Parse(/2. template bytes[]/)
		// ... tmpl.Execute(/file writer/, /1. map/)

		os.Chdir(originalDir)
	}

	// Write the updated ecos.yaml
	if err := utils.WriteYaml(filepath.Join(a.config.TempPaths.Base, "ecos.yaml"), a.config.Spec, 0644); err != nil {
		return fmt.Errorf("Unable to save updated ecos.yaml spec: %w", err)
	}

	utcTime := time.Now().UTC()
	updatedArchive := "ecos-" + a.config.Spec.Metadata.Name + "-" + a.config.Spec.Metadata.Version + "-" + strings.ReplaceAll(utcTime.Format(time.RFC3339Nano), ":", ".") + ".tar.zst"

	return utils.Archive(a.config.TempPaths.Base, updatedArchive)
}
