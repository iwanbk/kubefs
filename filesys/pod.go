package filesys

import (
	"context"
	"os"

	"bazil.org/fuse"
)

// podDir dir represents a pod directory
type podDir struct {
	inode uint64
	name  string
}

func newPodDir(inode uint64, name string) *podDir {
	return &podDir{
		inode: inode,
		name:  name,
	}
}

func (pd *podDir) Attr(ctx context.Context, attr *fuse.Attr) error {
	attr.Inode = pd.inode
	attr.Mode = os.ModeDir | 0444
	return nil
}
