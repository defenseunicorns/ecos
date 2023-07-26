package archive

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mholt/archiver/v3"
	"github.com/mikevanhemert/ecos/src/pkg/utils"
	"github.com/mikevanhemert/ecos/src/types"
)

type Archiver struct {
	config    *types.EcosConfig
	oldConfig *types.EcosConfig // Used when updating
}

func New(config *types.EcosConfig) *Archiver {
	var (
		archiver = &Archiver{
			config:    config,
			oldConfig: &types.EcosConfig{},
		}
	)

	if archiver.config.PackageVariables == nil {
		archiver.config.PackageVariables = map[string]string{}
	}

	return archiver
}

func (a *Archiver) ClearTempPaths() {
	_ = os.RemoveAll(a.config.TempPaths.Base)
	_ = os.RemoveAll(a.oldConfig.TempPaths.Base)
}

func (a *Archiver) CopyOverrides(componentName string) {
	var (
		oldComponent  types.EcosComponent
		oldProperties types.EcosVariable
		ok            bool
	)

	for name, properties := range a.config.Spec.Components[componentName].Variables {
		if oldComponent, ok = a.oldConfig.Spec.Components[componentName]; ok != true {
			continue
		}
		if oldProperties, ok = oldComponent.Variables[name]; ok != true {
			continue
		}

		if len(oldProperties.Override) > 0 {
			properties.Override = oldProperties.Override
			a.config.Spec.Components[componentName].Variables[name] = properties
		}
	}
}

func (a *Archiver) HandleVariables(componentName string) []string {
	var envVars = []string{}

	// Iterate over the variables
	for name, properties := range a.config.Spec.Components[componentName].Variables {
		value := properties.Default

		if len(properties.Override) != 0 {
			value = properties.Override
		}

		// Overwrite default
		if provided, ok := a.config.PackageVariables[name]; ok {
			value = provided
			properties.Override = provided
			a.config.Spec.Components[componentName].Variables[name] = properties
		}

		if len(value) == 0 {
			fmt.Printf("Warning: variable %s's value is undefined and will be skipped\n", name)
			continue
		}

		envVars = append(envVars, fmt.Sprintf("TF_VAR_%s=%s", name, value))
	}

	return envVars
}

func (a *Archiver) LoadArchive(archiveName string) error {
	var err error

	// 1. create tempdir
	if a.config.TempPaths.Base, err = utils.MakeTempDir(""); err != nil {
		return fmt.Errorf("Unable to create temp path to unarchive %s into: %w", archiveName, err)
	}

	// 2. load the archive
	if err = archiver.Unarchive(archiveName, a.config.TempPaths.Base); err != nil {
		return fmt.Errorf("Unable to unarchive %s: %w", archiveName, err)
	}

	// 3. load spec
	if err := utils.ReadYaml(filepath.Join(a.config.TempPaths.Base, "ecos.yaml"), &a.config.Spec); err != nil {
		return fmt.Errorf("Unable to read Ecos spec: %s", err)
	}

	return err
}

func (a *Archiver) LoadOldArchive(archiveName string) error {
	var err error

	// 1. create tempdir
	if a.oldConfig.TempPaths.Base, err = utils.MakeTempDir(""); err != nil {
		return fmt.Errorf("Unable to create temp path to unarchive old archive %s into: %w", archiveName, err)
	}

	// 2. load the old archive
	if err = archiver.Unarchive(archiveName, a.oldConfig.TempPaths.Base); err != nil {
		return fmt.Errorf("Unable to unarchive old archive %s: %w", archiveName, err)
	}

	// 3. load old spec
	if err := utils.ReadYaml(filepath.Join(a.oldConfig.TempPaths.Base, "ecos.yaml"), &a.oldConfig.Spec); err != nil {
		return fmt.Errorf("Unable to read old Ecos spec: %s", err)
	}

	return err
}
