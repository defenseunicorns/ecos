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

type Destroy struct {
	archiveName string
	config      *types.EcosConfig
}

func NewDestroy(archiveName string) (*Destroy, error) {
	// Create a temp directory
	var (
		err     error
		destroy = &Destroy{
			config:      &types.EcosConfig{},
			archiveName: archiveName,
		}
	)

	if destroy.config.TempPaths.Base, err = utils.MakeTempDir(""); err != nil {
		return nil, fmt.Errorf("Unable to create temp path: %w", err)
	}

	return destroy, nil
}

func NewOrDieDestroy(archiveName string) *Destroy {
	var (
		err     error
		destroy *Destroy
	)

	if destroy, err = NewDestroy(archiveName); err != nil {
		fmt.Printf("Unable to prepare for destroy: %s", err)
		os.Exit(1)
	}

	return destroy
}

func (d *Destroy) ClearTempPaths() {
	_ = os.RemoveAll(d.config.TempPaths.Base)
}

func (d *Destroy) Destroy() error {
	var err error

	// Unarchive to the temp dir and read spec
	if d.config.Spec, err = utils.Unarchive(d.archiveName, d.config.TempPaths.Base); err != nil {
		return err
	}

	for _, component := range d.config.Spec.Components {
		fmt.Printf("\nCOMPONENT %s\n\n", strings.ToUpper(component.Name))

		originalDir, _ := os.Getwd()

		// terraform destoy
		if err := os.Chdir(filepath.Join(d.config.TempPaths.Base, "components", component.Name)); err != nil {
			return fmt.Errorf("Unable to access directory '%s': %s", component.Name, err)
		}

		if err := utils.ExecCommand("terraform", "destroy", "-auto-approve"); err != nil {
			return fmt.Errorf("Unable to destroy Terraform: %w", err)
		}

		os.Chdir(originalDir)
	}

	utcTime := time.Now().UTC()
	archiveName := "ecos-" + d.config.Spec.Metadata.Name + "-" + d.config.Spec.Metadata.Version + "-" + utcTime.Format(time.RFC3339Nano) + ".tar.zst"

	return utils.Archive(d.config.TempPaths.Base, archiveName)
}
