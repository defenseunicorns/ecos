package archive

import (
	"fmt"

	"github.com/mikevanhemert/ecos/src/pkg/utils"
)

func (a *Archiver) CollectTransitives(componentName string) error {
	envVars := []string{}

	for name, properties := range a.config.Spec.Components[componentName].Transitives {
		options := []string{"output"}

		if len(properties.TerraformOptions) > 0 {
			options = append(options, properties.TerraformOptions...)
		}

		options = append(options, properties.TerraformName)

		if value, err := utils.ExecCommand("terraform", envVars, options...); err == nil {
			a.config.PackageVariables[name] = value
		} else {
			return fmt.Errorf("Unable to extract output %s: %w", properties.TerraformName, err)
		}
	}

	return nil
}
