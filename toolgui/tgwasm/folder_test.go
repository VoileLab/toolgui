//go:build js && wasm

package tgwasm

import (
	"bytes"
	"crypto/rand"
	"errors"
	"io"
	"io/fs"
	"strings"
	"syscall/js"
	"testing"
	"testing/fstest"

	"github.com/voilelab/toolgui/toolgui/internal/opfs"
	"github.com/voilelab/toolgui/toolgui/tgframe"
)

// Use real browser handles, exercising the standard asynchronous APIs also
// available on picker grants. No sync access handle is opened by these tests.
func testFolder(t *testing.T) *browserFolder {
	t.Helper()
	origin, err := opfs.Origin()
	if err != nil {
		t.Fatal(err)
	}
	name := "folder-test-" + strings.ToLower(rand.Text())
	root, err := opfs.AwaitCall(origin, "getDirectoryHandle", name, opfs.Create)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if _, err := opfs.AwaitCall(origin, "removeEntry", name, opfs.Recursive); err != nil {
			t.Error(err)
		}
	})
	return &browserFolder{root: root, write: true}
}

func putFolderFile(t *testing.T, f *browserFolder, name string, body []byte) {
	t.Helper()
	w, err := f.Create(name)
	if err != nil {
		t.Fatal(err)
	}
	n, err := w.Write(body)
	if err != nil || n != len(body) {
		t.Fatalf("Write = %d, %v", n, err)
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	if _, err := w.Write([]byte("late")); !errors.Is(err, fs.ErrClosed) {
		t.Fatalf("closed Write: %v", err)
	}
}

func TestFolderFS(t *testing.T) {
	f := testFolder(t)
	if !f.Writable() || f.Path() != "" {
		t.Fatal("grant metadata")
	}
	if err := f.MkdirAll("nested/deep"); err != nil {
		t.Fatal(err)
	}
	if err := f.MkdirAll("nested/deep"); err != nil {
		t.Fatal(err)
	}
	if err := f.MkdirAll("."); err != nil {
		t.Fatal(err)
	}
	// Several chunks, including a trailing partial chunk and non-ASCII names.
	large := bytes.Repeat([]byte("abc123\n"), opfs.ChunkSize/2)
	putFolderFile(t, f, "z.txt", []byte("last"))
	putFolderFile(t, f, "nested/deep/資料.txt", large)
	putFolderFile(t, f, "empty", nil)
	if err := fstest.TestFS(f, "z.txt", "nested/deep/資料.txt", "empty"); err != nil {
		t.Fatal(err)
	}
	got, err := fs.ReadFile(f, "nested/deep/資料.txt")
	if err != nil || !bytes.Equal(got, large) {
		t.Fatalf("ReadFile: len=%d, err=%v", len(got), err)
	}
	entries, err := f.ReadDir(".")
	if err != nil || len(entries) != 3 || entries[0].Name() != "empty" || entries[1].Name() != "nested" || !entries[1].IsDir() {
		t.Fatalf("ReadDir: %v, %v", entries, err)
	}
	// A shorter replacement must not leave the old tail behind.
	putFolderFile(t, f, "nested/deep/資料.txt", []byte("new"))
	got, err = fs.ReadFile(f, "nested/deep/資料.txt")
	if err != nil || string(got) != "new" {
		t.Fatalf("truncate: %q, %v", got, err)
	}
	info, err := f.Stat("nested/deep/資料.txt")
	if err != nil || info.Size() != 3 || info.IsDir() || info.ModTime().IsZero() {
		t.Fatalf("Stat: %v, %v", info, err)
	}
	fp, err := f.Open("empty")
	if err != nil {
		t.Fatal(err)
	}
	if err := fp.Close(); err != nil {
		t.Fatal(err)
	}
	if _, err := fp.Read(make([]byte, 1)); !errors.Is(err, fs.ErrClosed) {
		t.Fatalf("closed Read: %v", err)
	}
	if _, err := fp.Stat(); !errors.Is(err, fs.ErrClosed) {
		t.Fatalf("closed Stat: %v", err)
	}
}

func TestFolderErrorsAndReadOnly(t *testing.T) {
	f := testFolder(t)
	putFolderFile(t, f, "file", []byte("hello"))
	for _, name := range []string{"", "/absolute", "../escape", "a/../b", "a//b", "a/"} {
		if _, err := f.Open(name); !errors.Is(err, fs.ErrInvalid) {
			t.Errorf("Open(%q): %v", name, err)
		}
		if _, err := f.Stat(name); !errors.Is(err, fs.ErrInvalid) {
			t.Errorf("Stat(%q): %v", name, err)
		}
		if _, err := f.ReadDir(name); !errors.Is(err, fs.ErrInvalid) {
			t.Errorf("ReadDir(%q): %v", name, err)
		}
		if _, err := f.Create(name); !errors.Is(err, fs.ErrInvalid) {
			t.Errorf("Create(%q): %v", name, err)
		}
		if err := f.MkdirAll(name); !errors.Is(err, fs.ErrInvalid) {
			t.Errorf("MkdirAll(%q): %v", name, err)
		}
	}
	for _, name := range []string{"missing", "missing/file"} {
		if _, err := f.Open(name); !errors.Is(err, fs.ErrNotExist) {
			t.Errorf("Open(%q): %v", name, err)
		}
		if _, err := f.Stat(name); !errors.Is(err, fs.ErrNotExist) {
			t.Errorf("Stat(%q): %v", name, err)
		}
		if _, err := f.ReadDir(name); !errors.Is(err, fs.ErrNotExist) {
			t.Errorf("ReadDir(%q): %v", name, err)
		}
	}
	if _, err := f.Create("missing/file"); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("missing parent: %v", err)
	}
	if _, err := f.Create("."); !errors.Is(err, fs.ErrInvalid) {
		t.Fatalf("create root: %v", err)
	}
	if _, err := f.ReadDir("file"); !errors.Is(err, fs.ErrInvalid) {
		t.Fatalf("readdir file: %v", err)
	}
	ro := &browserFolder{root: f.root, write: false}
	if ro.Writable() {
		t.Fatal("read-only metadata")
	}
	if _, err := ro.Create("new"); !errors.Is(err, tgframe.ErrFolderReadOnly) {
		t.Fatalf("readonly Create: %v", err)
	}
	if err := ro.MkdirAll("new"); !errors.Is(err, tgframe.ErrFolderReadOnly) {
		t.Fatalf("readonly MkdirAll: %v", err)
	}
	if err := fstest.TestFS(ro, "file"); err != nil {
		t.Fatal(err)
	}
	// Real error names must stay recognizable to fs callers.
	for _, name := range []string{"NotAllowedError", "SecurityError"} {
		err := folderError("open", "x", opfs.Err(js.Global().Get("DOMException").New("denied", name)))
		if !errors.Is(err, fs.ErrPermission) {
			t.Fatal(err)
		}
	}
}

