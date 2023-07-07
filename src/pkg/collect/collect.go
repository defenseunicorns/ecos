package collect

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
		err     error
		collect = &Collect{
			config: config,
		}
	)

	// Create a temp directory for the package
	if collect.config.TempPaths, err = createPaths(); err != nil {
		return nil, fmt.Errorf("Unable to create package temp paths: %w", err)
	}

	return collect, nil
}

func NewOrDie(config *types.EcosConfig) *Collect {
	var (
		err     error
		collect *Collect
	)

	if collect, err = New(config); err != nil {
		fmt.Printf("Unable to prepare for collect: %s", err.Error())
		os.Exit(1)
	}

	return collect
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
		return fmt.Errorf("Unable to read the ecos.yaml file: %w", err)
	}

	// Write out the Ecos Spec
	if err := utils.WriteYaml(filepath.Join(c.config.TempPaths.Base, "ecos.yaml"), c.config.Spec, 0766); err != nil {
		return fmt.Errorf("Unable to write the ecos.yaml file: %w", err)
	}

	for _, component := range c.config.Spec.Components {
		fmt.Printf("\nCOMPONENT %s\n\n", strings.ToUpper(component.Name))

		// copy dir to temp location
		if err := utils.Copy(component.Name, filepath.Join(c.config.TempPaths.Components, component.Name)); err != nil {
			return fmt.Errorf("Unable to copy component %s: %w", component.Name, err)
		}

		originalDir, _ := os.Getwd()

		// Providers: terraform providers mirror
		if err := os.Chdir(filepath.Join(c.config.TempPaths.Components, component.Name)); err != nil {
			return fmt.Errorf("Unable to access directory '%s': %w", component.Name, err)
		}

		if err := utils.ExecCommand("terraform", "providers", "mirror", "providers"); err != nil {
			return fmt.Errorf("Unable to mirror Terraform providers: %w", err)
		}

		// Modules: terraform get
		if err := utils.ExecCommand("terraform", "get"); err != nil {
			return fmt.Errorf("Unable to get terraform modules: %w", err)
		}

		os.Chdir(originalDir)
	}

	archiveName := "ecos-" + c.config.Spec.Metadata.Name + "-" + c.config.Spec.Metadata.Version + ".tar.zst"

	return utils.Archive(c.config.TempPaths.Base, archiveName)
}
