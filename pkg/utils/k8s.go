package utils

import (
	"github.com/flanksource/commons/logger"
	"github.com/flanksource/kommons"
	"k8s.io/client-go/kubernetes"
)

func SetupK8sClientSet() (*kubernetes.Clientset, error) {
	client, err := kommons.NewClientFromDefaults(logger.StandardLogger())
	if err != nil {
		return nil, err
	}
	cc, err := client.GetClientset()
	if err != nil {
		return nil, err
	}
	return cc, nil
}
