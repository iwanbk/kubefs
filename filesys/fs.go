package filesys

import (
	"bazil.org/fuse/fs"
)

type FS struct {
	kubeCtx string
	nss     []string
}

func NewFS(kubeCtx string, nss []string) *FS {
	return &FS{
		kubeCtx: kubeCtx,
		nss:     nss,
	}
}

func (f *FS) Root() (fs.Node, error) {
	return newRootDir(f.kubeCtx, f.nss), nil
}
