package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/b4b4r07/hal-ops/command"
	"github.com/b4b4r07/hal-ops/command/git"
	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Checkout branch and validate config",
	Long:  "Check if there are problems in config file",
	RunE:  check,
}

func check(cmd *cobra.Command, args []string) error {
	var c *command.Command

	branch, err := git.GetCurrentBranchName()
	if err != nil {
		return err
	}

	// Checkout to it if an argument is given
	if len(args) > 0 {
		branch = args[0]
		if err := git.Checkout(branch); err != nil {
			return err
		}
	}

	if branch == git.MASTER {
		return fmt.Errorf("Error: you are on master branch")
	}

	// Copy new hal config to hal directory
	var (
		src  = filepath.Join(".hal", "config")
		dest = filepath.Join(os.Getenv("HOME"), ".hal", "config")
	)
	// Backup
	data, err := ioutil.ReadFile(dest)
	if err != nil {
		return err
	}
	if err := command.Cp(src, dest); err != nil {
		return err
	}

	// Validation
	c = command.New("hal config >/dev/null")
	defer func() {
		// Restore dest file if syntax is wrong
		ioutil.WriteFile(dest, data, 0644)
	}()
	return c.RunWithTTY()
}

func init() {
	RootCmd.AddCommand(checkCmd)
}
