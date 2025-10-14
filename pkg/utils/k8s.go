package utils

import (
	"github.com/flanksource/commons/logger"
	"github.com/flanksource/kommons"
	"k8s.io/client-go/kubernetes"
)

var K8sClientSet *kubernetes.Clientset

func init() {
	client, err := kommons.NewClientFromDefaults(logger.StandardLogger())
	if err != nil {
		// If this errors app should not start
		panic(err)
	}
	cc, err := client.GetClientset()
	if err != nil {
		// If this errors app should not start
		panic(err)
	}
	K8sClientSet = cc
}
