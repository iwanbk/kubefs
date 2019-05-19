package filesys

import (
	"context"
	"os"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"github.com/iwanbk/kubefs/kube"
)

// rootDir represents kubefs root directory.
// this directory contains all namespace in the cluster
type rootDir struct {
	cli *kube.Client
}

func newRootDir(kubeCli *kube.Client) *rootDir {
	return &rootDir{
		cli: kubeCli,
	}
}

// Attr returns file attr of root dir.
// the inode is always 1
func (rd *rootDir) Attr(ctx context.Context, attr *fuse.Attr) error {
	attr.Inode = 1
	attr.Mode = os.ModeDir | 0555
	return nil
}

func (rd *rootDir) Lookup(ctx context.Context, name string) (fs.Node, error) {
	inode, ok := inoMgr.get(prefixNamespace, name)
	if !ok {
		return nil, fuse.ENOENT
	}

	return newNamespaceDir(inode, name), nil
}

// ReadDirAll
func (rd *rootDir) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	// get namespaces
	nss, err := rd.cli.GetNamespacesName()
	if err != nil {
		return nil, fuse.EIO
	}

	var (
		dirs = make([]fuse.Dirent, 0, len(nss))
	)

	for _, ns := range nss {
		dirs = append(dirs, fuse.Dirent{
			Inode: inoMgr.getOrCreate(prefixNamespace, ns),
			Name:  ns,
			Type:  fuse.DT_Dir,
		})
	}
	return dirs, nil
}
