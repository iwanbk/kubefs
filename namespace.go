package main

import (
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func getAllNamespacesName(kcli kubernetes.Interface) ([]string, error) {
	var (
		api        = kcli.CoreV1().Namespaces()
		namespaces []string
	)

	nsList, err := api.List(meta_v1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, ns := range nsList.Items {
		namespaces = append(namespaces, ns.ObjectMeta.Name)
	}
	return namespaces, nil
}
