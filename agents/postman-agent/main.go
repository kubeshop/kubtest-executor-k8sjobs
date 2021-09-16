package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/kubeshop/kubtest-executor-k8sjobs/pkg/newman"
	"github.com/kubeshop/kubtest/pkg/api/kubtest"
)

func main() {

	args := os.Args
	if len(args) == 1 {
		fmt.Println("missing input argument")
		os.Exit(1)
	}

	script := args[1]

	e := kubtest.Execution{}
	json.Unmarshal([]byte(script), &e)
	runner := &newman.NewmanRunner{}
	result := runner.Run(e)

	fmt.Println(result)
	fmt.Printf("$$$%s$$$", e.Id)
}
