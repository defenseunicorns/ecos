package archive

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mikevanhemert/ecos/src/pkg/utils"
	"github.com/mikevanhemert/ecos/src/types"
)

type Apply struct {
	archiveName string
	config      *types.EcosConfig
}

func NewApply(archiveName string) (*Apply, error) {
	// Create a temp directory
	var (
		err   error
		apply = &Apply{
			config:      &types.EcosConfig{},
			archiveName: archiveName,
		}
	)

	if apply.config.TempPaths.Base, err = utils.MakeTempDir(""); err != nil {
		return nil, fmt.Errorf("Unable to create temp path: %w", err)
	}

	return apply, nil
}

func NewOrDieApply(archiveName string) *Apply {
	var (
		err   error
		apply *Apply
	)

	if apply, err = NewApply(archiveName); err != nil {
		fmt.Printf("Unable to prepare for apply: %s", err)
		os.Exit(1)
	}

	return apply
}

func (a *Apply) Apply() error {
	var err error

	// Unarchive to the temp dir and read spec
	if a.config.Spec, err = utils.Unarchive(a.archiveName, a.config.TempPaths.Base); err != nil {
		return err
	}

	for _, component := range a.config.Spec.Components {
		fmt.Printf("\nCOMPONENT %s\n\n", strings.ToUpper(component.Name))

		originalDir, _ := os.Getwd()

		// terraform init
		if err := os.Chdir(filepath.Join(a.config.TempPaths.Base, "components", component.Name)); err != nil {
			return fmt.Errorf("Unable to access directory '%s': %w", component.Name, err)
		}

		if err := utils.ExecCommand("terraform", "init", "-plugin-dir", "providers", "-get=false"); err != nil {
			return fmt.Errorf("Unable to initialize Terraform: %w", err)
		}

		_ = os.RemoveAll(filepath.Join(a.config.TempPaths.Base, "components", component.Name, "providers"))

		// terraform apply
		// TODO use variables defined in spec.component[].variables
		if err := utils.ExecCommand("terraform", "apply", "-auto-approve", "-input=false"); err != nil {
			return fmt.Errorf("Unable to apply Terraform: %w", err)
		}

		// TODO extract transitives (spec.component[].transitives) and save in variables map (a.config.PackageVariables.VariableMap)

		// TODO write out ecos Template files to originalDir/out/[component]/templates/

		os.Chdir(originalDir)
	}

	utcTime := time.Now().UTC()
	archiveName := "ecos-" + a.config.Spec.Metadata.Name + "-" + a.config.Spec.Metadata.Version + "-" + utcTime.Format(time.RFC3339Nano) + ".tar.zst"

	return utils.Archive(a.config.TempPaths.Base, archiveName)
}

func (a *Apply) ClearTempPaths() {
	_ = os.RemoveAll(a.config.TempPaths.Base)
}
