package main

import (
	"os"

	"github.com/tendermint/tmlibs/cli"

	cmd "kchain/cmd/commands"
)

func main() {
	rootCmd := cmd.RootCmd

	rootCmd.AddCommand(
		cmd.GenValidatorCmd,
		cmd.InitFilesCmd,
		//cmd.ProbeUpnpCmd,
		//cmd.LiteCmd,
		//cmd.ReplayCmd,
		//cmd.ReplayConsoleCmd,
		cmd.ResetAllCmd,
		cmd.ResetPrivValidatorCmd,
		cmd.ShowValidatorCmd,
		//cmd.TestnetFilesCmd,
		cmd.VersionCmd,

		cmd.NewRunNodeCmd(),
		cmd.ServicesCmd(),
	)

	cmd1 := cli.PrepareBaseCmd(rootCmd, "TM", os.ExpandEnv("$HOME/.tendermint"))
	if err := cmd1.Execute(); err != nil {
		panic(err)
	}
}
