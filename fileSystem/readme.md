# LLD (Low Level Design) patterns used:

## 1. 🏭 Singleton Pattern
**Where:** `FileSystemManager`
```go
var (
    instance *FileSystemManager
    once     sync.Once
)
func GetInstance(cfg Config) *FileSystemManager {
    once.Do(func() { ... })
    return instance
}
```
Ensures only one filesystem manager exists across the entire application lifetime.

---

## 2. 🧩 Composite Pattern
**Where:** `Directory` + `FileSystemEntry`
```go
type Directory struct {
    children map[string]FileSystemEntry  // holds Files, Dirs, SymLinks uniformly
}
```
Both `File` and `Directory` implement `FileSystemEntry`, so a directory can contain directories (tree structure). Classic composite — treat leaf and container uniformly.

---

## 3. 🎭 Facade Pattern
**Where:** `FileSystemManager`
```go
fsm.CreateFile("/a/b/c.txt")   // hides INode alloc + Block alloc + Dir traversal
fsm.Delete("/a/b/c.txt")       // hides recursive cleanup + inode free + block free
fsm.Move("/src", "/dst")       // hides detach + re-attach logic
```
The client never touches `BlockManager`, `InodeTable`, or raw `INode` directly — `FileSystemManager` hides all that complexity.

---

## 4. 🪆 Template Method / Inheritance Pattern
**Where:** `entryBase` → `File`, `Directory`, `SymLink`
```go
type entryBase struct { name, inode, parent }  // shared fields + shared methods

type File struct { entryBase; blocks []Block }         // extends + specialises
type Directory struct { entryBase; children map... }   // extends + specialises
type SymLink struct { entryBase; targetPath string }   // extends + specialises
```
Common behaviour (`Name()`, `INode()`, `Parent()`) lives in the base. Each subtype overrides only what differs (`Delete()`, `Rename()`).

---

## 5. 🔌 Strategy / Interface Pattern
**Where:** `FileSystemEntry` interface
```go
type FileSystemEntry interface {
    Name() string
    Delete() error
    Rename(newName string) error
    ...
}
```
`FileSystemManager.Delete(path)` calls `entry.Delete()` without knowing if it's a `File`, `Directory`, or `SymLink`. The correct delete behaviour is resolved at runtime — classic polymorphic dispatch.

---

## 6. 🏊 Object Pool Pattern
**Where:** `BlockManager`
```go
type BlockManager struct {
    blocks  []*Block   // pre-allocated pool
    freeMap []bool     // bitmap tracking free/used
}
func (bm *BlockManager) AllocateBlock() (*Block, error) { ... } // borrow
func (bm *BlockManager) FreeBlock(id int) error { ... }         // return
```
Blocks are pre-allocated at startup and reused rather than created/destroyed, avoiding fragmentation — exactly how real filesystems (ext4, NTFS) manage disk blocks.

---

## 7. 🗂️ Registry Pattern
**Where:** `InodeTable`
```go
type InodeTable struct {
    inodes map[uint64]*INode  // central registry
    nextID uint64             // atomic ID generator
}
func (t *InodeTable) Allocate(...) (*INode, error) { ... }
func (t *InodeTable) Get(id uint64) (*INode, error) { ... }
func (t *InodeTable) Free(id uint64) error { ... }
```
Acts as a centralised lookup table — every inode in the system is registered here, just like a real OS inode table.

---

## 8. 🔗 Proxy Pattern
**Where:** `SymLink`
```go
func (s *SymLink) Resolve() (FileSystemEntry, error) {
    return s.fs.Lookup(s.targetPath)  // delegates to the real target
}
```
`SymLink` acts as a proxy — it looks like any other `FileSystemEntry` but transparently forwards operations to its target.

---

## 9. 🔒 Monitor / Thread-Safety Pattern
**Where:** `File`, `Directory`, `BlockManager`, `InodeTable`
```go
type File      struct { mu sync.RWMutex; ... }
type Directory struct { mu sync.RWMutex; ... }
type BlockManager struct { mu sync.Mutex; ... }
type InodeTable   struct { mu sync.RWMutex; ... }
```
Every shared resource guards itself with a mutex — readers acquire `RLock`, writers acquire `Lock`. This is the **Monitor pattern** applied per-object.

---

## Summary Table

| Pattern | Class(es) | Purpose |
|---|---|---|
| Singleton | `FileSystemManager` | One global FS instance |
| Composite | `Directory` + `FileSystemEntry` | Recursive tree of entries |
| Facade | `FileSystemManager` | Simple API over complex subsystems |
| Template Method | `entryBase` + subtypes | Shared base, specialised behaviour |
| Strategy/Interface | `FileSystemEntry` | Polymorphic delete/rename/etc |
| Object Pool | `BlockManager` | Reuse pre-allocated blocks |
| Registry | `InodeTable` | Central inode lookup & lifecycle |
| Proxy | `SymLink` | Transparent forwarding to target |
| Monitor | All shared structs | Thread-safe concurrent access |