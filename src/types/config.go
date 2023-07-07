package types

type EcosConfig struct {
	Spec             EcosSpec
	TempPaths        TempPaths
	PackageVariables EcosPackageVariables
}

type TempPaths struct {
	Base       string
	Components string
}

type EcosPackageVariables struct {
	VariableMap map[string]string
}
