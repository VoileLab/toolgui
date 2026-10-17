# toolgui-wails

Runs a ToolGUI app in a desktop window with [Wails v2](https://wails.io),
instead of serving it over HTTP.

```go
app := tgframe.NewApp()
app.AddPage("index", "Index", Main)

e := tgwails.NewExecutor(app, &tgwails.Conf{Title: "My Tool"})
e.Run()
```

The pages, components and state are the same as the web executor's: only the
transport changes. `tgframe.Session` does the work either way, and the frontend
is the same `@toolgui-web/lib` components with a different adapter around them.

## Conf

`Conf.WindowStartState` selects `tgwails.Normal` (the zero value),
`tgwails.Maximised`, `tgwails.Minimised` or `tgwails.Fullscreen`. The type is an
alias of Wails' `options.WindowStartState`; apps need no Wails import:

```go
e := tgwails.NewExecutor(app, &tgwails.Conf{
	Title:            "My Tool",
	WindowStartState: tgwails.Maximised,
})
```

A nil or zero `Conf`, or one that omits this field, starts normally.
`Width` and `Height` still set the normal window size, defaulting to 1024x768.

`DisableResize: true` is passed to Wails together with the start state. In
Wails v2.15.0, Windows ignores `Maximised` for a non-resizable window and opens
it at its normal size; `Fullscreen` still takes the fullscreen path. On macOS,
Wails sends the zoom/fullscreen request to AppKit even when resizing is disabled.
These are upstream code paths, not results of a manual GUI check. Before release,
manually check `Maximised` and `Fullscreen` on both macOS and Windows, with
`DisableResize` both false and true: confirm the initial state, that startup
does not panic, and the platform's non-resizable-window behaviour. Also check
that omitting the state still uses the configured `Width` and `Height`.

## Why a separate module

Wails needs cgo, GTK and WebKit. Keeping it in its own module means people who
only build web apps never pull any of that in.

In-repo development builds against the working tree:

```
replace github.com/voilelab/toolgui => ../
```

A replace only applies in the main module, so it does not follow the module to
anyone depending on it. What they resolve is the `require` above it, which
`release.yml` stamps with the version being tagged.

## Installing it

```shell
go get github.com/voilelab/toolgui/toolgui-wails
```

Being a module in a subdirectory, it is tagged `toolgui-wails/vX.Y.Z` rather
than `vX.Y.Z`; `go get` finds that on its own. Desktop builds still need the
webview dependencies below.

## What crosses the boundary

| Web | Desktop |
| --- | --- |
| `GET /api/app` | `window.go.tgwails.ToolGUI.AppConf()` |
| update websocket | `Update(eventJSON)` + the `toolgui:pack` event |
| Native folder selection (desktop only) | `PickFolder(componentID, defaultDir)`; trusted Go grant, then rerun; cancel keeps state |
| `POST /api/files` | `UploadFileStart` + `UploadFileChunk` + `UploadFileFinish` |
| a page load | `Start(pageName, query)` |
| the menubar the frontend draws | the window's own, from `options.App.Menu` |

Payloads cross as JSON strings — the same ones the websocket carries, so both
transports share a wire format. An upload is the exception: bound methods take
strings, so a file crosses in chunks rather than as one base64 blob. The state
writes them to disk as they land, and the component sees the file only once
the last chunk is in. Packs go out on a single event name so
create/update/delete/result keep the order the page produced them in.

## Links

A desktop window has no tabs and no back button, so letting an external link
navigate would replace the running app, page state and all. The adapter
catches clicks on the document and hands `http(s)` urls whose origin is not
the app's own to `runtime.BrowserOpenURL`, which opens them in the system
browser. Page navigation and other schemes are left alone, and `Link` and
`LinkButton` render the same plain anchors they do on the web — nothing to set
on them.

Wails serves the frontend from its own origin, so the asset server could carry
the plain HTTP endpoints too. Bound methods handle all four instead, to keep
this module off `tgexec`, which embeds the whole web bundle.

## Menus

### Fullscreen toggle

Set `Conf.FullscreenToggle` to opt in to a native `View > Toggle Fullscreen`
command, even when the app declares no menu:

```go
e := tgwails.NewExecutor(app, &tgwails.Conf{FullscreenToggle: true})
```

Click the command or press **F11** on Windows/Linux, or **Ctrl+Cmd+F** on
macOS, to enter or leave fullscreen. This is a window action handled directly
in Go: it does not send a session event, rerun the page or reset page state.
Clicks before Wails startup, after shutdown or with a cancelled runtime
context are ignored. The switch defaults to false, preserving the existing
menus and shortcuts.

An existing top-level `View` submenu is reused. If an app command already
uses the platform's fullscreen shortcut, the app keeps it and the built-in
fullscreen command has no accelerator; use the menu to toggle fullscreen.
Shortcuts in nested submenus count too.

When enabled on macOS, the App and Edit roles remain, but the standard Window
role is omitted because it includes a competing fullscreen command with the
same shortcut. Its Minimize and Zoom menu entries are therefore absent too;
the new View command is available for both decorated and frameless windows.
When disabled, the usual Window role is kept for decorated windows.

For release validation on each target OS, enter and exit using both the
menu and shortcut, with and without an app menu;
check that macOS shows only one fullscreen command and that page state stays
unchanged. Repeat with `Frameless: true` and with the toggle disabled.

### App menus

`App.SetMenu` declares one tree, and the desktop window draws it itself:

```go
app.SetMenu(tgframe.NewMenu().
	Submenu("File", func(m *tgframe.Menu) {
		m.Text("Open", "file_open")
		m.Separator()
		m.Text("Quit", "file_quit")
	}))
```

The same declaration is a row of buttons above the app in a browser. Here it
is translated into a `menu.Menu` and handed to `options.App.Menu`, so what the
window shows is the platform's real menubar. The tree is then kept out of the
app config the frontend receives — otherwise it would draw its own menubar
inside the window, under the real one.

Nothing new carries the click back. A Wails `MenuItem.Click` callback happens
in Go, so it goes into the running session as the click event the web menubar
would have sent over the wire:

```go
session.HandleEvent(&tgframe.EventClick{ID: "menu_item_file_open"})
```

which is the ordinary rerun, with the packs going out to the frontend as
usual. The page reads the pick with `tgframe.MenuClicked(p, "file_open")`
whichever executor it is running under. The wire format does not change and
there is no new event type.

The picks go through a queue that one goroutine works through, rather than a
goroutine each. Windows calls the callback on the message loop, where applying
the pick there and then freezes the window for as long as the page takes to
run; but a run cuts the one before it, so picks that overtake each other would
leave the page showing the older one. Queueing costs the callback a mutex held
for the length of an append and keeps both. It can only order what reaches it
in order, which on Windows is everything — macOS and Linux hand each callback
its own goroutine before it gets there.

A pick made before `Start` has opened a session — the window is up, the
menubar with it, and no page has loaded yet — has nothing to run and is
ignored.

### What the platform adds

macOS actions like Cmd+C, Cmd+V and Cmd+Q belong to the standard App and Edit
menus rather than to the webview, so a window that replaces the menubar
without them loses those shortcuts. The App and Edit menus therefore go in
front of the app's own entries there, plus the Window menu unless the window
is frameless — the same set, under the same rule, that Wails gives a window
with no menu at all. So an app keeps the standard menus whether or not it
declares one of its own.

This is the one place the menubar is not the same on every platform, and it
cannot be helped: Wails v2 serves menu roles on macOS only. The translation
does not expose roles at all — a `tgframe.Menu` is text items, separators and
submenus, and nothing else.

### What it does not do

- **A menu item cannot change page.** On the desktop the page switch happens
  in the frontend, and Go has no way to ask it for one: `SendPackFunc` carries
  create, update, delete, ready and result, none of which navigate. That would
  need a new Go-to-frontend pack type.
- **Nothing is enabled or disabled while running.** The tree is declared once
  on the app, which is the shape a native menu — with no diff to apply — can
  take. `runtime.MenuUpdateApplicationMenu` rebuilds the whole menu, closing
  any dropdown that happens to be open, so a menu that changed on every run
  would be a menu that cannot be used.

## Building

Install the desktop package and build tool in your app's Go module:

```shell
go get github.com/voilelab/toolgui/toolgui-wails
go get -tool github.com/voilelab/toolgui/toolgui-wails/cmd/toolgui-desktop
```

Then build or run the package that calls `tgwails.NewExecutor(...).Run()`:

```shell
go tool toolgui-desktop build                  # ./build/bin/<package>
go tool toolgui-desktop build -o build/mytool ./cmd/mytool
go tool toolgui-desktop run ./cmd/mytool
go tool toolgui-desktop run ./cmd/mytool -- -input "data file.csv"
```

Flags precede the package; the package defaults to `.`. Output is a native
executable, with `.exe` on Windows. `-o` names a file, relative to the current
working directory. `run` builds in a temporary directory, launches the app in
the caller's working directory and removes the executable when it exits. App
arguments go after `--`, and the CLI preserves the app's exit code.

Tagged releases already contain the frontend. These commands need no Node,
Yarn, Wails CLI, `wails.json` or generated JS bindings. In a ToolGUI source
checkout, build real assets with `task asset_wails` first: the CLI rejects
missing assets and the placeholders made by `task stub_assets`. It checks the
desktop package your app actually imports, including local `replace` paths.

The tool builds for the host OS and architecture and rejects cross compilation.
It produces executables; `.app` bundles, installers, icons, signing and
notarization need a separate packaging step. `run` does one build and launch,
without file watching.

### Platform requirements and flags

- macOS needs Xcode Command Line Tools and cgo. The CLI adds the
  `UniformTypeIdentifiers` linker framework and preserves existing CGO flags.
- Linux needs a C compiler, cgo, `pkg-config`, GTK3 and WebKitGTK development
  packages. The CLI prefers WebKitGTK 4.1 and falls back to 4.0. Select one
  explicitly with `-webkit 4.0` or `-webkit 4.1`.
- Windows uses WebView2. The executable hides the console by default;
  `-console` keeps it. Wails retains its default WebView2 runtime installation
  behavior. The CLI itself does not require a C compiler on Windows.

`production` is automatic. `-tags custom,another` adds tags to the required
ones and those in `GOFLAGS`. `dev` and `bindings` conflict with this build
mode. On Linux, `-tags webkit2_41` also selects WebKitGTK 4.1.

`-ldflags` passes a value to the Go linker, overriding any `GOFLAGS` linker
flags. For example:

```shell
go tool toolgui-desktop build -ldflags="-X main.version=1.0" ./cmd/mytool
```

### Working on ToolGUI itself

From the repository root:

```shell
task run_desktop_hello     # build frontend, then compile and launch
task build_desktop_hello   # executable in toolgui-wails/build/bin
```

Frontend development still uses the Wails CLI, pinned as a tool dependency
of this module:

```shell
task run_wails_hello       # frontend from disk, Go files watched
task build_wails_hello     # Wails packaging, example/hello/build/bin
```

These tasks use `example/hello/wails.json`, which points at the frontend
workspace. They stub embedded assets before Wails generates bindings, then
build the frontend with Yarn. The Wails tasks add `webkit2_41` on Linux.

### Without the CLI

A plain Go build needs the platform flags supplied by the CLI, for example:

```shell
# macOS 11+
CGO_LDFLAGS="-framework UniformTypeIdentifiers" go build -tags production ./example/hello

# Linux with WebKitGTK 4.1
go build -tags production,webkit2_41 ./example/hello
```

### Webview dependencies

On Debian and Ubuntu:

```shell
sudo apt-get install libgtk-3-dev libwebkit2gtk-4.1-dev
```

`.github/workflows/wails.yml` does the same and builds with both tags.
