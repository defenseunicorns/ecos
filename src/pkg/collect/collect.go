package collect

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mholt/archiver/v3"
	"github.com/mikevanhemert/ecos/src/pkg/utils"
	"github.com/mikevanhemert/ecos/src/types"
)

type Collect struct {
	config *types.EcosConfig
}

func New(config *types.EcosConfig) (*Collect, error) {
	if config == nil {
		return nil, fmt.Errorf("Ecos spec not provided")
	}

	var (
		err       error
		pkgConfig = &Collect{
			config: config,
		}
	)

	// Create a temp directory for the package
	if pkgConfig.config.TempPaths, err = createPaths(); err != nil {
		return nil, fmt.Errorf("unable to create package temp paths: %w", err)
	}

	return pkgConfig, nil
}

func NewOrDie(config *types.EcosConfig) *Collect {
	var (
		err       error
		pkgConfig *Collect
	)

	if pkgConfig, err = New(config); err != nil {
		fmt.Printf("Unable to setup the package config: %s", err)
		os.Exit(1)
	}

	return pkgConfig
}

func archive(sourceDir string, destinationTarball string) error {
	archiveSrc := []string{sourceDir + string(os.PathSeparator)}

	if err := archiver.Archive(archiveSrc, destinationTarball); err != nil {
		return fmt.Errorf("Unable to create archive: %w", err)
	}

	_, err := os.Stat(destinationTarball)
	if err != nil {
		return fmt.Errorf("Unable to read the package archive: %w", err)
	}

	return nil
}

func createPaths() (paths types.TempPaths, err error) {
	basePath, err := utils.MakeTempDir("")
	paths = types.TempPaths{
		Base: basePath,

		Components: filepath.Join(basePath, "components"),
	}

	return paths, err
}

func (c *Collect) ClearTempPaths() {
	_ = os.RemoveAll(c.config.TempPaths.Base)
}

func (c *Collect) Collect() error {
	// Read in the Ecos Spec
	if err := utils.ReadYaml("ecos.yaml", &c.config.Spec); err != nil {
		return fmt.Errorf("Unable to read the ecos.yaml file: %s", err.Error())
	}

	// Write out the Ecos Spec
	if err := utils.WriteYaml(filepath.Join(c.config.TempPaths.Base, "ecos.yaml"), c.config.Spec, 0766); err != nil {
		return fmt.Errorf("Unable to write the ecos.yaml file: %s", err.Error())
	}

	for _, component := range c.config.Spec.Components {
		fmt.Printf("\nCOMPONENT %s\n\ns", strings.ToUpper(component.Name))

		// copy dir to temp location
		if err := utils.Copy(component.Name, filepath.Join(c.config.TempPaths.Components, component.Name)); err != nil {
			return fmt.Errorf("Unable to copy component %s: %s", component.Name, err.Error())
		}

		// Providers: terraform providers mirror
		originalDir, _ := os.Getwd()
		if err := os.Chdir(filepath.Join(c.config.TempPaths.Components, component.Name)); err != nil {
			return fmt.Errorf("Unable to access directory '%s': %w", component.Name, err)
		}

		if err := utils.ExecCommand("terraform", "providers", "mirror", "providers"); err != nil {
			return fmt.Errorf("Unable to mirror Terraform providers: %s", err.Error())
		}

		// Modules: terraform get
		if err := utils.ExecCommand("terraform", "get"); err != nil {
			return fmt.Errorf("Unable to get terraform modules: %s", err.Error())
		}
		os.Chdir(originalDir)
	}

	archiveName := "ecos-" + c.config.Spec.Metadata.Name + "-" + c.config.Spec.Metadata.Version + ".tar.zst"
	fmt.Printf("Creating archive %s\n", archiveName)

	if err := archive(c.config.TempPaths.Base, archiveName); err != nil {
		return err
	}

	return nil
}
