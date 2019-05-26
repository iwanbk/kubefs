package filesys

import (
	"context"
	"os"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"github.com/iwanbk/kubefs/kube"
)

// deploymentDir dir represents a deployments directory
// it contains all deployments in the current namespace
type deploymentDir struct {
	inode uint64
	name  string
	ns    string
	cli   *kube.Client
}

func newDeploymentDir(inode uint64, ns, name string, cli *kube.Client) *deploymentDir {
	return &deploymentDir{
		inode: inode,
		ns:    ns,
		name:  name,
		cli:   cli,
	}
}

func (dd *deploymentDir) Attr(ctx context.Context, attr *fuse.Attr) error {
	attr.Inode = dd.inode
	attr.Mode = os.ModeDir | permDir
	return nil
}

func (dd *deploymentDir) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	names, err := dd.cli.GetDeploymentsName(dd.ns)
	if err != nil {
		return nil, fuse.EIO
	}

	dirs := make([]fuse.Dirent, 0, len(names))
	for _, name := range names {
		dirs = append(dirs, fuse.Dirent{
			Inode: inoMgr.getOrCreate(prefixNamespace, dd.ns, prefixDeployment, name),
			Name:  name,
			Type:  fuse.DT_Dir,
		})
	}
	return dirs, nil
}

func (dd *deploymentDir) Lookup(ctx context.Context, name string) (fs.Node, error) {
	inode, ok := inoMgr.get(prefixNamespace, dd.ns, prefixDeployment, name)
	if !ok {
		return nil, fuse.ENOENT
	}

	return newDeployment(inode, dd.ns, name, dd.cli), nil
}

// deployment entry represents a deployment
type deployment struct {
	inode uint64
	name  string
	ns    string
	cli   *kube.Client
}

func newDeployment(inode uint64, ns, name string, cli *kube.Client) *deployment {
	return &deployment{
		inode: inode,
		ns:    ns,
		name:  name,
		cli:   cli,
	}
}

func (d *deployment) Attr(ctx context.Context, attr *fuse.Attr) error {
	attr.Inode = d.inode
	attr.Mode = os.ModeDir | permDir
	return nil
}

func (d *deployment) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	var dirs []fuse.Dirent

	dirs = append(dirs,
		fuse.Dirent{
			Inode: inoMgr.getOrCreate(prefixNamespace, d.ns, prefixDeployment, d.name, deploymentDescribe),
			Name:  deploymentDescribe,
			Type:  fuse.DT_File,
		},
	)
	return dirs, nil
}

func (d *deployment) Lookup(ctx context.Context, name string) (fs.Node, error) {
	inode, ok := inoMgr.get(prefixNamespace, d.ns, prefixDeployment, d.name, name)
	if !ok {
		return nil, fuse.ENOENT
	}
	switch name {
	case deploymentDescribe:
		return newFile(inode, d.describe), nil
	default:
		return nil, fuse.ENOENT
	}
}

func (d *deployment) describe(ctx context.Context) ([]byte, error) {
	return d.cli.GetDeploymentDescribe(ctx, d.ns, d.name)
}

const (
	deploymentDescribe = "describe"
)
