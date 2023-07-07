package types

type EcosConfig struct {
	Spec             EcosSpec
	TempPaths        TempPaths
	PackageVariables EcosPackageVariables
}

type TempPaths struct {
	Base string
}

type EcosPackageVariables struct {
	VariableMap map[string]string
}
