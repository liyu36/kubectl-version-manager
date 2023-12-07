package cmd

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(currentCmd())
}

func currentCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "current",
		Short: "Display the current version of kubectl",
		Args:  cobra.MaximumNArgs(0),
		Run:   current,
	}
}

func current(cmd *cobra.Command, args []string) {
	path, err := filepath.EvalSymlinks(GetCurrentKubectlPath())
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%s -> %s\n", GetCurrentKubectlPath(), path)
}
