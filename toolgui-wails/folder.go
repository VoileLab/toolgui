package tgwails

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/voilelab/toolgui/toolgui/tgframe"
	"github.com/voilelab/toolgui/toolgui/tgutil"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type desktopFolder struct {
	fs.FS
	root *os.Root
	path string
}

var _ tgframe.Folder = (*desktopFolder)(nil)

func openFolder(path string) (*desktopFolder, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	root, err := os.OpenRoot(path)
	if err != nil {
		return nil, err
	}
	f := root.FS()
	return &desktopFolder{FS: f, root: root, path: path}, nil
}
func (f *desktopFolder) Create(name string) (io.WriteCloser, error) {
	if !fs.ValidPath(name) {
		return nil, &fs.PathError{Op: "create", Path: name, Err: fs.ErrInvalid}
	}
	return f.root.Create(filepath.FromSlash(name))
}
func (f *desktopFolder) MkdirAll(name string) error {
	if !fs.ValidPath(name) {
		return &fs.PathError{Op: "mkdir", Path: name, Err: fs.ErrInvalid}
	}
	return f.root.MkdirAll(filepath.FromSlash(name), 0755)
}
func (f *desktopFolder) Writable() bool { return true }
func (f *desktopFolder) Path() string   { return f.path }
func (f *desktopFolder) Close() error   { return f.root.Close() }

// PickFolder opens a native dialog and reruns with a trusted folder grant.
// Cancellation leaves the state and current run untouched.
func (t *ToolGUI) PickFolder(componentID, defaultDir string) error {
	t.folderLock.Lock()
	defer t.folderLock.Unlock()
	session := t.host.Session()
	if session == nil {
		return ErrNoSession
	}
	if err := session.ValidateFolderPick(componentID); err != nil {
		return err
	}
	t.windowLock.Lock()
	ctx := t.ctx
	t.windowLock.Unlock()
	if ctx == nil {
		return tgutil.NewError("wails runtime is not ready")
	}
	path, err := t.openDirectoryDialog(ctx, runtime.OpenDialogOptions{DefaultDirectory: defaultDir})
	if err != nil || path == "" {
		return err
	}
	folder, err := openFolder(path)
	if err != nil {
		return err
	}
	return t.host.HandleFolder(session, &tgframe.EventFolder{ID: componentID, Folder: folder})
}

func (f *desktopFolder) ReadDir(name string) ([]fs.DirEntry, error) { return fs.ReadDir(f.FS, name) }
func (f *desktopFolder) Stat(name string) (fs.FileInfo, error)      { return fs.Stat(f.FS, name) }
