package lang

const (
	// ecos root command
	RootCmdShort = "Terraform for Air Gap"
	RootCmdLong  = "Ecos eliminates the complexity of air gap infrastructure provisioning, patching, and maintaining with Terraform using a declarative packaging strategy to support DevSecOps in offline and semi-connected environments."

	RootCmdFlagLogLevel = "Log level when running Ecos. Valid options are: warn, info, debug, trace"

	RootCmdErrInvalidLogLevel = "Invalid log level. Valid options are: warn, info, debug, trace."

	// ecos package command
	CmdPackageShort = "Ecos package commands for creating, applying, and repairing packages"

	CmdPackageCreateShort = "Creates an Ecos package from a given directory of the current dirctory"
	CmdPackageCreateLong  = "Builds an archive of resources and dependencies defined by the 'ecos.yaml' in the specified directory."

	CmdPackageApplyShort = "Applies the Ecos package Terraform from a local file or URL (runs offline)"
	CmdPackageApplyLong  = "Unpacks Terraform resources and dependencies from an Ecos package archive and applies them."

	CmdPackageApplyFlagSet = "Specify package variables (KEY=value)"

	// ecos state command

)
