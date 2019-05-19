package main

import (
	"log"

	"github.com/iwanbk/kubefs/kube"
)

func main() {
	const (
		mountPoint = "/tmp/kfs"
	)

	// init client
	cli, err := kube.NewClient()
	if err != nil {
		log.Fatalf("failed to creates kube client: %v", err)
	}

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
