package cmd

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

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
	s := strings.TrimPrefix(path, filepath.Join(GetBaseDir(), "kubectl-"))
	fmt.Println(s)
}
