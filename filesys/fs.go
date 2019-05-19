package filesys

import (
	"bazil.org/fuse/fs"
	"github.com/iwanbk/kubefs/kube"
)

type FS struct {
	kubeCtx string
	cli     *kube.Client
}

func NewFS(kubeCtx string, kubeCli *kube.Client) *FS {
	return &FS{
		kubeCtx: kubeCtx,
		cli:     kubeCli,
	}
}

func (f *FS) Root() (fs.Node, error) {
	return newRootDir(f.kubeCtx, f.cli), nil
}
