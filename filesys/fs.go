package filesys

import (
	"bazil.org/fuse/fs"
	"github.com/iwanbk/kubefs/kube"
)

type FS struct {
	cli  *kube.Client
	root *rootDir
}

func NewFS(kubeCli *kube.Client) *FS {
	return &FS{
		cli:  kubeCli,
		root: newRootDir(kubeCli),
	}
}

func (f *FS) Root() (fs.Node, error) {
	return f.root, nil
}

const (
	permFile    = 0444
	permDir     = 0444
	permRootDir = 0555 // TODO: not sure why it is not 0444
)

const (
	prefixNamespace = "root"
	prefixPod       = "pods"
)

const (
	dirPodName = "pods"
)
