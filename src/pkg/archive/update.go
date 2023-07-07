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

type Update struct {
	originalArchive string
	updatedArchive  string
	originalConfig  *types.EcosConfig
	updatedConfig   *types.EcosConfig
}

func NewUpdate(originalArchive string, updatedArhive string) (*Update, error) {
	var (
		err    error
		update = &Update{
			originalConfig:  &types.EcosConfig{},
			updatedConfig:   &types.EcosConfig{},
			originalArchive: originalArchive,
			updatedArchive:  updatedArhive,
		}
	)

	if update.originalConfig.TempPaths.Base, err = utils.MakeTempDir(""); err != nil {
		return nil, fmt.Errorf("Unable to create temp path: %w", err)
	}

	if update.updatedConfig.TempPaths.Base, err = utils.MakeTempDir(""); err != nil {
		return nil, fmt.Errorf("Unable to create temp path: %w", err)
	}

	return update, nil
}

func NewOrDieUpdate(originalArchive string, updatedArchive string) *Update {
	var (
		err    error
		update *Update
	)

	if update, err = NewUpdate(originalArchive, updatedArchive); err != nil {
		fmt.Printf("Unable to prepare for update: %s", err)
		os.Exit(1)
	}

	return update
}

func (u *Update) ClearTempPaths() {
	_ = os.RemoveAll(u.originalConfig.TempPaths.Base)
	_ = os.RemoveAll(u.updatedConfig.TempPaths.Base)
}

func (u *Update) Update() error {
	var err error

	// Unarchive the original archive and read spec
	if u.originalConfig.Spec, err = utils.Unarchive(u.originalArchive, u.originalConfig.TempPaths.Base); err != nil {
		return err
	}

	// Unarchive the updated archive and read spec
	if u.updatedConfig.Spec, err = utils.Unarchive(u.updatedArchive, u.updatedConfig.TempPaths.Base); err != nil {
		return err
	}

	for _, updatedComponent := range u.updatedConfig.Spec.Components {
		fmt.Printf("\nCOMPONENT %s\n\n", strings.ToUpper(updatedComponent.Name))

		originalDir, _ := os.Getwd()

		// Copy the state and lock file from the original component
		if err := utils.Copy(filepath.Join(u.originalConfig.TempPaths.Base, "components", updatedComponent.Name, "terraform.tfstate"), filepath.Join(u.updatedConfig.TempPaths.Base, "components", updatedComponent.Name, "terraform.tfstate")); err != nil {
			return fmt.Errorf("Unable to extract Terraform state from original package: %w", err)
		}

		if err := utils.Copy(filepath.Join(u.originalConfig.TempPaths.Base, "components", updatedComponent.Name, ".terraform.lock.hcl"), filepath.Join(u.updatedConfig.TempPaths.Base, "components", updatedComponent.Name, ".terraform.lock.hcl")); err != nil {
			return fmt.Errorf("Unable to extract Terraform lock file from original package: %w", err)
		}

		if err := utils.Copy(filepath.Join(u.originalConfig.TempPaths.Base, "components", updatedComponent.Name, ".terraform"), filepath.Join(u.updatedConfig.TempPaths.Base, "components", updatedComponent.Name, ".terraform")); err != nil {
			return fmt.Errorf("Unable to extract Terraform providers and modules from original package: %w", err)
		}

		// terraform init
		if err := os.Chdir(filepath.Join(u.updatedConfig.TempPaths.Base, "components", updatedComponent.Name)); err != nil {
			return fmt.Errorf("Unable to access directory '%s': %w", updatedComponent.Name, err)
		}

		if err := utils.ExecCommand("terraform", "init", "-plugin-dir", "providers", "-get=false"); err != nil {
			return fmt.Errorf("Unable to initialize Terraform: %w", err)
		}

		_ = os.RemoveAll(filepath.Join(u.updatedConfig.TempPaths.Base, "components", updatedComponent.Name, "providers"))

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
	archiveName := "ecos-" + u.updatedConfig.Spec.Metadata.Name + "-" + u.updatedConfig.Spec.Metadata.Version + "-" + utcTime.Format(time.RFC3339Nano) + ".tar.zst"

	return utils.Archive(u.updatedConfig.TempPaths.Base, archiveName)
}
