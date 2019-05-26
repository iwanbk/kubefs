package filesys

import (
	"context"

	"bazil.org/fuse"
)

type file struct {
	inode   uint64
	readAll readAllFn
}

type readAllFn func(context.Context) ([]byte, error)

func newFile(inode uint64, readAll readAllFn) *file {
	return &file{
		inode:   inode,
		readAll: readAll,
	}
}
func (f *file) Attr(ctx context.Context, a *fuse.Attr) error {
	a.Inode = f.inode
	a.Mode = permFile
	a.Size = uint64(10000) // TODO: give proper number, could we?
	return nil
}

func (f *file) ReadAll(ctx context.Context) ([]byte, error) {
	if f.readAll == nil {
		return nil, fuse.ENOTSUP
	}
	return f.readAll(ctx)
}
