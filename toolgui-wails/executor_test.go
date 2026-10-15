package tgwails

import (
	"testing"

	"github.com/voilelab/toolgui/toolgui/tgframe"
	"github.com/wailsapp/wails/v2/pkg/options"
)

func TestConfDefaults(t *testing.T) {
	def := DefaultConf()

	conf := (&Conf{Title: "Mine"}).withDefaults()
	if conf.Title != "Mine" {
		t.Fatalf("expect the given title, got %q", conf.Title)
	}
	if conf.Width != def.Width || conf.Height != def.Height {
		t.Fatalf("expect the default size, got %dx%d", conf.Width, conf.Height)
	}
	// Compare the colours, not the pointers.
	if *conf.Background != *def.Background {
		t.Fatalf("expect the default background, got %v", conf.Background)
	}

	// A window is resizable unless the caller opts out, so the zero value of
	// Conf is a usable config.
	if conf.DisableResize {
		t.Fatal("expect a resizable window by default")
	}

	if (*Conf)(nil).withDefaults().Title != def.Title {
		t.Fatal("expect a nil conf to fall back to the defaults")
	}
}

func TestConfWindowStartState(t *testing.T) {
	if got := DefaultConf().WindowStartState; got != Normal {
		t.Fatalf("default WindowStartState = %v, want Normal", got)
	}
	for _, tt := range []struct {
		name string
		conf *Conf
		want WindowStartState
	}{
		{name: "nil", want: Normal},
		{name: "zero", conf: &Conf{}, want: Normal},
		{name: "omitted", conf: &Conf{Width: 1200, Height: 800}, want: Normal},
		{name: "normal", conf: &Conf{WindowStartState: Normal}, want: Normal},
		{name: "maximised", conf: &Conf{WindowStartState: Maximised}, want: Maximised},
		{name: "minimised", conf: &Conf{WindowStartState: Minimised}, want: Minimised},
		{name: "fullscreen", conf: &Conf{WindowStartState: Fullscreen}, want: Fullscreen},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.conf.withDefaults().WindowStartState; got != tt.want {
				t.Errorf("WindowStartState = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExecutorWindowOptions(t *testing.T) {
	for _, tt := range []struct {
		name  string
		state WindowStartState
		want  options.WindowStartState
	}{
		{"normal", Normal, options.Normal},
		{"maximised", Maximised, options.Maximised},
		{"minimised", Minimised, options.Minimised},
		{"fullscreen", Fullscreen, options.Fullscreen},
	} {
		for _, disableResize := range []bool{false, true} {
			name := tt.name
			if disableResize {
				name += "/disable-resize"
			}
			t.Run(name, func(t *testing.T) {
				app := tgframe.NewApp()
				e := NewExecutor(app, &Conf{
					Width: 1200, Height: 800,
					WindowStartState: tt.state,
					DisableResize:    disableResize,
				})
				got := e.appOptions(NewToolGUI(app))
				if got.WindowStartState != tt.want {
					t.Errorf("WindowStartState = %v, want %v", got.WindowStartState, tt.want)
				}
				if got.Width != 1200 || got.Height != 800 {
					t.Errorf("window size = %dx%d, want 1200x800", got.Width, got.Height)
				}
				if got.DisableResize != disableResize {
					t.Errorf("DisableResize = %v, want %v", got.DisableResize, disableResize)
				}
			})
		}
	}
}

// The app title names the window, so a title is set in one place for both
// executors. A conf that gives its own still wins.
func TestExecutorTitleFromApp(t *testing.T) {
	app := tgframe.NewApp()
	app.SetTitle("My Tool")

	if got := NewExecutor(app, nil).conf.Title; got != "My Tool" {
		t.Errorf("Title = %q, want %q", got, "My Tool")
	}

	if got := NewExecutor(app, &Conf{Title: "Mine"}).conf.Title; got != "Mine" {
		t.Errorf("Title = %q, want %q", got, "Mine")
	}

	// An app with no title falls back to the default.
	if got := NewExecutor(tgframe.NewApp(), nil).conf.Title; got != DefaultConf().Title {
		t.Errorf("Title = %q, want %q", got, DefaultConf().Title)
	}
}
