package main

import (
	"github.com/iwanbk/kubefs/kube"
	log "github.com/sirupsen/logrus"
)

func main() {
	const (
		mountPoint = "/tmp/kfs"
	)

	log.Info("initializing kubernetes client")
	// init client
	cli, err := kube.NewClient()
	if err != nil {
		log.Fatalf("failed to creates kube client: %v", err)
	}

	log.Info("pinging kubernetes cluster")
	// test connection
	err = cli.Ping()
	if err != nil {
		log.Fatalf("failed to get namespace: %v", err)
	}

	err = mount(mountPoint, "", cli)
	if err != nil {
		log.Fatalf("err = %v", err)
	}
	log.Println("OK")
}
