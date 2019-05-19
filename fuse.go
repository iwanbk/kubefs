package main

import (
	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"github.com/iwanbk/kubefs/filesys"
)

func mount(mountPoint, kubeCtx string, nss []string) error {
	c, err := fuse.Mount(mountPoint)
	if err != nil {
		return err
	}
	defer c.Close()

	fileSys := filesys.NewFS(kubeCtx, nss)

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
