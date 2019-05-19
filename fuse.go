package main

import (
	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"github.com/iwanbk/kubefs/filesys"
	"github.com/iwanbk/kubefs/kube"
)

func mount(mountPoint, kubeCtx string, kubeCli *kube.Client) error {
	c, err := fuse.Mount(mountPoint)
	if err != nil {
		return err
	}
	defer c.Close()

	fileSys := filesys.NewFS(kubeCtx, kubeCli)

	err = fs.Serve(c, fileSys)
	if err != nil {
		return nil
	}

	// check if the mount process has an error to report
	<-c.Ready
	if err := c.MountError; err != nil {
		return err
	}

	return nil
}
