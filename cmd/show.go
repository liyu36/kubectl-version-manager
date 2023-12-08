package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(showCmd())
}

func showCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show",
		Short: "View the currently installed version of kubectl",
		Run:   show,
	}
}

func show(cmd *cobra.Command, args []string) {
	files, err := os.ReadDir(GetBaseDir())
	if err != nil {
		log.Fatalln(err)
	}

	path, err := filepath.EvalSymlinks(GetCurrentKubectlPath())
	if err != nil {
		log.Fatalln(err)
	}

	basename := filepath.Base(path)

	for _, file := range files {
		if !strings.HasPrefix(file.Name(), "kubectl-") {
			continue
		}

		if strings.EqualFold(file.Name(), basename) {
			fmt.Printf("%s *\n", strings.TrimPrefix(file.Name(), "kubectl-"))
			continue
		}
		fmt.Println(strings.TrimPrefix(file.Name(), "kubectl-"))
	}
}
