package kube

import (
	"context"

	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetDeploymentsName get all deployments name in a namespace
func (c *Client) GetDeploymentsName(ns string) ([]string, error) {
	deployments, err := c.cli.AppsV1().Deployments(ns).List(meta_v1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var names []string
	for _, p := range deployments.Items {
		names = append(names, p.ObjectMeta.Name)
	}
	return names, nil
}

func (c *Client) GetDeploymentDescribe(ctx context.Context, ns, name string) ([]byte, error) {
	return c.kubeCliNs(ctx, ns, []string{"describe", "deployment", name}...)
}
