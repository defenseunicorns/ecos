package types

type EcosSpec struct {
	Kind        string          `json:"kind" jsonschema:"description=The kind of Ecos package,enum=EcosPackage"`
	EcosVersion string          `json:"ecosversion" jsonschema:"description=The Ecos binary version"`
	Metadata    EcosMetadata    `json:"metadata,omitempty" jsonschema:"description=Package metadata"`
	Components  []EcosComponent `json:"components" jsonschema:"description=Ecos component"`
}

type EcosMetadata struct {
	Name          string `json:"name" jsonschema:"description=Name that identifies this package,pattern=^[a-z0-9\\-]+$"`
	Description   string `json:"description,omitempty" jsonschema:"description=Additional information about this package"`
	Version       string `json:"version" jsonschema:"description=Generic string set by package author to track the package version"`
	Authors       string `json:"authors,omitempty" jsonschema:"description=Comma-separated list of package authors (including contact info),example=Josh &#60;hello@defenseunicorns.com&#62; Gedd Josh &#60;hello@defenseunicorns.com&#62;"`
	Documentation string `json:"documenation,omitempty" jsonschema:"description=Link to package documentation"`
}

type EcosComponent struct {
	Name            string               `json:"name" jsonschema:"description=Name that identifies this component,pattern=^[a-z0-9\\-]+$"`
	Description     string               `json:"description,omitempty" jsonschema:"description=Additional information about this component"`
	Variables       []EcosVariable       `json:"variables,omitempty" jsonschema:"description=Variable template values"`
	Transitives     []EcosVariable       `json:"transitives,omitempty" jsonschema:"description=Transitive variables to pass to subsequent components"`
	OutputTemplates []EcosOutputTemplate `json:"templates,omitempty" jsonschema:"description=Templated output files and Terraform output mappings"`
}

type EcosVariable struct {
	Name        string `json:"name" jsonschema:"description=The variable name,pattern=^[A-Z0-9_]+$"`
	Description string `json:"description,omitempty" jsonschema:"description=A description of the variable"`
	Default     string `json:"default,omitempty" jsonschema:"description=The default value to use fo the variable"`
}

type EcosOutputTemplate struct {
	Template          string                 `json:"template" jsonschema:"description=The source template file with relative path (if needed)"`
	Description       string                 `json:"description,omitempty" jsonschema:"description=A description of the template file"`
	TemplateVariables []EcosTemplateVariable `json:"variables" jsonschema:"description=Terraform outputs to write into this template"`
}

type EcosTemplateVariable struct {
	Name             string   `json:"name" jsonschema:"description=The template variable name to replaces,pattern=^[a-z0-9\\-]+$"`
	Description      string   `json:"description,omitaempty" jsonschema:"description=A description of the variable"`
	TerraformName    string   `json:"tfname" jsonschema:"description=Terraform output to capture"`
	TerraformOptions []string `json:"tfoptions,omitempty" jsonschema:"description=Terraform output options (-json, -raw)"`
}