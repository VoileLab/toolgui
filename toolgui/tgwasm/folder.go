//go:build js && wasm

package tgwasm

import (
	"io"
	"io/fs"
	"path"
	"sort"
	"strings"
	"sync"
	"syscall/js"
	"time"

	"github.com/voilelab/toolgui/toolgui/internal/opfs"
	"github.com/voilelab/toolgui/toolgui/tgframe"
)

// User-picked files are outside OPFS: createSyncAccessHandle is unavailable.
// Await the standard handle APIs on the page goroutine, never a JS callback.
// Reads use slices of a File snapshot; writes stream chunks to a swap file,
// committed by Close. Neither operation loads an entire file into the Go heap.
type browserFolder struct {
	root  js.Value
	write bool
}

var _ tgframe.Folder = (*browserFolder)(nil)

func (f *browserFolder) Writable() bool { return f.write }
func (f *browserFolder) Path() string   { return "" }

func folderError(op, name string, err error) error {
	if err == nil {
		return nil
	}
	switch {
	case strings.Contains(err.Error(), "NotFoundError:"):
		err = fs.ErrNotExist
	case strings.Contains(err.Error(), "NotAllowedError:"), strings.Contains(err.Error(), "SecurityError:"):
		err = fs.ErrPermission
	case strings.Contains(err.Error(), "TypeMismatchError:"):
		err = fs.ErrInvalid
	}
	return &fs.PathError{Op: op, Path: name, Err: err}
}

func (f *browserFolder) directory(name string, create bool) (js.Value, error) {
	if !fs.ValidPath(name) {
		return js.Undefined(), fs.ErrInvalid
	}
	dir := f.root
	if name == "." {
		return dir, nil
	}
	for part := range strings.SplitSeq(name, "/") {
		var err error
		dir, err = opfs.AwaitCall(dir, "getDirectoryHandle", part, map[string]any{"create": create})
		if err != nil {
			return js.Undefined(), err
		}
	}
	return dir, nil
}

func (f *browserFolder) handle(name string) (js.Value, error) {
	if !fs.ValidPath(name) {
		return js.Undefined(), fs.ErrInvalid
	}
	if name == "." {
		return f.root, nil
	}
	parent, err := f.directory(path.Dir(name), false)
	if err != nil {
		return js.Undefined(), err
	}
	handle, err := opfs.AwaitCall(parent, "getFileHandle", path.Base(name))
	// Only a type mismatch means it is a directory. Do not hide missing or
	// permission failures behind a second lookup.
	if err != nil && strings.Contains(err.Error(), "TypeMismatchError:") {
		return opfs.AwaitCall(parent, "getDirectoryHandle", path.Base(name))
	}
	return handle, err
}

type folderInfo struct {
	name string
	size int64
	mod  time.Time
	dir  bool
}

func (i folderInfo) Name() string { return i.name }
func (i folderInfo) Size() int64  { return i.size }
func (i folderInfo) Mode() fs.FileMode {
	if i.dir {
		return fs.ModeDir | 0755
	}
	return 0644
}
func (i folderInfo) ModTime() time.Time { return i.mod }
func (i folderInfo) IsDir() bool        { return i.dir }
func (i folderInfo) Sys() any           { return nil }

func handleInfo(name string, handle js.Value) (folderInfo, js.Value, error) {
	info := folderInfo{name: path.Base(name), dir: handle.Get("kind").String() == "directory"}
	if info.dir {
		return info, js.Undefined(), nil
	}
	file, err := opfs.AwaitCall(handle, "getFile")
	if err != nil {
		return info, file, err
	}
	info.size = int64(file.Get("size").Float())
	info.mod = time.UnixMilli(int64(file.Get("lastModified").Float()))
	return info, file, nil
}

func (f *browserFolder) Stat(name string) (fs.FileInfo, error) {
	h, err := f.handle(name)
	if err != nil {
		return nil, folderError("stat", name, err)
	}
	info, _, err := handleInfo(name, h)
	if err != nil {
		return nil, folderError("stat", name, err)
	}
	return info, nil
}

func (f *browserFolder) ReadDir(name string) ([]fs.DirEntry, error) {
	dir, err := f.directory(name, false)
	if err != nil {
		return nil, folderError("readdir", name, err)
	}
	iter, err := opfs.Call(dir, "entries")
	if err != nil {
		return nil, folderError("readdir", name, err)
	}
	entries := make([]fs.DirEntry, 0)
	for {
		item, err := opfs.AwaitCall(iter, "next")
		if err != nil {
			return nil, folderError("readdir", name, err)
		}
		if item.Get("done").Bool() {
			break
		}
		pair := item.Get("value")
		info, _, err := handleInfo(pair.Index(0).String(), pair.Index(1))
		if err != nil {
			return nil, folderError("readdir", name, err)
		}
		entries = append(entries, fs.FileInfoToDirEntry(info))
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].Name() < entries[j].Name() })
	return entries, nil
}

