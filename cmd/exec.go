package cmd

import (
	"github.com/b4b4r07/hal-ops/cli"
	"github.com/spf13/cobra"
)

var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "exec",
	Long:  "exec",
	RunE:  exec,
}

func exec(cmd *cobra.Command, args []string) error {
	return cli.Run("echo", "test")
}

func init() {
	RootCmd.AddCommand(execCmd)
}
