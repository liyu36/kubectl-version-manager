package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(useCmd())
}

func useCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "use",
		Aliases: []string{"set"},
		Short:   "Use the specified version of kubectl",
		Args:    cobra.ExactArgs(1),
		Example: fmt.Sprintf("  %s %s v1.28.4", rootCmd.Name(), "use"),
		Run:     use,
	}
}

func use(cmd *cobra.Command, args []string) {
	version := args[0]
	if KubectlIsExist(GetKubectlPathWithVersion(version)) {
		if err := SetKubectlSysLink(args[0]); err != nil {
			log.Fatalln(err)
		}
		return
	}

	log.Printf("%s not found, will be automatically downloaded", GetKubectlPathWithVersion(version))
	for i := 0; i < 3; i++ {
		fmt.Print(".")
		time.Sleep(1 * time.Second)
	}
	fmt.Println()
	install(cmd, args)
}
