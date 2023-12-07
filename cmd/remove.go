package cmd

import (
	"log"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(removeCmd())
}

var removeCmdForceFlag bool

func removeCmd() *cobra.Command {
	removeCmd := &cobra.Command{
		Use:   "remove",
		Short: "Delete the specified version of kubectl",
		Args:  cobra.ExactArgs(1),
		Run:   remove,
	}

	removeCmd.Flags().BoolVarP(&removeCmdForceFlag, "force", "f", false, "Force delete this version")
	return removeCmd
}

func remove(cmd *cobra.Command, args []string) {
	version := args[0]
	path, err := filepath.EvalSymlinks(GetCurrentKubectlPath())
	if err != nil {
		log.Fatalln(err)
	}

	if !strings.HasSuffix(path, version) || removeCmdForceFlag {
		if KubectlIsExist(GetKubectlPathWithVersion(version)) {
			RemoveKubectl(GetKubectlPathWithVersion(version))
		}
		if strings.HasSuffix(path, version) {
			RemoveKubectl(GetCurrentKubectlPath())
		}
		return
	}

	log.Fatalln("This version is in use, deletion failed. Please add '-f' for a force remove")
}
