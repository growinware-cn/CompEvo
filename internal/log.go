package internal

import (
	"bytes"
	"io"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

func GetLog(kubeClient kubernetes.Interface, namespace string, podName string) (string, error) {
	req := kubeClient.CoreV1().Pods(namespace).GetLogs(podName, &v1.PodLogOptions{})
	podLogs, err := req.Stream()
	if err != nil {
		return "", err
	}
	defer podLogs.Close()

	buffer := new(bytes.Buffer)
	if _, err := io.Copy(buffer, podLogs); err != nil {
		return "", nil
	}

	return buffer.String(), nil
}
