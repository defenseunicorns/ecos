package types

type EcosConfig struct {
	Spec             EcosSpec
	TempPaths        TempPaths
	PackageVariables map[string]string
}

type TempPaths struct {
	Base string
}
