package filesystem

import (
	"sync"
	"sync/atomic"
)

// ─── InodeTable ───────────────────────────────────────────────────────────────

// InodeTable is a central registry of all live inodes.
type InodeTable struct {
	mu     sync.RWMutex
	inodes map[uint64]*INode
	nextID uint64 // used with atomic ops
}

// NewInodeTable creates an empty inode table.
func NewInodeTable() *InodeTable {
	return &InodeTable{
		inodes: make(map[uint64]*INode),
		nextID: 1, // inode 0 is reserved / null
	}
}

// Allocate creates and registers a new INode.
func (t *InodeTable) Allocate(ft FileType, perm Permission, uid int) (*INode, error) {
	id := atomic.AddUint64(&t.nextID, 1) - 1
	inode := newINode(id, ft, perm, uid)

	t.mu.Lock()
	defer t.mu.Unlock()
	t.inodes[id] = inode
	return inode, nil
}

// Get returns the INode for a given id.
func (t *InodeTable) Get(id uint64) (*INode, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	inode, ok := t.inodes[id]
	if !ok {
		return nil, ErrNotFound
	}
	return inode, nil
}

// Free removes an inode from the table.
func (t *InodeTable) Free(id uint64) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if _, ok := t.inodes[id]; !ok {
		return ErrNotFound
	}
	delete(t.inodes, id)
	return nil
}

// Update replaces an inode entry (e.g. after size change).
func (t *InodeTable) Update(inode *INode) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.inodes[inode.ID] = inode
}

// Count returns how many inodes are currently allocated.
func (t *InodeTable) Count() int {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return len(t.inodes)
}
