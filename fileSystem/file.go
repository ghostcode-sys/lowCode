package filesystem

import (
	"sync"
	"time"
)

// ─── File ─────────────────────────────────────────────────────────────────────

// File represents a regular file in the filesystem.
type File struct {
	entryBase
	mu           sync.RWMutex
	blockManager *BlockManager
	blocks       []*Block // ordered list of blocks holding the file's data
}

func newFile(name string, inode *INode, bm *BlockManager) *File {
	return &File{
		entryBase:    entryBase{name: name, inode: inode},
		blockManager: bm,
		blocks:       []*Block{},
	}
}

// ── Read ──────────────────────────────────────────────────────────────────────

// Read reads `length` bytes starting at `offset`.
func (f *File) Read(offset, length int) ([]byte, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	size := int(f.inode.Size)
	if offset >= size || length == 0 {
		return []byte{}, nil
	}
	if offset+length > size {
		length = size - offset
	}

	bs := f.blockManager.BlockSize()
	result := make([]byte, 0, length)
	remaining := length
	pos := offset

	for remaining > 0 && pos < size {
		blockIdx := pos / bs
		blockOff := pos % bs

		if blockIdx >= len(f.blocks) {
			break
		}

		toRead := bs - blockOff
		if toRead > remaining {
			toRead = remaining
		}

		chunk, err := f.blocks[blockIdx].Read(blockOff, toRead)
		if err != nil {
			return nil, err
		}
		result = append(result, chunk...)
		pos += toRead
		remaining -= toRead
	}

	f.inode.AccessedAt = time.Now() // update access time (no mutation of data)
	return result, nil
}

// ── Write ─────────────────────────────────────────────────────────────────────

// Write writes data starting at offset, expanding the file if necessary.
func (f *File) Write(offset int, data []byte) error {
	if len(data) == 0 {
		return nil
	}

	f.mu.Lock()
	defer f.mu.Unlock()

	bs := f.blockManager.BlockSize()
	end := offset + len(data)

	// Ensure we have enough blocks to cover [0, end)
	neededBlocks := (end + bs - 1) / bs
	for len(f.blocks) < neededBlocks {
		blk, err := f.blockManager.AllocateBlock()
		if err != nil {
			return err
		}
		f.blocks = append(f.blocks, blk)
		f.inode.BlockIDs = append(f.inode.BlockIDs, blk.ID)
	}

	// Write data chunk by chunk into the correct blocks
	written := 0
	pos := offset
	for written < len(data) {
		blockIdx := pos / bs
		blockOff := pos % bs
		canWrite := bs - blockOff
		if canWrite > len(data)-written {
			canWrite = len(data) - written
		}
		n, err := f.blocks[blockIdx].Write(blockOff, data[written:written+canWrite])
		if err != nil {
			return err
		}
		written += n
		pos += n
	}

	if int64(end) > f.inode.Size {
		f.inode.Size = int64(end)
	}
	f.inode.touch()
	return nil
}

// ── Append ────────────────────────────────────────────────────────────────────

// Append adds data to the end of the file.
func (f *File) Append(data []byte) error {
	return f.Write(int(f.inode.Size), data)
}

// ── Truncate ──────────────────────────────────────────────────────────────────

// Truncate shrinks or grows the file to exactly `size` bytes.
func (f *File) Truncate(size int) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	bs := f.blockManager.BlockSize()
	neededBlocks := (size + bs - 1) / bs
	if size == 0 {
		neededBlocks = 0
	}

	// Free excess blocks
	for len(f.blocks) > neededBlocks {
		last := f.blocks[len(f.blocks)-1]
		_ = f.blockManager.FreeBlock(last.ID)
		f.blocks = f.blocks[:len(f.blocks)-1]
		f.inode.BlockIDs = f.inode.BlockIDs[:len(f.inode.BlockIDs)-1]
	}

	// Allocate missing blocks if growing
	for len(f.blocks) < neededBlocks {
		blk, err := f.blockManager.AllocateBlock()
		if err != nil {
			return err
		}
		f.blocks = append(f.blocks, blk)
		f.inode.BlockIDs = append(f.inode.BlockIDs, blk.ID)
	}

	f.inode.Size = int64(size)
	f.inode.touch()
	return nil
}

// ── Delete ────────────────────────────────────────────────────────────────────

// Delete frees all blocks and removes the file from its parent directory.
func (f *File) Delete() error {
	f.mu.Lock()
	defer f.mu.Unlock()

	for _, blk := range f.blocks {
		_ = f.blockManager.FreeBlock(blk.ID)
	}
	f.blocks = nil
	f.inode.BlockIDs = nil
	f.inode.Size = 0

	if f.parent != nil {
		f.parent.mu.Lock()
		defer f.parent.mu.Unlock()
		delete(f.parent.children, f.name)
	}
	return nil
}

// Rename renames the file within its parent directory.
func (f *File) Rename(newName string) error {
	if newName == "" {
		return ErrInvalidPath
	}
	if f.parent == nil {
		f.name = newName
		return nil
	}
	f.parent.mu.Lock()
	defer f.parent.mu.Unlock()
	if _, exists := f.parent.children[newName]; exists {
		return ErrAlreadyExists
	}
	delete(f.parent.children, f.name)
	f.name = newName
	f.parent.children[newName] = f
	return nil
}

// GetSize returns the file size in bytes.
func (f *File) GetSize() int64 { return f.inode.Size }

// BlockCount returns the number of blocks allocated to this file.
func (f *File) BlockCount() int {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return len(f.blocks)
}
