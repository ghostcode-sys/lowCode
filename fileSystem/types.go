package filesystem

import (
	"errors"
	"time"
)

var (
	ErrNotFound          = errors.New("entry not found")
	ErrAlreadyExists     = errors.New("entry already exists")
	ErrNotDirectory      = errors.New("not a directory")
	ErrNotFile           = errors.New("not a file")
	ErrPermission        = errors.New("permission denied")
	ErrInvalidPath       = errors.New("invalid path")
	ErrNoFreeBlocks      = errors.New("no free blocks available")
	ErrNoFreeInodes      = errors.New("no free inodes available")
	ErrDirectoryNotEmpty = errors.New("directory not empty")
	ErrIsDirectory       = errors.New("is a directory")
	ErrReadOnly          = errors.New("filesystem is read-only")
	ErrOutOfRange        = errors.New("offset out of range")
)

type FileType int

const (
	FileTypeRegular FileType = iota
	FileTypeDirectory
	FileTypeSysLink
)

func (ft FileType) String() string {
	switch ft {
	case FileTypeRegular:
		return "FILE"
	case FileTypeDirectory:
		return "DIRECTORY"
	case FileTypeSysLink:
		return "SYMLINK"
	default:
		return "UNKNOWN"
	}
}

type Permission uint16

const (
	PermOwnerRead    Permission = 0400
	PermOwnerWrite   Permission = 0200
	PermOwnerExecute Permission = 0100
	PermGroupRead    Permission = 0040
	PermGroupWrite   Permission = 0020
	PermGroupExecute Permission = 0010
	PermOtherRead    Permission = 0004
	PermOtherWrite   Permission = 0002
	PermOtherExecute Permission = 0001

	// Convenience combos
	PermDefault    Permission = 0644 // rw-r--r--
	PermDirDefault Permission = 0755 // rwxr-xr-x
)

func (p Permission) String() string {
	chars := []byte("---------")
	bits := []struct {
		perm Permission
		pos  int
		ch   byte
	}{
		{PermOwnerRead, 0, 'r'}, {PermOwnerWrite, 1, 'w'}, {PermOwnerExecute, 2, 'x'},
		{PermGroupRead, 3, 'r'}, {PermGroupWrite, 4, 'w'}, {PermGroupExecute, 5, 'x'},
		{PermOtherRead, 6, 'r'}, {PermOtherWrite, 7, 'w'}, {PermOtherExecute, 8, 'x'},
	}
	for _, b := range bits {
		if p&b.perm != 0 {
			chars[b.pos] = b.ch
		}
	}
	return string(chars)
}

type INode struct {
	ID         uint64
	FileType   FileType
	Permission Permission
	Size       int64
	OwnerUID   int
	LinkCount  int
	CreatedAt  time.Time
	ModifiedAt time.Time
	AccessedAt time.Time
	BlockIDs   []int // list of block indices assigned to this inode
}

func newINode(id uint64, ft FileType, perm Permission, uid int) *INode {
	now := time.Now()
	return &INode{
		ID:         id,
		FileType:   ft,
		Permission: perm,
		OwnerUID:   uid,
		LinkCount:  1,
		CreatedAt:  now,
		ModifiedAt: now,
		AccessedAt: now,
		BlockIDs:   []int{},
	}
}

func (n *INode) touch() {
	n.ModifiedAt = time.Now()
	n.AccessedAt = time.Now()
}

// ─── FileSystemEntry interface ────────────────────────────────────────────────

type FileSystemEntry interface {
	Name() string
	INode() *INode
	Parent() *Directory
	SetParent(d *Directory)
	Type() FileType
	Delete() error
	Rename(newName string) error
}

// ─── entryBase — shared fields for File, Directory, SymLink ──────────────────

type entryBase struct {
	name   string
	inode  *INode
	parent *Directory
}

func (e *entryBase) Name() string           { return e.name }
func (e *entryBase) INode() *INode          { return e.inode }
func (e *entryBase) Parent() *Directory     { return e.parent }
func (e *entryBase) SetParent(d *Directory) { e.parent = d }
func (e *entryBase) Type() FileType         { return e.inode.FileType }
