package filesys

import (
	"context"
	"os"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"github.com/iwanbk/kubefs/kube"
)

// namespaceDir dir represents a namespace directory
type namespaceDir struct {
	inode uint64
	name  string
	cli   *kube.Client
}

func newNamespaceDir(inode uint64, name string, cli *kube.Client) *namespaceDir {
	return &namespaceDir{
		inode: inode,
		name:  name,
		cli:   cli,
	}
}

func (nd *namespaceDir) Attr(ctx context.Context, attr *fuse.Attr) error {
	attr.Inode = nd.inode
	attr.Mode = os.ModeDir | 0444
	return nil
}

func (nd *namespaceDir) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	dirs := []fuse.Dirent{
		{
			Inode: inoMgr.getOrCreate(prefixNamespace, nd.name, prefixPod),
			Name:  dirPodName,
			Type:  fuse.DT_Dir,
		},
	}
	return dirs, nil
}

func (nd *namespaceDir) Lookup(ctx context.Context, name string) (fs.Node, error) {
	inode, ok := inoMgr.get(prefixNamespace, nd.name, name)
	if !ok {
		return nil, fuse.ENOENT
	}

	return newPodDir(inode, nd.name, name, nd.cli), nil
}
