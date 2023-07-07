package lang

const (
	// ecos root command
	RootCmdShort = "Terraform for Air Gap"
	RootCmdLong  = "Ecos eliminates the complexity of air gap infrastructure provisioning, patching, and maintaining with Terraform using a declarative packaging strategy to support DevSecOps in offline and semi-connected environments."

	RootCmdFlagLogLevel = "Log level when running Ecos. Valid options are: warn, info, debug, trace"

	RootCmdErrInvalidLogLevel = "Invalid log level. Valid options are: warn, info, debug, trace."

	// ecos collect command
	CmdCollectShort = "Collects and packages infrastructure resources into an Ecos package"
	CmdCollectLong  = "Builds an archive of resources and dependencies defined by the 'ecos.yaml' in the specified directory."

	// ecos apply command
	CmdApplyShort = "Applies the Ecos package Terraform to a new environment"
	CmdApplyLong  = "Unpacks Terraform resources and dependencies from an Ecos package archive and applies them to a new environment."

	CmdApplyFlagSet = "Specify package variables (KEY=value)"

	// ecos update command
	CmdUpdateShort = "Applies the Ecos package Terraform as an update to an existing environment"
	CmdUpdateLong  = "Unpacks Terraform resources and dependencies from an Ecos package archive and applies them as an update to an existing environment."

	CmdUpdateFlagSet = "Specify package variables (KEY=value)"
)
