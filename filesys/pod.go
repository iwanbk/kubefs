package filesys

import (
	"context"
	"os"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"github.com/iwanbk/kubefs/kube"
)

// podDir dir represents a pod directory
type podDir struct {
	inode uint64
	name  string
	ns    string
	cli   *kube.Client
}

func newPodDir(inode uint64, ns, name string, cli *kube.Client) *podDir {
	return &podDir{
		inode: inode,
		ns:    ns,
		name:  name,
		cli:   cli,
	}
}

func (pd *podDir) Attr(ctx context.Context, attr *fuse.Attr) error {
	attr.Inode = pd.inode
	attr.Mode = os.ModeDir | permDir
	return nil
}

func (pd *podDir) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	names, err := pd.cli.GetPodsName(pd.ns)
	if err != nil {
		return nil, fuse.EIO
	}

	dirs := make([]fuse.Dirent, 0, len(names))
	for _, name := range names {
		dirs = append(dirs, fuse.Dirent{
			Inode: inoMgr.getOrCreate(prefixNamespace, pd.ns, prefixPod, name),
			Name:  name,
			Type:  fuse.DT_Dir,
		})
	}
	return dirs, nil
}

func (pd *podDir) Lookup(ctx context.Context, name string) (fs.Node, error) {
	inode, ok := inoMgr.get(prefixNamespace, pd.ns, prefixPod, name)
	if !ok {
		return nil, fuse.ENOENT
	}

	return newPodEntryDir(inode, pd.ns, name, pd.cli), nil
}

// pod entry represents a pod
type podEntryDir struct {
	inode uint64
	name  string
	ns    string
	cli   *kube.Client
}

func newPodEntryDir(inode uint64, ns, name string, cli *kube.Client) *podEntryDir {
	return &podEntryDir{
		inode: inode,
		ns:    ns,
		name:  name,
		cli:   cli,
	}
}

func (p *podEntryDir) Attr(ctx context.Context, attr *fuse.Attr) error {
	attr.Inode = p.inode
	attr.Mode = os.ModeDir | permDir
	return nil
}

func (p *podEntryDir) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	var dirs []fuse.Dirent

	dirs = append(dirs, fuse.Dirent{
		Inode: inoMgr.getOrCreate(prefixNamespace, p.ns, prefixPod, p.name, "describe"),
		Name:  "describe",
		Type:  fuse.DT_File,
	})
	return dirs, nil
}

func (p *podEntryDir) Lookup(ctx context.Context, name string) (fs.Node, error) {
	inode, ok := inoMgr.get(prefixNamespace, p.ns, prefixPod, p.name, name)
	if !ok {
		return nil, fuse.ENOENT
	}

	return newFile(inode, nil), nil
}
