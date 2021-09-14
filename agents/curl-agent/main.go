package main

import (
	"fmt"

	"github.com/kubeshop/kubtest-executor-k8sjobs/pkg/curl"
	"github.com/kubeshop/kubtest/pkg/api/kubtest"
	"github.com/spf13/cobra"
)

var (
	id         string
	script     string
	repository string
	params     []string
)

func main() {

	var RootCmd = &cobra.Command{
		Use:   "kubtest",
		Short: "kubtest entrypoint for plugin",
		Long:  `kubtest`,
		Run: func(cmd *cobra.Command, args []string) {
			runner := &curl.NewCurlRunner{}
			result := runner.Run(kubtest.Execution{
				ScriptContent: script,
			})
			fmt.Println(result)
			fmt.Printf("$$$%s$$$", id)
		},
	}

	RootCmd.SilenceUsage = true
	RootCmd.Flags().StringVarP(&id, "id", "i", "", "input script")
	RootCmd.Flags().StringVarP(&script, "script", "s", "", "input script")
	RootCmd.Flags().StringVarP(&repository, "repository", "r", "", "repository path")
	RootCmd.Flags().StringArrayVarP(&params, "parameters", "p", []string{""}, "input script")

	RootCmd.Execute()
}
