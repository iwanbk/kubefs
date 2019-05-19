package filesys

import (
	"bazil.org/fuse/fs"
	"github.com/iwanbk/kubefs/kube"
)

type FS struct {
	kubeCtx string
	cli     *kube.Client
	root    *rootDir
}

func NewFS(kubeCtx string, kubeCli *kube.Client) *FS {
	return &FS{
		kubeCtx: kubeCtx,
		cli:     kubeCli,
		root:    newRootDir(kubeCtx, kubeCli),
	}
}

func (f *FS) Root() (fs.Node, error) {
	return f.root, nil
}
