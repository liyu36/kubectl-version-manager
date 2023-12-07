package cmd

import (
	"fmt"
	"log"
	"os"
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

	for _, file := range files {
		if strings.HasPrefix(file.Name(), "kubectl-") {
			fmt.Println(strings.TrimPrefix(file.Name(), "kubectl-"))
		}
	}
}
