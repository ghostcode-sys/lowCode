package filesystem

import (
	"sync"
)

// ─── FileSystemManager ────────────────────────────────────────────────────────

// FileSystemManager is the top-level façade for all filesystem operations.
// It is a singleton; use GetInstance() to obtain the single instance.
type FileSystemManager struct {
	mu           sync.Mutex
	root         *Directory
	blockManager *BlockManager
	inodeTable   *InodeTable
	ownerUID     int // default UID for new entries
}

var (
	instance *FileSystemManager
	once     sync.Once
)

// Config holds initialisation parameters.
type Config struct {
	TotalBlocks int
	BlockSize   int
	OwnerUID    int
}

// DefaultConfig provides sensible defaults: 1024 × 4 KB blocks.
func DefaultConfig() Config {
	return Config{
		TotalBlocks: 1024,
		BlockSize:   DefaulBlockSize,
		OwnerUID:    1000,
	}
}

// GetInstance returns (or initialises) the singleton FileSystemManager.
// Subsequent calls ignore cfg.
func GetInstance(cfg Config) *FileSystemManager {
	once.Do(func() {
		bm := NewBlockManager(cfg.TotalBlocks, cfg.BlockSize)
		it := NewInodeTable()

		rootInode, _ := it.Allocate(FileTypeDirectory, PermDirDefault, cfg.OwnerUID)
		root := newDirectory("/", rootInode)

		instance = &FileSystemManager{
			root:         root,
			blockManager: bm,
			inodeTable:   it,
			ownerUID:     cfg.OwnerUID,
		}
	})
	return instance
}

// ResetInstance tears down the singleton (useful in tests).
func ResetInstance() {
	once = sync.Once{}
	instance = nil
}

// ─── Internal helpers ─────────────────────────────────────────────────────────

// resolveDirAndBase resolves a path to its parent directory and the final component.
// e.g. "/a/b/c" → (dirFor("/a/b"), "c")
func (fsm *FileSystemManager) resolveDirAndBase(path string) (*Directory, string, error) {
	parts := splitPath(path)
	if len(parts) == 0 {
		return nil, "", ErrInvalidPath
	}

	cur := fsm.root
	for _, part := range parts[:len(parts)-1] {
		entry, err := cur.GetEntry(part)
		if err != nil {
			return nil, "", err
		}
		dir, ok := entry.(*Directory)
		if !ok {
			return nil, "", ErrNotDirectory
		}
		cur = dir
	}
	return cur, parts[len(parts)-1], nil
}

// Lookup resolves a full absolute path to a FileSystemEntry.
func (fsm *FileSystemManager) Lookup(path string) (FileSystemEntry, error) {
	parts := splitPath(path)
	if len(parts) == 0 {
		return fsm.root, nil
	}

	var cur FileSystemEntry = fsm.root
	for _, part := range parts {
		dir, ok := cur.(*Directory)
		if !ok {
			return nil, ErrNotDirectory
		}
		next, err := dir.GetEntry(part)
		if err != nil {
			return nil, err
		}
		cur = next
	}
	return cur, nil
}

// ─── Public API ───────────────────────────────────────────────────────────────

// CreateFile creates a new regular file at the given absolute path.
func (fsm *FileSystemManager) CreateFile(path string) (*File, error) {
	fsm.mu.Lock()
	defer fsm.mu.Unlock()

	dir, name, err := fsm.resolveDirAndBase(path)
	if err != nil {
		return nil, err
	}

	inode, err := fsm.inodeTable.Allocate(FileTypeRegular, PermDefault, fsm.ownerUID)
	if err != nil {
		return nil, err
	}

	f := newFile(name, inode, fsm.blockManager)
	if err := dir.AddEntry(f); err != nil {
		_ = fsm.inodeTable.Free(inode.ID)
		return nil, err
	}
	return f, nil
}

// CreateDirectory creates a new directory at the given absolute path.
func (fsm *FileSystemManager) CreateDirectory(path string) (*Directory, error) {
	fsm.mu.Lock()
	defer fsm.mu.Unlock()

	dir, name, err := fsm.resolveDirAndBase(path)
	if err != nil {
		return nil, err
	}

	inode, err := fsm.inodeTable.Allocate(FileTypeDirectory, PermDirDefault, fsm.ownerUID)
	if err != nil {
		return nil, err
	}

	newDir := newDirectory(name, inode)
	if err := dir.AddEntry(newDir); err != nil {
		_ = fsm.inodeTable.Free(inode.ID)
		return nil, err
	}
	return newDir, nil
}

// CreateSymLink creates a symbolic link at `linkPath` pointing to `targetPath`.
func (fsm *FileSystemManager) CreateSymLink(linkPath, targetPath string) (*SymLink, error) {
	fsm.mu.Lock()
	defer fsm.mu.Unlock()

	dir, name, err := fsm.resolveDirAndBase(linkPath)
	if err != nil {
		return nil, err
	}

	inode, err := fsm.inodeTable.Allocate(FileTypeSysLink, PermDefault, fsm.ownerUID)
	if err != nil {
		return nil, err
	}

	link := newSymLink(name, targetPath, inode, fsm)
	if err := dir.AddEntry(link); err != nil {
		_ = fsm.inodeTable.Free(inode.ID)
		return nil, err
	}
	return link, nil
}

// Delete removes the entry at path (recursively for directories).
func (fsm *FileSystemManager) Delete(path string) error {
	fsm.mu.Lock()
	defer fsm.mu.Unlock()

	entry, err := fsm.Lookup(path)
	if err != nil {
		return err
	}
	inodeID := entry.INode().ID
	if err := entry.Delete(); err != nil {
		return err
	}
	return fsm.inodeTable.Free(inodeID)
}

// Move moves an entry from src to dst path (rename/move semantics).
func (fsm *FileSystemManager) Move(srcPath, dstPath string) error {
	fsm.mu.Lock()
	defer fsm.mu.Unlock()

	src, err := fsm.Lookup(srcPath)
	if err != nil {
		return err
	}

	dstDir, dstName, err := fsm.resolveDirAndBase(dstPath)
	if err != nil {
		return err
	}

	// Detach from current parent
	if src.Parent() != nil {
		if err := src.Parent().RemoveEntry(src.Name()); err != nil {
			return err
		}
	}

	// Change name and attach to new parent
	switch e := src.(type) {
	case *File:
		e.name = dstName
	case *Directory:
		e.name = dstName
	case *SymLink:
		e.name = dstName
	}

	return dstDir.AddEntry(src)
}

// Rename renames the last component of path to newName (same directory).
func (fsm *FileSystemManager) Rename(path, newName string) error {
	entry, err := fsm.Lookup(path)
	if err != nil {
		return err
	}
	return entry.Rename(newName)
}

// ─── Stats ────────────────────────────────────────────────────────────────────

// Stats returns a snapshot of filesystem utilisation.
type Stats struct {
	TotalBlocks int
	FreeBlocks  int
	UsedBlocks  int
	TotalInodes int
	BlockSize   int
}

// GetStats returns current filesystem statistics.
func (fsm *FileSystemManager) GetStats() Stats {
	free := fsm.blockManager.FreeCount()
	total := fsm.blockManager.TotalCount()
	return Stats{
		TotalBlocks: total,
		FreeBlocks:  free,
		UsedBlocks:  total - free,
		TotalInodes: fsm.inodeTable.Count(),
		BlockSize:   fsm.blockManager.BlockSize(),
	}
}

// Root returns the root directory.
func (fsm *FileSystemManager) Root() *Directory { return fsm.root }
