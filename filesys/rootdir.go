package filesys

import (
	"context"
	"os"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"github.com/iwanbk/kubefs/kube"
	log "github.com/sirupsen/logrus"
)

type rootDir struct {
	kubeCtx string
	cli     *kube.Client
}

func newRootDir(kubeCtx string, kubeCli *kube.Client) *rootDir {
	return &rootDir{
		kubeCtx: kubeCtx,
		cli:     kubeCli,
	}
}

func (rd *rootDir) Attr(ctx context.Context, attr *fuse.Attr) error {
	attr.Inode = 1
	attr.Mode = os.ModeDir | 0555
	return nil
}

func (rd *rootDir) Lookup(ctx context.Context, name string) (fs.Node, error) {
	log.Infof("rd lookup %v", name)
	if name == "hello" {
		return File{}, nil
	}
	return nil, fuse.ENOENT
}

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
			Inode: inoMgr.get("ns", ns),
			Name:  ns,
			Type:  fuse.DT_Dir,
		})
	}
	return dirs, nil
}

// File implements both Node and Handle for the hello file.
type File struct{}

const greeting = "hello, world\n"

func (File) Attr(ctx context.Context, a *fuse.Attr) error {
	a.Inode = 2
	a.Mode = 0444
	a.Size = uint64(len(greeting))
	return nil
}

func (File) ReadAll(ctx context.Context) ([]byte, error) {
	return []byte(greeting), nil
}
