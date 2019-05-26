package kube

import (
	"context"
	"os"

	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type Client struct {
	cli *kubernetes.Clientset
}

func NewClient() (*Client, error) {
	// get config
	cfg, err := rest.InClusterConfig()
	if err != nil {
		cfg, err = buildOutOfClusterConfig()
	}
	if err != nil {
		return nil, err
	}

	// creates client
	cli, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}

	return &Client{
		cli: cli,
	}, nil
}

func (c *Client) Ping() error {
	_, err := c.GetNamespacesName()
	return err
}

func (c *Client) GetNamespacesName() ([]string, error) {
	var (
		api        = c.cli.CoreV1().Namespaces()
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

// GetPodsName get all pods name in a namespace
func (c *Client) GetPodsName(ns string) ([]string, error) {
	pods, err := c.cli.CoreV1().Pods(ns).List(meta_v1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var names []string
	for _, p := range pods.Items {
		names = append(names, p.ObjectMeta.Name)
	}
	return names, nil
}

func (c *Client) GetPodDescribe(ctx context.Context, ns, name string) ([]byte, error) {
	return c.kubeCliNs(ctx, ns, []string{"describe", "pod", name}...)
}

func (c *Client) GetPodLogs(ctx context.Context, ns, name string) ([]byte, error) {
	return c.kubeCliNs(ctx, ns, []string{"logs", name}...)
}

/*
func (c *Client) getPodDescribeKubeCli(ns, name string) ([]byte, error) {
	pod, err := c.cli.CoreV1().Pods(ns).Get(name, meta_v1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return json.MarshalIndent(pod, "", "\t")
}*/

func buildOutOfClusterConfig() (*rest.Config, error) {
	kubeconfigPath := os.Getenv("KUBECONFIG")
	if kubeconfigPath == "" {
		kubeconfigPath = os.Getenv("HOME") + "/.kube/config"
	}
	return clientcmd.BuildConfigFromFlags("", kubeconfigPath)
}
