package types

type EcosConfig struct {
	PackageVariables EcosPackageVariables
	StateOpts        EcosStateOptions
}

type EcosPackageVariables struct {
	VariableMap map[string]string
}

type EcosStateOptions struct {
}
