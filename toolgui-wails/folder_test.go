package tgwails

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"
)

func TestDesktopFolder(t *testing.T) {
	dir := t.TempDir()
	if err := os.Mkdir(filepath.Join(dir, "sub"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "sub", "a.txt"), []byte("hello"), 0644); err != nil {
		t.Fatal(err)
	}
	f, err := openFolder(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if err := fstest.TestFS(f, "sub", "sub/a.txt"); err != nil {
		t.Fatal(err)
	}
	if f.Path() != dir || !f.Writable() {
		t.Fatalf("unexpected grant: %q, writable %v", f.Path(), f.Writable())
	}
	if err := f.MkdirAll("new/nested"); err != nil {
		t.Fatal(err)
	}
	w, err := f.Create("new/nested/b.txt")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := io.WriteString(w, "written"); err != nil {
		t.Fatal(err)
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	bs, err := fs.ReadFile(f, "new/nested/b.txt")
	if err != nil || string(bs) != "written" {
		t.Fatalf("read = %q, %v", bs, err)
	}
	w, err = f.Create("new/nested/b.txt")
	if err != nil {
		t.Fatal(err)
	}
	w.Close()
	info, err := f.Stat("new/nested/b.txt")
	if err != nil || info.Size() != 0 {
		t.Fatalf("truncate: %v, %v", info, err)
	}
	if _, err := f.Create("missing/a"); err == nil {
		t.Fatal("Create must not create parents")
	}
	for _, name := range []string{"../escape", "/absolute", "a/../escape", ""} {
		if w, err := f.Create(name); err == nil {
			w.Close()
			t.Errorf("Create(%q) succeeded", name)
		}
		if err := f.MkdirAll(name); err == nil {
			t.Errorf("MkdirAll(%q) succeeded", name)
		}
	}
	if _, err := os.Stat(filepath.Join(dir, "..", "escape")); !os.IsNotExist(err) {
		t.Fatalf("escaped file: %v", err)
	}
}

func TestDesktopFolderSymlink(t *testing.T) {
	dir, outside := t.TempDir(), t.TempDir()
	if err := os.WriteFile(filepath.Join(outside, "secret"), []byte("secret"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(dir, "out")); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}
	f, err := openFolder(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if file, err := f.Open("out/secret"); err == nil {
		file.Close()
		t.Error("read escaped root")
	}
	if w, err := f.Create("out/new"); err == nil {
		w.Close()
		t.Error("write escaped root")
	}
	if err := f.MkdirAll("out/nested"); err == nil {
		t.Error("mkdir escaped root")
	}
	if _, err := os.Stat(filepath.Join(outside, "new")); !os.IsNotExist(err) {
		t.Fatalf("outside changed: %v", err)
	}
}
