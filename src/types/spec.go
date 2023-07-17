package types

type EcosSpec struct {
	Kind        string                   `json:"kind" jsonschema:"description=The kind of Ecos package,enum=EcosPackage"`
	EcosVersion string                   `json:"ecosversion" jsonschema:"description=The Ecos binary version"`
	Metadata    EcosMetadata             `json:"metadata,omitempty" jsonschema:"description=Package metadata"`
	Components  map[string]EcosComponent `json:"components" jsonschema:"description=Ecos component"`
}

type EcosMetadata struct {
	Name          string `json:"name" jsonschema:"description=Name that identifies this package,pattern=^[a-z0-9\\-]+$"`
	Description   string `json:"description,omitempty" jsonschema:"description=Additional information about this package"`
	Version       string `json:"version" jsonschema:"description=Generic string set by package author to track the package version"`
	Authors       string `json:"authors,omitempty" jsonschema:"description=Comma-separated list of package authors (including contact info),example=Josh &#60;hello@defenseunicorns.com&#62; Gedd Josh &#60;hello@defenseunicorns.com&#62;"`
	Documentation string `json:"documenation,omitempty" jsonschema:"description=Link to package documentation"`
}

type EcosComponent struct {
	Description     string                          `json:"description,omitempty" jsonschema:"description=Additional information about this component"`
	Variables       map[string]EcosVariable         `json:"variables,omitempty" jsonschema:"description=Variable template values"`
	Transitives     map[string]EcosTemplateVariable `json:"transitives,omitempty" jsonschema:"description=Transitive variables to pass to subsequent components"`
	OutputTemplates map[string]EcosOutputTemplate   `json:"templates,omitempty" jsonschema:"description=Templated output files and Terraform output mappings"`
}

type EcosVariable struct {
	Description string `json:"description,omitempty" jsonschema:"description=A description of the variable"`
	Default     string `json:"default,omitempty" jsonschema:"description=The default value to use fo the variable"`
	Override    string `json:"override,omitempty" jsonschema:"description=An override supplied by the user during ecos apply or ecos update, reserved for ecos use"`
}

type EcosOutputTemplate struct {
	File              string                          `json:"file" jsonschema:"description=The source template file with relative path (if needed)"`
	Description       string                          `json:"description,omitempty" jsonschema:"description=A description of the template file"`
	TemplateVariables map[string]EcosTemplateVariable `json:"variables" jsonschema:"description=Terraform outputs to write into this template"`
}

type EcosTemplateVariable struct {
	OutputName       string   `json:"outputname" jsonschema:"description=The template variable name to replaces,pattern=^[a-z0-9\\-]+$"`
	Description      string   `json:"description,omitaempty" jsonschema:"description=A description of the variable"`
	TerraformName    string   `json:"tfname" jsonschema:"description=Terraform output to capture"`
	TerraformOptions []string `json:"tfoptions,omitempty" jsonschema:"description=Terraform output options (-json, -raw)"`
}
