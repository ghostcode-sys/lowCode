package filesystem

import (
	"bytes"
	"fmt"
	"testing"
)

// helper: get a fresh FSM for each test
func newFSM(t *testing.T) *FileSystemManager {
	t.Helper()
	ResetInstance()
	return GetInstance(DefaultConfig())
}

// ─── INode ────────────────────────────────────────────────────────────────────

func TestINodeAllocation(t *testing.T) {
	table := NewInodeTable()
	inode, err := table.Allocate(FileTypeRegular, PermDefault, 1000)
	if err != nil {
		t.Fatal(err)
	}
	if inode.FileType != FileTypeRegular {
		t.Errorf("expected FileTypeRegular, got %v", inode.FileType)
	}
	if inode.Permission != PermDefault {
		t.Errorf("expected PermDefault, got %v", inode.Permission)
	}

	got, err := table.Get(inode.ID)
	if err != nil || got.ID != inode.ID {
		t.Errorf("Get failed: %v", err)
	}

	if err := table.Free(inode.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := table.Get(inode.ID); err != ErrNotFound {
		t.Errorf("expected ErrNotFound after Free, got %v", err)
	}
}

// ─── Block / BlockManager ─────────────────────────────────────────────────────

func TestBlockReadWrite(t *testing.T) {
	blk := newBlock(0, 512)
	data := []byte("hello, filesystem!")
	n, err := blk.Write(0, data)
	if err != nil || n != len(data) {
		t.Fatalf("Write error: %v (n=%d)", err, n)
	}

	out, err := blk.Read(0, len(data))
	if err != nil || !bytes.Equal(out, data) {
		t.Errorf("Read mismatch: got %q, want %q, err=%v", out, data, err)
	}
}

func TestBlockManagerAllocFree(t *testing.T) {
	bm := NewBlockManager(10, 512)
	if bm.FreeCount() != 10 {
		t.Fatalf("expected 10 free blocks")
	}

	blk, err := bm.AllocateBlock()
	if err != nil {
		t.Fatal(err)
	}
	if bm.FreeCount() != 9 {
		t.Errorf("expected 9 free after alloc")
	}

	_ = bm.FreeBlock(blk.ID)
	if bm.FreeCount() != 10 {
		t.Errorf("expected 10 free after free")
	}
}

func TestBlockManagerExhausted(t *testing.T) {
	bm := NewBlockManager(2, 512)
	_, _ = bm.AllocateBlock()
	_, _ = bm.AllocateBlock()
	_, err := bm.AllocateBlock()
	if err != ErrNoFreeBlocks {
		t.Errorf("expected ErrNoFreeBlocks, got %v", err)
	}
}

// ─── File ─────────────────────────────────────────────────────────────────────

func TestFileWriteRead(t *testing.T) {
	fsm := newFSM(t)

	f, err := fsm.CreateFile("/hello.txt")
	if err != nil {
		t.Fatal(err)
	}

	payload := []byte("The quick brown fox jumps over the lazy dog")
	if err := f.Write(0, payload); err != nil {
		t.Fatal(err)
	}
	if f.GetSize() != int64(len(payload)) {
		t.Errorf("size mismatch: got %d, want %d", f.GetSize(), len(payload))
	}

	got, err := f.Read(0, len(payload))
	if err != nil || !bytes.Equal(got, payload) {
		t.Errorf("Read mismatch: %q vs %q, err=%v", got, payload, err)
	}
}

func TestFileAppend(t *testing.T) {
	fsm := newFSM(t)
	f, _ := fsm.CreateFile("/append.txt")

	_ = f.Write(0, []byte("Hello"))
	_ = f.Append([]byte(", World!"))

	want := []byte("Hello, World!")
	got, _ := f.Read(0, int(f.GetSize()))
	if !bytes.Equal(got, want) {
		t.Errorf("Append: got %q, want %q", got, want)
	}
}

func TestFileTruncate(t *testing.T) {
	fsm := newFSM(t)
	f, _ := fsm.CreateFile("/trunc.txt")

	_ = f.Write(0, bytes.Repeat([]byte("A"), 8192)) // 2 blocks
	if f.BlockCount() < 2 {
		t.Skip("not enough blocks allocated for truncate test")
	}

	if err := f.Truncate(100); err != nil {
		t.Fatal(err)
	}
	if f.GetSize() != 100 {
		t.Errorf("expected size 100, got %d", f.GetSize())
	}
}

func TestFileCrossBlockRead(t *testing.T) {
	ResetInstance()
	cfg := DefaultConfig()
	cfg.BlockSize = 16 // tiny blocks for easy testing
	fsm := GetInstance(cfg)

	f, _ := fsm.CreateFile("/cross.txt")
	data := []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123") // 30 bytes → ~2 blocks
	_ = f.Write(0, data)

	got, err := f.Read(0, len(data))
	if err != nil || !bytes.Equal(got, data) {
		t.Errorf("cross-block read: got %q, want %q, err=%v", got, data, err)
	}
}

// ─── Directory ────────────────────────────────────────────────────────────────

func TestDirectoryAddRemove(t *testing.T) {
	fsm := newFSM(t)

	dir, err := fsm.CreateDirectory("/docs")
	if err != nil {
		t.Fatal(err)
	}

	f, _ := fsm.CreateFile("/docs/readme.txt")
	if f.Parent() != dir {
		t.Error("parent not set correctly")
	}

	entries := dir.ListEntries()
	if len(entries) != 1 {
		t.Errorf("expected 1 entry, got %d", len(entries))
	}

	if err := dir.RemoveEntry("readme.txt"); err != nil {
		t.Fatal(err)
	}
	if !dir.IsEmpty() {
		t.Error("expected directory to be empty after removal")
	}
}

func TestDirectoryNestedSearch(t *testing.T) {
	fsm := newFSM(t)
	_, _ = fsm.CreateDirectory("/a")
	_, _ = fsm.CreateDirectory("/a/b")
	_, _ = fsm.CreateFile("/a/b/target.txt")

	found, err := fsm.root.Search("target.txt")
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}
	if found.Name() != "target.txt" {
		t.Errorf("expected target.txt, got %s", found.Name())
	}
}

