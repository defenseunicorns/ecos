package archive

import (
	"bytes"
	"fmt"
	"html/template"
	"path/filepath"
	"strings"

	"github.com/mikevanhemert/ecos/src/pkg/utils"
	"github.com/mikevanhemert/ecos/src/types"
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

func (a *Archiver) HandleTemplates(componentName string, originalDir string, componentDir string) error {
	for name, properties := range a.config.Spec.Components[componentName].OutputTemplates {
		// 1. Extract outputs to map[string]interface{}{}
		outputs, err := collectTemplateVariables(properties.TemplateVariables)
		if err != nil {
			return fmt.Errorf("Unable to process template %s: %w", name, err)
		}

		// TODO delete
		fmt.Printf("outputs %s\n", outputs)

		// 2. Read the template file as bytes[] from spec.components[].templates[name].file
		sourcePath := filepath.Join(originalDir, componentName, properties.Source)
		templateFile, err := utils.ReadFile(sourcePath)
		if err != nil {
			return fmt.Errorf("Unable to read template file %s for template %s: %w", properties.Source, name, err)
		}

		// 3. Load the template
		tmpl, err := template.New(name).Parse(string(templateFile))
		if err != nil {
			return fmt.Errorf("unable to parse template %s: %w", name, err)
		}

		var tmplBuf bytes.Buffer
		tmpl.Execute(&tmplBuf, outputs)
		tmplBytes := tmplBuf.Bytes()

		// 4. Write the file
		utils.MakeDir(filepath.Join(originalDir, "out"))
		originalDestPath := filepath.Join(originalDir, "out", properties.Destination)
		componentDestPath := filepath.Join(componentDir, properties.Destination)
		utils.WriteFile(originalDestPath, tmplBytes)
		utils.WriteFile(componentDestPath, tmplBytes)
	}

	return nil
}

func collectTemplateVariables(variables map[string]types.EcosTemplateVariable) (map[string]interface{}, error) {
	outputs := map[string]interface{}{}
	envVars := []string{}

	for name, properties := range variables {
		options := []string{"output"}

		if len(properties.TerraformOptions) > 0 {
			options = append(options, properties.TerraformOptions...)
		}

		options = append(options, properties.TerraformName)

		if value, err := utils.ExecCommand("terraform", envVars, options...); err == nil {
			outputs[name] = strings.ReplaceAll(strings.TrimSpace(value), "\"", "")
		} else {
			return outputs, fmt.Errorf("Unable to extract output %s: %w", properties.TerraformName, err)
		}
	}

	return outputs, nil
}
