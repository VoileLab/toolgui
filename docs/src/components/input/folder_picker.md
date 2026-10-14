# FolderPicker

`FolderPicker` draws a button that opens the desktop's native directory dialog.
It returns `nil` before a selection. Selecting a directory reruns the page, and
that run receives the new folder. Cancelling keeps the previous selection and
does not rerun.

## API

```go
type FolderObject = tgframe.Folder

type FolderPickerConf struct {
    tgframe.Base // ID
    Write bool
    DefaultDir string
    Disabled bool
}

func FolderPicker(c *tgframe.Container, label string, conf ...*FolderPickerConf) FolderObject
```

`conf` is optional, with at most one configuration. `Write` requests write
access; desktop grants it unconditionally. `DefaultDir` sets the desktop
dialog's initial directory. `Disabled` disables the button.

The shared interface lives in `tgframe`:

```go
type Folder interface {
    fs.ReadDirFS // Open and ReadDir
    fs.StatFS
    Create(name string) (io.WriteCloser, error)
    MkdirAll(name string) error
    Writable() bool
    Path() string
}

var ErrFolderReadOnly = tgutil.NewError("folder is read-only")
```

Names follow `fs.ValidPath`: slash-separated paths relative to the selected
root, with no `..` components. Desktop uses `os.OpenRoot`, which also prevents
symlinks from escaping the root. `Create` creates or truncates a file and
requires an existing parent. `MkdirAll` creates missing parents. Read-only
implementations return `tgframe.ErrFolderReadOnly` for either write operation.
Close files and writers returned by `Open` and `Create` when finished.

`Writable` reports granted write access. `Path` returns an absolute OS path on
desktop; other executors return an empty string. The session owns the folder
and closes it when replaced or the session ends; use it within the current
page run rather than keeping it across selections or sessions.

Web and Wasm currently show a disabled button with an explanation and return
`nil`. Their directory access implementations are planned separately.

## Example

```go
folder := tgcomp.FolderPicker(p.Main, "Output folder", &tgcomp.FolderPickerConf{
    Write: true,
})
if folder != nil && tgcomp.Button(p.Main, "Export") {
    if err := folder.MkdirAll("results"); err != nil {
        return err
    }
    file, err := folder.Create("results/report.txt")
    if err != nil {
        return err
    }
    defer file.Close()
    if _, err := io.WriteString(file, "Export complete\n"); err != nil {
        return err
    }
    return file.Close()
}
```
