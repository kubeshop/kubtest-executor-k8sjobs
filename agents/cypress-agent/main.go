package main

import (
	"fmt"

	"github.com/kubeshop/kubtest-executor-k8sjobs/pkg/cypress"
	"github.com/kubeshop/kubtest/pkg/api/kubtest"
	"github.com/spf13/cobra"
)

var (
	id         string
	script     string
	repository string
	branch     string
	path       string
)

func main() {

	var RootCmd = &cobra.Command{
		Use:   "kubtest",
		Short: "kubtest entrypoint for plugin",
		Long:  `kubtest`,
		Run: func(cmd *cobra.Command, args []string) {
			runner := &cypress.CypressRunner{}
			result := runner.Run(kubtest.Execution{
				Repository: &kubtest.Repository{
					Uri:    repository,
					Branch: branch,
					Path:   path,
				},
			})
			fmt.Println(result)
			fmt.Printf("$$$%s$$$", id)
		},
	}

	RootCmd.SilenceUsage = true
	RootCmd.Flags().StringVarP(&id, "id", "i", "", "input script")
	RootCmd.Flags().StringVarP(&script, "script", "s", "", "input script")
	RootCmd.Flags().StringVarP(&repository, "repository", "r", "", "repository path")
	RootCmd.Flags().StringVarP(&branch, "branch", "b", "", "branch")
	RootCmd.Flags().StringVarP(&path, "path", "p", "", "path")

	RootCmd.Execute()
}
