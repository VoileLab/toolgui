# ToolGUI

[![Go Reference](https://pkg.go.dev/badge/github.com/voilelab/toolgui.svg)](https://pkg.go.dev/github.com/voilelab/toolgui)

ToolGUI is a Go framework for quickly building interactive dashboards and
internal tools — a Streamlit-like experience for Go. Write a plain Go
function; ToolGUI renders it as a web page and reruns it on each interaction.

The same app runs three ways:

- **Web server** — a Go backend serving the page.
- **Desktop** — a native window via Wails ([`toolgui-wails`](toolgui-wails/README.md)).
- **Browser only** — compiled to WebAssembly, no server ([`toolgui/tgwasm`](toolgui/tgwasm/README.md)).

> [!WARNING]
> The API is still under development and may change.

- **Live demo:** <https://voilelab.github.io/toolgui/demo/> — every component,
  running in your browser as WebAssembly.
- **Docs:** <https://voilelab.github.io/toolgui/>, built from [`docs/`](docs/).
  They track the latest release; `docs/` on `dev` may be ahead.

## Install

```shell
go get github.com/voilelab/toolgui
```

## Hello world

```go
package main

import (
	"log"

	"github.com/voilelab/toolgui/toolgui/tgcomp"
	"github.com/voilelab/toolgui/toolgui/tgexec"
	"github.com/voilelab/toolgui/toolgui/tgframe"
)

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
	e := tgexec.NewWebExecutor(app)
	log.Println("Starting service...")
	e.StartService("127.0.0.1:3000")
}
```

Open <http://127.0.0.1:3000>.

> [!IMPORTANT]
> ToolGUI has no built-in auth: anyone who reaches the page can run the tool.
> That is why the example binds to `127.0.0.1`. To expose it further, put it
> behind an authenticating proxy — see
> [Serving It Safely](https://voilelab.github.io/toolgui/architecture/exposure.html).

## Development

Branch from `dev` and open pull requests against `dev`; see
[CONTRIBUTING.md](CONTRIBUTING.md).

### Prerequisites

- [Go](https://go.dev/) — version in [`go.mod`](go.mod)
- [Task](https://taskfile.dev/) — task runner
- [Yarn](https://yarnpkg.com/) — frontend
- [Cypress](https://www.cypress.io/) — E2E tests
- [mdBook](https://rust-lang.github.io/mdBook/) — docs, see [docs/README.md](docs/README.md)

### Run the examples

```shell
task run_demo         # component demo (web server)
task run_todo         # small todo app built on state (cmd/toolgui-todo)
task run_wasm_hello   # hello world in the browser (WebAssembly)
task run_wasm_demo    # component demo in the browser
task run_desktop_hello # hello world as a desktop app
```

`toolgui-wails` is a separate Go module, since Wails v2 needs GTK and WebKit.
For your own desktop app, use `go tool toolgui-desktop build` or `run`; see
the [desktop build instructions](toolgui-wails/README.md#building). Tagged
releases include the frontend, so consumers need no Node or Yarn.
`cmd/toolgui-wasm` builds the static site for the WebAssembly runtime; the
[live demo](https://voilelab.github.io/toolgui/demo/) is published from it.

The `asset_app`, `asset_wails` and `asset_wasm` tasks each prepare their own
prerequisites. When building several frontends, pass them to one invocation:

```shell
task asset_app asset_wails asset_wasm
```

The workspace dependencies are installed once, and the shared library is built
and type-checked once, including its TypeScript declarations. Each frontend
then builds and runs its own typecheck.

### Build without the web assets

`toolgui-web/web.go` and `toolgui-wails/assets.go` embed build output that is
not committed on `dev`. Stub it to compile Go without a frontend build:

```shell
task stub_assets
go build ./...
```

### Test

```shell
task stub_assets
go test ./...
task test_wasm   # tgframe/tgwasm under js/wasm; needs Chrome on PATH
```

E2E against the server demo (start `task run_demo` first):

```shell
cd toolgui-e2e
yarn
yarn e2e:chrome
yarn e2e:firefox
```

E2E against the WebAssembly demo (start `task run_wasm_demo` first):

```shell
task test_e2e_wasm
```

## Release

Run the **Release** GitHub Action (`Actions` → `Release` → `Run workflow`) with
a `vX.Y.Z` version. It builds the web assets from `dev`, commits them, tags,
and pushes.

- `go get` resolves tags, so the tagged commit must carry
  `toolgui-web/app/build` — the only reason built assets live in git.
- `main` is a CI-owned mirror of the latest release; never commit to it by hand.
- Each release pushes two tags: `vX.Y.Z` and `toolgui-wails/vX.Y.Z`. Go resolves
  a submodule through a tag prefixed with its directory; without it, `go get`
  on the desktop module falls back to `dev`'s head, where the assets are
  gitignored.
