# Desktop App

The same App can run in a desktop window instead of a browser, through
[Wails v2](https://wails.io). Only the executor changes: pages, components and
state work exactly as they do on the web.

```go
package main

import (
	"log"

	"github.com/voilelab/toolgui/toolgui/tgcomp"
	"github.com/voilelab/toolgui/toolgui/tgframe"

	tgwails "github.com/voilelab/toolgui/toolgui-wails"
)

func main() {
	app := tgframe.NewApp()
	app.AddPage("index", "Index", func(p *tgframe.Params) error {
		tgcomp.Text(p.Main, "Hello world")
		return nil
	})

	e := tgwails.NewExecutor(app, &tgwails.Conf{Title: "ToolGUI Hello"})
	if err := e.Run(); err != nil {
		log.Fatal(err)
	}
}
```

`Run` opens the window and blocks until the user closes it.

## The example app

[`toolgui-wails/example/hello`](https://github.com/voilelab/toolgui/tree/dev/toolgui-wails/example/hello)
is a runnable version of that program, with a sidebar textbox and a button to
show the state round trip working in a window:

```go
func Main(p *tgframe.Params) error {
	name := tgcomp.Textbox(p.Sidebar, "What's your name?")
	if name != "" {
		tgcomp.Text(p.Sidebar, "Hi "+name+"~")
	}

	tgcomp.Text(p.Main, "hello ")
	if tgcomp.Button(p.Main, "keep going") {
		tgcomp.Text(p.Main, "world")
	}

	return nil
}

func main() {
	app := tgframe.NewApp()
	app.AddPage("index", "Index", Main)

	e := tgwails.NewExecutor(app, &tgwails.Conf{Title: "ToolGUI Hello"})
	err := e.Run()
	if err != nil {
		log.Fatal(err)
	}
}
```

## Build your app

In your own Go module, install the desktop package and the build tool:

```shell
go get github.com/voilelab/toolgui/toolgui-wails
go get -tool github.com/voilelab/toolgui/toolgui-wails/cmd/toolgui-desktop
```

The app only needs its Go source. From the directory containing `main.go`:

```shell
go tool toolgui-desktop build   # executable in build/bin/
go tool toolgui-desktop run     # compile, then open the window
```

For an entry point in another directory:

```shell
go tool toolgui-desktop build -o build/mytool ./cmd/mytool
go tool toolgui-desktop run ./cmd/mytool -- -input "data file.csv"
```

Put flags before the package and application arguments after `--`. `run`
uses a temporary executable, keeps the caller's working directory and removes
the executable on exit. It does not watch files or automatically restart.

The frontend is already embedded in tagged releases, so building your app
needs no Node, Yarn, `wails.json` or Wails CLI. The tool supplies `desktop,production`
and the platform flags, including the macOS linker framework and Linux's
WebKitGTK selection. `-tags` adds custom tags; `-ldflags` passes Go linker flags.

Builds target the host OS and architecture. Output is an executable (`.exe`
on Windows); macOS `.app` bundles, installers, icons and signing require a
separate packaging step.

### System dependencies

macOS needs Xcode Command Line Tools and cgo. Windows needs WebView2 to run
the app; Wails uses its default runtime installation behavior. Windows builds
hide the console unless `-console` is set.

Linux needs a C compiler, cgo, `pkg-config`, GTK3 and WebKitGTK development
packages. On Debian and Ubuntu:

```shell
sudo apt-get install pkg-config libgtk-3-dev libwebkit2gtk-4.1-dev
```

The tool prefers WebKitGTK 4.1, then falls back to 4.0. Use `-webkit 4.0` or
`-webkit 4.1` to select a version explicitly. Missing dependencies produce an
error before compilation.

## Running the repository example

From the ToolGUI repository root:

```shell
task run_desktop_hello
task build_desktop_hello   # executable in toolgui-wails/build/bin/
```

These tasks build the frontend before invoking the desktop CLI. Assets are
not committed on `dev`; if invoking the CLI directly from a checkout, run
`task asset_wails` first. The CLI rejects missing assets and the placeholders
created by `task stub_assets`.

For frontend development, `task run_wails_hello` still uses Wails' file
watcher and serves assets from disk. `task build_wails_hello` retains Wails'
packaging flow. Those tasks use the example's `wails.json` and require Node
and Yarn; details are in the
[desktop README](https://github.com/voilelab/toolgui/blob/dev/toolgui-wails/README.md#working-on-toolgui-itself).

## Conf

```go
type Conf struct {
	// Title shows in the title bar.
	Title string

	// Width and Height are the window size in pixels.
	Width, Height int

	// WindowStartState sets the initial window state. The zero value is Normal.
	WindowStartState WindowStartState

	// FullscreenToggle adds View > Toggle Fullscreen to the native menu.
	// The shortcut is F11, or Ctrl+Cmd+F on macOS. Disabled by default.
	FullscreenToggle bool

	// Background is the window background colour.
	Background *RGBA

	// DisableResize stops the user resizing the window.
	DisableResize bool

	// Frameless drops the window decorations.
	Frameless bool
}
```

Empty fields fall back to `DefaultConf`: a 1024x768 white window titled
`ToolGUI`. An [app title](../app/index.md#title) names the window ahead of that
default, so a `Conf` only sets `Title` to override it.

Set `WindowStartState` to `tgwails.Maximised` to open maximised, or
`tgwails.Fullscreen` to open fullscreen. `tgwails.Normal` is the zero value,
so omitting the field keeps the normal window size from `Width` and `Height`.
`tgwails.Minimised` is also exported. The type aliases Wails' window start state;
apps do not need to import Wails' `options` package.

With `DisableResize: true`, Wails v2.15.0 ignores `Maximised` on Windows and
uses the normal size; `Fullscreen` still takes the fullscreen path. macOS
passes both requests to AppKit. See the
[desktop README](https://github.com/voilelab/toolgui/blob/dev/toolgui-wails/README.md#conf)
for the platform details and outstanding manual GUI checks.

Set `FullscreenToggle: true` to add `View > Toggle Fullscreen`, using F11 on
Windows/Linux and Ctrl+Cmd+F on macOS. It works without an app menu and does
not rerun the page or change session state. An existing `View` submenu is
reused. App shortcuts take precedence: if the fullscreen shortcut is already
assigned, the built-in command remains clickable without an accelerator.
The default is false. On macOS,
enabling it replaces the standard Window role (including its Minimize and
Zoom menu entries) to avoid a duplicate fullscreen command. See the
[menu notes](https://github.com/voilelab/toolgui/blob/dev/toolgui-wails/README.md#fullscreen-toggle).

## A separate module

`github.com/voilelab/toolgui/toolgui-wails` is its own Go module, because Wails
needs cgo, GTK and WebKit. Apps that only target the web never pull any of that
in.
