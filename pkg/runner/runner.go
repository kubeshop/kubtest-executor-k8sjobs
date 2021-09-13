package runner

import (
	"fmt"

	"github.com/kubeshop/kubtest-executor-template/pkg/k8s"
	"github.com/kubeshop/kubtest/pkg/api/kubtest"
)

func NewRunner() *JobsRunner {
	return &JobsRunner{}
}

// ExampleRunner for template - change me to some valid runner
type JobsRunner struct {
}

func (r *JobsRunner) Run(execution kubtest.Execution) kubtest.ExecutionResult {
	fmt.Println("TADA")
	client, err := k8s.NewClient()
	client.Namespace = "default"
	if err != nil {
		return kubtest.ExecutionResult{
			Status:    kubtest.ExecutionStatusError,
			RawOutput: err.Error(),
		}
	}

	return *client.LaunchK8sJob(execution.Id, execution)
}
