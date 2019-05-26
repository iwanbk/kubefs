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
	attr.Mode = os.ModeDir | permDir
	return nil
}

func (nd *namespaceDir) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	dirs := []fuse.Dirent{
		{
			Inode: inoMgr.getOrCreate(dirNamespace, nd.name, dirPod),
			Name:  dirPod,
			Type:  fuse.DT_Dir,
		},
		{
			Inode: inoMgr.getOrCreate(dirNamespace, nd.name, dirDeployment),
			Name:  dirDeployment,
			Type:  fuse.DT_Dir,
		},
	}
	return dirs, nil
}

func (nd *namespaceDir) Lookup(ctx context.Context, name string) (fs.Node, error) {
	inode, ok := inoMgr.get(dirNamespace, nd.name, name)
	if !ok {
		return nil, fuse.ENOENT
	}

	switch name {
	case dirDeployment:
		return newDeploymentDir(inode, nd.name, name, nd.cli), nil
	case dirPod:
		return newPodDir(inode, nd.name, name, nd.cli), nil
	default:
		return nil, fuse.ENOENT
	}
}