func TestFolderBridgeGrantsAndNavigation(t *testing.T) {
	f := testFolder(t)
	b := newBridge(testApp())
	discardPacks(b)
	b.jsStart(js.Undefined(), []js.Value{js.ValueOf("index")})
	settle(b)
	begin := func(write bool) string {
		t.Helper()
		value, err := opfs.Await(b.jsBeginFolderPick(js.Undefined(), []js.Value{js.ValueOf("comp"), js.ValueOf(write)}).(js.Value))
		if err != nil {
			t.Fatal(err)
		}
		return value.String()
	}
	slot := begin(false)
	if hostState(b).GetFolder("comp") != nil {
		t.Fatal("reservation must not change state")
	}
	b.jsCancelFolderPick(js.Undefined(), []js.Value{js.ValueOf(slot)})
	if hostState(b).GetFolder("comp") != nil {
		t.Fatal("cancelled picker changed state")
	}
	slot = begin(false)
	_, err := opfs.Await(b.jsFolderPicked(js.Undefined(), []js.Value{js.ValueOf(slot), f.root}).(js.Value))
	if err != nil {
		t.Fatal(err)
	}
	settle(b)
	grant := hostState(b).GetFolder("comp")
	if grant == nil || grant.Writable() || grant.Path() != "" {
		t.Fatal("read grant not installed")
	}
	slot = begin(true)
	_, err = opfs.Await(b.jsFolderPicked(js.Undefined(), []js.Value{js.ValueOf(slot), f.root}).(js.Value))
	if err != nil {
		t.Fatal(err)
	}
	settle(b)
	grant = hostState(b).GetFolder("comp")
	if grant == nil || !grant.Writable() {
		t.Fatal("write grant not installed")
	}
	w, err := grant.Create("bridge.txt")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := io.WriteString(w, "bridge"); err != nil {
		t.Fatal(err)
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	slot = begin(true)
	b.jsStart(js.Undefined(), []js.Value{js.ValueOf("index")})
	settle(b)
	_, err = opfs.Await(b.jsFolderPicked(js.Undefined(), []js.Value{js.ValueOf(slot), f.root}).(js.Value))
	if err == nil || hostState(b).GetFolder("comp") != nil {
		t.Fatal("stale dialog must not grant the replacement session")
	}
	_, err = opfs.Await(b.jsBeginFolderPick(js.Undefined(), []js.Value{js.ValueOf("missing"), js.ValueOf(true)}).(js.Value))
	if err == nil {
		t.Fatal("invalid target reserved")
	}
}