func (f *browserFolder) Open(name string) (fs.File, error) {
	h, err := f.handle(name)
	if err != nil {
		return nil, folderError("open", name, err)
	}
	info, file, err := handleInfo(name, h)
	if err != nil {
		return nil, folderError("open", name, err)
	}
	return &folderFile{folder: f, name: name, info: info, blob: file}, nil
}

type folderFile struct {
	mu             sync.Mutex
	folder         *browserFolder
	name           string
	info           folderInfo
	blob           js.Value
	buffer         []byte
	offset         int64
	entries        []fs.DirEntry
	listed, closed bool
}

func (f *folderFile) Stat() (fs.FileInfo, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.closed {
		return nil, folderError("stat", f.name, fs.ErrClosed)
	}
	return f.info, nil
}
func (f *folderFile) Close() error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.closed = true
	f.blob = js.Undefined()
	f.buffer = nil
	f.entries = nil
	return nil
}
func (f *folderFile) Read(p []byte) (int, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.closed {
		return 0, folderError("read", f.name, fs.ErrClosed)
	}
	if f.info.dir {
		return 0, folderError("read", f.name, fs.ErrInvalid)
	}
	if len(p) == 0 {
		return 0, nil
	}
	if f.offset >= f.info.size {
		return 0, io.EOF
	}
	if len(f.buffer) == 0 {
		// Read ahead one bounded chunk: small Go reads should not require a
		// disk operation and a Promise for every byte.
		size := min(f.info.size-f.offset, int64(opfs.ChunkSize))
		slice, err := opfs.Call(f.blob, "slice", float64(f.offset), float64(f.offset+size))
		if err != nil {
			return 0, folderError("read", f.name, err)
		}
		buffer, err := opfs.AwaitCall(slice, "arrayBuffer")
		if err != nil {
			return 0, folderError("read", f.name, err)
		}
		f.buffer = make([]byte, size)
		n := js.CopyBytesToGo(f.buffer, opfs.Uint8Array.New(buffer))
		if n == 0 {
			return 0, folderError("read", f.name, io.ErrNoProgress)
		}
		f.buffer = f.buffer[:n]
	}
	n := copy(p, f.buffer)
	f.buffer = f.buffer[n:]
	f.offset += int64(n)
	return n, nil
}
func (f *folderFile) ReadDir(n int) ([]fs.DirEntry, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.closed {
		return nil, folderError("readdir", f.name, fs.ErrClosed)
	}
	if !f.info.dir {
		return nil, folderError("readdir", f.name, fs.ErrInvalid)
	}
	if !f.listed {
		var err error
		f.entries, err = f.folder.ReadDir(f.name)
		if err != nil {
			return nil, err
		}
		f.listed = true
	}
	if n > 0 && len(f.entries) == 0 {
		return nil, io.EOF
	}
	count := len(f.entries)
	if n > 0 {
		count = min(n, count)
	}
	entries := f.entries[:count]
	f.entries = f.entries[count:]
	return entries, nil
}

func (f *browserFolder) MkdirAll(name string) error {
	if !f.write {
		return folderError("mkdir", name, tgframe.ErrFolderReadOnly)
	}
	_, err := f.directory(name, true)
	return folderError("mkdir", name, err)
}

func (f *browserFolder) Create(name string) (io.WriteCloser, error) {
	if !f.write {
		return nil, folderError("create", name, tgframe.ErrFolderReadOnly)
	}
	if !fs.ValidPath(name) || name == "." {
		return nil, folderError("create", name, fs.ErrInvalid)
	}
	dir, err := f.directory(path.Dir(name), false)
	if err != nil {
		return nil, folderError("create", name, err)
	}
	h, err := opfs.AwaitCall(dir, "getFileHandle", path.Base(name), opfs.Create)
	if err != nil {
		return nil, folderError("create", name, err)
	}
	stream, err := opfs.AwaitCall(h, "createWritable", map[string]any{"keepExistingData": false})
	if err != nil {
		return nil, folderError("create", name, err)
	}
	return &folderWriter{name: name, stream: stream}, nil
}

type folderWriter struct {
	mu     sync.Mutex
	name   string
	stream js.Value
	closed bool
	err    error
}

func (w *folderWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.closed {
		return 0, folderError("write", w.name, fs.ErrClosed)
	}
	if w.err != nil {
		return 0, w.err
	}
	written := 0
	for len(p) > 0 {
		n := min(len(p), opfs.ChunkSize)
		chunk := opfs.Uint8Array.New(n)
		js.CopyBytesToJS(chunk, p[:n])
		if _, err := opfs.AwaitCall(w.stream, "write", chunk); err != nil {
			w.err = folderError("write", w.name, err)
			// A failed stream must never commit the successfully written prefix.
			_, _ = opfs.AwaitCall(w.stream, "abort")
			return written, w.err
		}
		written += n
		p = p[n:]
	}
	return written, nil
}
func (w *folderWriter) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.closed {
		return w.err
	}
	w.closed = true
	if w.err == nil {
		_, err := opfs.AwaitCall(w.stream, "close")
		w.err = folderError("close", w.name, err)
	}
	w.stream = js.Undefined()
	return w.err
}
