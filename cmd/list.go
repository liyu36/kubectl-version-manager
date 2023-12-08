package cmd

import (
	"fmt"
	"os/exec"
	"regexp"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listCmd())
}

func listCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   fmt.Sprintf("Retrieve the list of kubectl versions from %s", GetkubernetesSourceGitURL()),
		RunE:    list,
	}
}

func list(cmd *cobra.Command, args []string) error {
	var regexp string
	if len(args) > 0 {
		regexp = args[0]
	}
	vers, err := ListKubectlVersions(regexp)
	if err != nil {
		return err
	}

	for _, ver := range vers {
		fmt.Println(ver)
	}

	return nil
}

func ListKubectlVersions(re string) ([]string, error) {
	if re == "" {
		re = ".+"
	} else {
		re = fmt.Sprintf(`.*%s.*`, re)
	}

	cmd := exec.Command("git", "ls-remote", "--sort=version:refname", "--tags", GetkubernetesSourceGitURL())
	refs, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	r := regexp.MustCompile(fmt.Sprintf(`refs/tags/(%s)`, re))
	match := r.FindAllStringSubmatch(string(refs), -1)
	if match == nil {
		return nil, fmt.Errorf("No kubectl version found")
	}

	var vers []string
	for _, m := range match {
		r := regexp.MustCompile(".+{}") // 过滤掉以^{}结尾的tag
		if r.MatchString(m[1]) {
			continue
		}
		vers = append(vers, m[1])
	}

	return vers, nil
}