func TestDirectoryDelete(t *testing.T) {
	fsm := newFSM(t)
	_, _ = fsm.CreateDirectory("/todelete")
	_, _ = fsm.CreateFile("/todelete/file1.txt")
	_, _ = fsm.CreateFile("/todelete/file2.txt")

	if err := fsm.Delete("/todelete"); err != nil {
		t.Fatal(err)
	}

	if _, err := fsm.Lookup("/todelete"); err != ErrNotFound {
		t.Errorf("expected ErrNotFound after delete, got %v", err)
	}
}

// ─── SymLink ──────────────────────────────────────────────────────────────────

func TestSymLinkResolve(t *testing.T) {
	fsm := newFSM(t)
	f, _ := fsm.CreateFile("/original.txt")
	_ = f.Write(0, []byte("data"))

	link, err := fsm.CreateSymLink("/link.txt", "/original.txt")
	if err != nil {
		t.Fatal(err)
	}
	if link.GetTarget() != "/original.txt" {
		t.Errorf("unexpected target: %s", link.GetTarget())
	}

	resolved, err := link.Resolve()
	if err != nil {
		t.Fatalf("Resolve failed: %v", err)
	}
	if resolved.Name() != "original.txt" {
		t.Errorf("expected original.txt, got %s", resolved.Name())
	}
}

// ─── FileSystemManager ────────────────────────────────────────────────────────

func TestFSManagerSingleton(t *testing.T) {
	ResetInstance()
	a := GetInstance(DefaultConfig())
	b := GetInstance(DefaultConfig())
	if a != b {
		t.Error("GetInstance should return the same pointer")
	}
}

func TestFSManagerLookup(t *testing.T) {
	fsm := newFSM(t)
	_, _ = fsm.CreateDirectory("/usr")
	_, _ = fsm.CreateDirectory("/usr/local")
	_, _ = fsm.CreateFile("/usr/local/bin.txt")

	entry, err := fsm.Lookup("/usr/local/bin.txt")
	if err != nil || entry.Name() != "bin.txt" {
		t.Errorf("Lookup failed: %v", err)
	}
}

func TestFSManagerMove(t *testing.T) {
	fsm := newFSM(t)
	_, _ = fsm.CreateDirectory("/src")
	_, _ = fsm.CreateDirectory("/dst")
	_, _ = fsm.CreateFile("/src/move_me.txt")

	if err := fsm.Move("/src/move_me.txt", "/dst/move_me.txt"); err != nil {
		t.Fatal(err)
	}

	if _, err := fsm.Lookup("/src/move_me.txt"); err != ErrNotFound {
		t.Error("expected src entry to be gone")
	}
	if _, err := fsm.Lookup("/dst/move_me.txt"); err != nil {
		t.Errorf("expected dst entry to exist: %v", err)
	}
}

func TestFSManagerStats(t *testing.T) {
	fsm := newFSM(t)
	stats := fsm.GetStats()
	if stats.TotalBlocks != 1024 {
		t.Errorf("expected 1024 total blocks, got %d", stats.TotalBlocks)
	}
	if stats.FreeBlocks != 1024 {
		t.Errorf("expected 1024 free blocks initially, got %d", stats.FreeBlocks)
	}
}

// ─── Permission string ────────────────────────────────────────────────────────

func TestPermissionString(t *testing.T) {
	cases := []struct {
		perm Permission
		want string
	}{
		{PermDefault, "rw-r--r--"},
		{PermDirDefault, "rwxr-xr-x"},
		{0777, "rwxrwxrwx"},
		{0, "---------"},
	}
	for _, c := range cases {
		got := c.perm.String()
		if got != c.want {
			t.Errorf("Permission(%04o).String() = %q, want %q", c.perm, got, c.want)
		}
	}
}

// ─── Stress / concurrent writes ───────────────────────────────────────────────

func TestConcurrentWrites(t *testing.T) {
	ResetInstance()
	cfg := DefaultConfig()
	cfg.TotalBlocks = 4096
	fsm := GetInstance(cfg)

	const numFiles = 20
	done := make(chan error, numFiles)

	for i := 0; i < numFiles; i++ {
		go func(id int) {
			path := fmt.Sprintf("/concurrent_%d.txt", id)
			f, err := fsm.CreateFile(path)
			if err != nil {
				done <- err
				return
			}
			done <- f.Write(0, []byte(fmt.Sprintf("file %d content", id)))
		}(i)
	}

	for i := 0; i < numFiles; i++ {
		if err := <-done; err != nil {
			t.Errorf("concurrent write error: %v", err)
		}
	}

	entries := fsm.root.ListEntries()
	if len(entries) != numFiles {
		t.Errorf("expected %d root entries, got %d", numFiles, len(entries))
	}
}
