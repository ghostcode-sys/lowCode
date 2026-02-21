package filesystem

import "sync"

const DefaulBlockSize = 4096

type Block struct {
	ID        int
	data      []byte
	blockSize int
}

func newBlock(id, size int) *Block {
	return &Block{
		ID:        id,
		data:      make([]byte, size),
		blockSize: size,
	}
}

func (b *Block) Read(offset, length int) ([]byte, error) {
	if offset < 0 || offset > b.blockSize {
		return nil, ErrOutOfRange
	}

	end := offset + length

	if end > b.blockSize {
		end = b.blockSize
	}

	out := make([]byte, length)

	copy(out, b.data[offset:end])
	return out, nil
}

func (b *Block) Write(offset int, data []byte) (int, error) {
	if offset < 0 || offset > b.blockSize {
		return 0, ErrOutOfRange
	}
	n := copy(b.data[offset:], data)
	return n, nil
}

func (b *Block) Zero() {
	for i := range b.data {
		b.data[i] = 0
	}
}

func (b *Block) Size() int { return b.blockSize }

type BlockManager struct {
	mu        sync.Mutex
	blocks    []*Block
	freeMap   []bool
	blockSize int
	total     int
}

func NewBlockManager(total, blockSize int) *BlockManager {
	blocks := make([]*Block, total)
	free := make([]bool, total)

	for i := range total {
		blocks[i] = newBlock(i, blockSize)
		free[i] = true
	}

	return &BlockManager{
		blocks:    blocks,
		freeMap:   free,
		blockSize: blockSize,
		total:     total,
	}
}

func (bm *BlockManager) AllocateBlock() (*Block, error) {
	bm.mu.Lock()
	defer bm.mu.Unlock()

	for i, free := range bm.freeMap {
		if free {
			bm.freeMap[i] = false
			bm.blocks[i].Zero()
			return bm.blocks[i], nil
		}
	}
	return nil, ErrNoFreeBlocks
}

// AllocateN allocates n contiguous (or scattered) blocks atomically.
func (bm *BlockManager) AllocateN(n int) ([]*Block, error) {
	bm.mu.Lock()
	defer bm.mu.Unlock()

	indices := make([]int, 0, n)
	for i, free := range bm.freeMap {
		if free {
			indices = append(indices, i)
			if len(indices) == n {
				break
			}
		}
	}
	if len(indices) < n {
		return nil, ErrNoFreeBlocks
	}
	result := make([]*Block, n)
	for i, idx := range indices {
		bm.freeMap[idx] = false
		bm.blocks[idx].Zero()
		result[i] = bm.blocks[idx]
	}
	return result, nil
}

// FreeBlock returns a block back to the free pool.
func (bm *BlockManager) FreeBlock(id int) error {
	bm.mu.Lock()
	defer bm.mu.Unlock()

	if id < 0 || id >= bm.total {
		return ErrOutOfRange
	}
	bm.freeMap[id] = true
	return nil
}

// GetBlock returns the block for a given id.
func (bm *BlockManager) GetBlock(id int) (*Block, error) {
	if id < 0 || id >= bm.total {
		return nil, ErrOutOfRange
	}
	return bm.blocks[id], nil
}

// FreeCount returns the number of available blocks.
func (bm *BlockManager) FreeCount() int {
	bm.mu.Lock()
	defer bm.mu.Unlock()
	count := 0
	for _, f := range bm.freeMap {
		if f {
			count++
		}
	}
	return count
}

// TotalCount returns total blocks managed.
func (bm *BlockManager) TotalCount() int { return bm.total }

// BlockSize returns the configured block size.
func (bm *BlockManager) BlockSize() int { return bm.blockSize }
