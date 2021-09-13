package k8s

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/kubeshop/kubtest/pkg/api/kubtest"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type Client struct {
	ClientSet *kubernetes.Clientset
	Namespace string
	Cmd       string
}

func NewClient() (*Client, error) {
	clientSet, err := connectToK8s()
	if err != nil {
		return nil, err
	}

	return &Client{
		ClientSet: clientSet,
	}, nil
}

func (c *Client) LaunchK8sJob(jobName string, execution kubtest.Execution) *kubtest.ExecutionResult {
	jobs := c.ClientSet.BatchV1().Jobs(c.Namespace)
	var result string

	image := "jasmingacic/test"
	id := fmt.Sprintf("--id=%s", execution.Id)
	script := fmt.Sprintf("--script=%s", execution.ScriptContent)

	var backOffLimit int32 = 0
	jobSpec := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      jobName,
			Namespace: c.Namespace,
		},
		Spec: batchv1.JobSpec{
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:            jobName,
							Image:           image,
							Command:         []string{"agent", id, script},
							ImagePullPolicy: v1.PullAlways,
						},
					},
					RestartPolicy: v1.RestartPolicyNever,
				},
			},
			BackoffLimit: &backOffLimit,
		},
	}

	job, err := jobs.Create(context.TODO(), jobSpec, metav1.CreateOptions{})
	if err != nil {
		log.Println("Failed to create K8s job.", err)
	}

	//print job details
	time.Sleep(2 * time.Second)

	pods, err := c.ClientSet.CoreV1().Pods("default").List(context.TODO(), metav1.ListOptions{LabelSelector: "job-name=" + job.Name})
	// pods, err := c.ClientSet.CoreV1().Pods("default").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return &kubtest.ExecutionResult{
			Status:       kubtest.ExecutionStatusError,
			ErrorMessage: err.Error(),
		}
	}

	for _, pod := range pods.Items {
		if pod.Labels["job-name"] == jobName {
			if pod.Status.Phase == v1.PodSucceeded {
				if err := wait.PollImmediate(time.Second, time.Duration(0)*time.Second, isPodRunning(c.ClientSet, pod.Name, c.Namespace)); err != nil {
					fmt.Println(err)
					return &kubtest.ExecutionResult{
						Status:       kubtest.ExecutionStatusError,
						ErrorMessage: err.Error(),
					}
				}
			}
			result, err = c.GetPodLogs(pod.Name, jobName, execution.Id)
			if err != nil {
				return &kubtest.ExecutionResult{
					Status:       kubtest.ExecutionStatusError,
					ErrorMessage: err.Error(),
				}
			}
		}
	}

	return &kubtest.ExecutionResult{
		Status:    kubtest.ExecutionStatusSuceess,
		RawOutput: result,
	}
}

// connectToK8s returns ClientSet
func connectToK8s() (*kubernetes.Clientset, error) {
	config, err := clientcmd.BuildConfigFromFlags("", os.Getenv("KUBECONFIG"))
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}

// isPodRunning check if the pod in question is running state
func isPodRunning(c *kubernetes.Clientset, podName, namespace string) wait.ConditionFunc {
	return func() (bool, error) {
		pod, err := c.CoreV1().Pods(namespace).Get(context.Background(), podName, metav1.GetOptions{})
		if err != nil {
			return false, err
		}

		switch pod.Status.Phase {
		case v1.PodRunning, v1.PodSucceeded:
			return true, nil
		case v1.PodFailed:
			return false, nil
		}
		return false, nil
	}
}

func (c *Client) GetPodLogs(podName string, containerName string, endMessage string) (string, error) {
	count := int64(100)
	var toReturn string
	var message string
	podLogOptions := v1.PodLogOptions{
		Follow:    true,
		TailLines: &count,
	}

	podLogRequest := c.ClientSet.CoreV1().
		Pods(c.Namespace).
		GetLogs(podName, &podLogOptions)
	stream, err := podLogRequest.Stream(context.TODO())
	if err != nil {
		return "", err
	}

	defer stream.Close()

	for {
		buf := make([]byte, 2000)
		numBytes, err := stream.Read(buf)
		if numBytes == 0 {
			break
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}

		message = string(buf[:numBytes])
		if strings.Contains(message, fmt.Sprintf("$$$%s$$$", endMessage)) {
			message = ""
			break
		} else {
			toReturn += message
		}
	}
	return toReturn, nil
}
