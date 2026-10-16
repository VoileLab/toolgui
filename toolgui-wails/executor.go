// Package tgwails runs a ToolGUI app in a desktop window with Wails v2.
//
// It's the desktop counterpart of tgexec.WebExecutor: same app, same pages,
// same frontend components, with bound methods and Wails events in place of
// HTTP and websockets.
//
//	app := tgframe.NewApp()
//	app.AddPage("index", "Index", Main)
//	tgwails.NewExecutor(app, &tgwails.Conf{Title: "My Tool"}).Run()
package tgwails

import (
	"github.com/voilelab/toolgui/internal/tgutil"
	"github.com/voilelab/toolgui/toolgui/tgframe"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

// RGBA is a window colour. [options.NewRGBA] builds one.
type RGBA = options.RGBA

// WindowStartState is the initial state of the desktop window.
type WindowStartState = options.WindowStartState

// Window start states, matching Wails.
const (
	Normal     WindowStartState = options.Normal
	Maximised  WindowStartState = options.Maximised
	Minimised  WindowStartState = options.Minimised
	Fullscreen WindowStartState = options.Fullscreen
)

// Conf configures the desktop window. The zero value is usable: see
// [DefaultConf] for what the empty fields fall back to.
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

// DefaultConf return the config used for the fields a caller leaves empty.
func DefaultConf() *Conf {
	return &Conf{
		Title:      "ToolGUI",
		Width:      1024,
		Height:     768,
		Background: options.NewRGBA(255, 255, 255, 1),
	}
}

func (c *Conf) withDefaults() *Conf {
	def := DefaultConf()
	if c == nil {
		return def
	}

	filled := *c
	if filled.Title == "" {
		filled.Title = def.Title
	}
	if filled.Width == 0 {
		filled.Width = def.Width
	}
	if filled.Height == 0 {
		filled.Height = def.Height
	}
	if filled.Background == nil {
		filled.Background = def.Background
	}

	return &filled
}

// Executor is a desktop ui executor for ToolGUI.
type Executor struct {
	app  *tgframe.App
	conf *Conf
}

// NewExecutor return an Executor serving app in a desktop window.
// A nil conf uses [DefaultConf].
func NewExecutor(app *tgframe.App, conf *Conf) *Executor {
	filled := conf.withDefaults()

	// The app title names the window too, so an app that sets one doesn't
	// repeat it here. A title in the conf still wins.
	if (conf == nil || conf.Title == "") && app.AppConf().Title != "" {
		filled.Title = app.AppConf().Title
	}

	return &Executor{
		app:  app,
		conf: filled,
	}
}

// Run open the window and block until the user closes it.
func (e *Executor) Run() error {
	backend := NewToolGUI(e.app)

	err := wails.Run(e.appOptions(backend))
	if err != nil {
		return tgutil.Errorf("%w", err)
	}

	return nil
}

func (e *Executor) appOptions(backend *ToolGUI) *options.App {
	var toggleFullscreen func()
	if e.conf.FullscreenToggle {
		toggleFullscreen = backend.toggleFullscreen
	}

	return &options.App{
		Title:            e.conf.Title,
		Width:            e.conf.Width,
		Height:           e.conf.Height,
		WindowStartState: e.conf.WindowStartState,
		BackgroundColour: e.conf.Background,
		DisableResize:    e.conf.DisableResize,
		Frameless:        e.conf.Frameless,

		// The tree the app declared becomes the window's own menubar, rather
		// than a row the frontend draws inside it. The way back needs nothing
		// new: the click callback happens in Go, so it goes straight into the
		// session as the click event the web menubar would have sent over the
		// wire, and the packs come back the usual way.
		Menu: nativeMenu(e.conf.Frameless, e.app.AppConf().Menu,
			backend.queueMenuClick, toggleFullscreen),

		// Handler catches what the embedded frontend does not have, which is
		// where plugin assets land.
		AssetServer: &assetserver.Options{
			Assets:  Assets,
			Handler: tgframe.PluginAssetHandler(e.app),
		},

		OnStartup:  backend.start,
		OnShutdown: backend.shutdown,
		Bind:       []any{backend},
	}
}
