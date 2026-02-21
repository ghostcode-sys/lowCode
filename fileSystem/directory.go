package filesystem

import (
	"strings"
	"sync"
)

// ─── Directory ────────────────────────────────────────────────────────────────

// Directory represents a filesystem directory that holds named children.
type Directory struct {
	entryBase
	mu       sync.RWMutex
	children map[string]FileSystemEntry
}

func newDirectory(name string, inode *INode) *Directory {
	return &Directory{
		entryBase: entryBase{name: name, inode: inode},
		children:  make(map[string]FileSystemEntry),
	}
}

// AddEntry adds a child entry (file, dir, or symlink) to this directory.
func (d *Directory) AddEntry(entry FileSystemEntry) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if _, exists := d.children[entry.Name()]; exists {
		return ErrAlreadyExists
	}
	d.children[entry.Name()] = entry
	entry.SetParent(d)
	d.inode.touch()
	return nil
}

// RemoveEntry removes a named child from this directory.
func (d *Directory) RemoveEntry(name string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if _, exists := d.children[name]; !exists {
		return ErrNotFound
	}
	delete(d.children, name)
	d.inode.touch()
	return nil
}

// GetEntry returns the child entry with the given name.
func (d *Directory) GetEntry(name string) (FileSystemEntry, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	entry, exists := d.children[name]
	if !exists {
		return nil, ErrNotFound
	}
	return entry, nil
}

// ListEntries returns all children as a slice.
func (d *Directory) ListEntries() []FileSystemEntry {
	d.mu.RLock()
	defer d.mu.RUnlock()

	result := make([]FileSystemEntry, 0, len(d.children))
	for _, e := range d.children {
		result = append(result, e)
	}
	return result
}

// IsEmpty returns true if the directory has no children.
func (d *Directory) IsEmpty() bool {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return len(d.children) == 0
}

// Search does a recursive depth-first search for an entry by name.
func (d *Directory) Search(name string) (FileSystemEntry, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if entry, ok := d.children[name]; ok {
		return entry, nil
	}
	for _, child := range d.children {
		if sub, ok := child.(*Directory); ok {
			if found, err := sub.Search(name); err == nil {
				return found, nil
			}
		}
	}
	return nil, ErrNotFound
}

// Delete removes all children recursively and then removes self from parent.
func (d *Directory) Delete() error {
	d.mu.Lock()

	// Collect children before unlocking to avoid nested lock
	children := make([]FileSystemEntry, 0, len(d.children))
	for _, c := range d.children {
		children = append(children, c)
	}
	d.mu.Unlock()

	for _, child := range children {
		if err := child.Delete(); err != nil {
			return err
		}
	}

	d.mu.Lock()
	d.children = make(map[string]FileSystemEntry)
	d.mu.Unlock()

	if d.parent != nil {
		d.parent.mu.Lock()
		defer d.parent.mu.Unlock()
		delete(d.parent.children, d.name)
	}
	return nil
}

// Rename renames this directory within its parent.
func (d *Directory) Rename(newName string) error {
	if newName == "" {
		return ErrInvalidPath
	}
	if d.parent == nil {
		d.name = newName
		return nil
	}
	d.parent.mu.Lock()
	defer d.parent.mu.Unlock()

	if _, exists := d.parent.children[newName]; exists {
		return ErrAlreadyExists
	}
	delete(d.parent.children, d.name)
	d.name = newName
	d.parent.children[newName] = d
	return nil
}

// ─── SymLink ──────────────────────────────────────────────────────────────────

// SymLink is a symbolic link that points to another path.
type SymLink struct {
	entryBase
	targetPath string
	fs         *FileSystemManager // back-reference to resolve the target
}

func newSymLink(name, target string, inode *INode, fs *FileSystemManager) *SymLink {
	return &SymLink{
		entryBase:  entryBase{name: name, inode: inode},
		targetPath: target,
		fs:         fs,
	}
}

// Resolve follows the link and returns the target entry.
func (s *SymLink) Resolve() (FileSystemEntry, error) {
	return s.fs.Lookup(s.targetPath)
}

// GetTarget returns the raw target path string.
func (s *SymLink) GetTarget() string { return s.targetPath }

// Delete removes the symlink from its parent.
func (s *SymLink) Delete() error {
	if s.parent != nil {
		s.parent.mu.Lock()
		defer s.parent.mu.Unlock()
		delete(s.parent.children, s.name)
	}
	return nil
}

// Rename renames the symlink within its parent.
func (s *SymLink) Rename(newName string) error {
	if newName == "" {
		return ErrInvalidPath
	}
	if s.parent == nil {
		s.name = newName
		return nil
	}
	s.parent.mu.Lock()
	defer s.parent.mu.Unlock()

	if _, exists := s.parent.children[newName]; exists {
		return ErrAlreadyExists
	}
	delete(s.parent.children, s.name)
	s.name = newName
	s.parent.children[newName] = s
	return nil
}

// ─── Path helpers ─────────────────────────────────────────────────────────────

// splitPath splits "/a/b/c" into ["a","b","c"], ignoring empty segments.
func splitPath(path string) []string {
	parts := strings.Split(path, "/")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}
