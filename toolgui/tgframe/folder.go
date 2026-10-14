package tgframe

import (
	"io"
	"io/fs"
	"sync"

	"github.com/voilelab/toolgui/toolgui/tgutil"
)

// Folder is a user-authorized directory. Names follow fs.ValidPath: slash-separated,
// relative to the root, with no .. components. Implementations must confine symlinks too.
type Folder interface {
	fs.ReadDirFS
	fs.StatFS
	// Create creates or truncates a file; its parent must exist.
	// Read-only grants return ErrFolderReadOnly.
	Create(name string) (io.WriteCloser, error)
	// MkdirAll creates a directory and missing parents.
	// Read-only grants return ErrFolderReadOnly.
	MkdirAll(name string) error
	// Writable reports whether the user granted write access.
	Writable() bool
	// Path is the absolute OS path on desktop, and empty elsewhere.
	Path() string
}

// ErrFolderReadOnly is returned when a read-only grant is asked to write.
var ErrFolderReadOnly = tgutil.NewError("folder is read-only")

type folderStore struct {
	sync.RWMutex
	values map[string]Folder
	closed bool
}

func (f *folderStore) remove(id string) {
	f.Lock()
	defer f.Unlock()
	closeFolder(f.values[id])
	delete(f.values, id)
}
func closeFolder(f Folder) {
	if c, ok := f.(io.Closer); ok {
		_ = c.Close()
	}
}
func (f *folderStore) destroy() {
	f.Lock()
	defer f.Unlock()
	f.closed = true
	for id, v := range f.values {
		closeFolder(v)
		delete(f.values, id)
	}
}

// GetFolder returns only a trusted executor grant, never a JSON input value.
func (s *State) GetFolder(id string) Folder {
	s.folders.RLock()
	defer s.folders.RUnlock()
	return s.folders.values[id]
}

// EventFolder is an executor-only event; ParseEvent cannot construct it.
// HandleFolder takes ownership of Folder, including closing rejected grants.
type EventFolder struct {
	ID     string
	Folder Folder
}

func (e *EventFolder) validateState(s *State) error { return validateComponentID(s, e.ID) }
func (e *EventFolder) ApplyState(s *State) {
	s.folders.Lock()
	defer s.folders.Unlock()
	if s.folders.closed {
		closeFolder(e.Folder)
		return
	}
	closeFolder(s.folders.values[e.ID])
	s.folders.values[e.ID] = e.Folder
}

// ValidateFolderPick checks the target before opening a blocking dialog.
func (s *Session) ValidateFolderPick(id string) error {
	s.handling.Lock()
	defer s.handling.Unlock()
	if s.closed.Load() {
		return ErrNoSession
	}
	return validateEvent(&EventFolder{ID: id}, s.state)
}

// HandleFolder accepts a grant and starts the run that reads it.
func (s *Session) HandleFolder(e *EventFolder) error {
	s.handling.Lock()
	defer s.handling.Unlock()
	if s.closed.Load() {
		closeFolder(e.Folder)
		return ErrNoSession
	}
	if err := validateEvent(e, s.state); err != nil {
		closeFolder(e.Folder)
		return err
	}
	s.startRun(e)
	return nil
}
