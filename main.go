package main

import (
	"log"
	"os"

	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	// init client
	var (
		kubeClient *kubernetes.Clientset
	)
	const (
		mountPoint = "/tmp/kfs"
	)

	_, err := rest.InClusterConfig()
	if err != nil {
		kubeClient = GetClientOutOfCluster()
	} else {
		kubeClient = GetClient()
	}

	if kubeClient == nil {
		log.Fatal("payah")
	}

	nss, err := getAllNamespacesName(kubeClient)
	if err != nil {
		log.Fatalf("failed to get namespaces: %v", err)
	}
	for i, ns := range nss {
		log.Printf("%d : %s", i, ns)
	}

	err = mount(mountPoint, "", nss)
	if err != nil {
		log.Fatalf("err = %v", err)
	}
	log.Println("OK")
}

func GetClient() *kubernetes.Clientset {
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatalf("Can not get kubernetes config: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Can not create kubernetes client: %v", err)
	}

	return clientset
}

// GetClientOutOfCluster returns a k8s clientset to the request from outside of cluster
func GetClientOutOfCluster() *kubernetes.Clientset {
	config, err := buildOutOfClusterConfig()
	if err != nil {
		log.Fatalf("Can not get kubernetes config: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("can't get out of cluster kube client: %v", err)
	}

	return clientset
}

func buildOutOfClusterConfig() (*rest.Config, error) {
	kubeconfigPath := os.Getenv("KUBECONFIG")
	if kubeconfigPath == "" {
		kubeconfigPath = os.Getenv("HOME") + "/.kube/config"
	}
	return clientcmd.BuildConfigFromFlags("", kubeconfigPath)
}
