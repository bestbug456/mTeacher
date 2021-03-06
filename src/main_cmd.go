package main

import (
	"fmt"
)

import (
	"github.com/spf13/cobra"
)

import (
	"command"
)

// Root command
var mDotaCmd = &cobra.Command{
	Use:   "mTeacher",
	Short: "main command for this project",
	Long:  `this command is the entrypoint for our project.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Execution implementation
		return fmt.Errorf("undefined command, select a valid one")
	},
}

func PrepareMainCommand() error {
	// add all the subcommand
	mDotaCmd.AddCommand(command.TeachCmd)
	return nil
}
