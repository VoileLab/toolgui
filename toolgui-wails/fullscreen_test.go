package tgwails

import (
	"context"
	"slices"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/voilelab/toolgui/toolgui/tgframe"
)

func TestToolGUIToggleFullscreen(t *testing.T) {
	for _, initiallyFullscreen := range []bool{false, true} {
		name := "enter-first"
		if initiallyFullscreen {
			name = "exit-first"
		}
		t.Run(name, func(t *testing.T) {
			backend := NewToolGUI(tgframe.NewApp())
			// A window action must not need a session host at all.
			backend.host = nil
			ctx := t.Context()
			backend.start(ctx)
			fullscreen := initiallyFullscreen
			var calls []string
			checkContext := func(got context.Context) {
				if got != ctx {
					t.Fatal("runtime did not receive the startup context")
				}
			}
			backend.windowIsFullscreen = func(ctx context.Context) bool {
				checkContext(ctx)
				calls = append(calls, "is-fullscreen")
				return fullscreen
			}
			backend.windowFullscreen = func(ctx context.Context) {
				checkContext(ctx)
				calls = append(calls, "enter")
				fullscreen = true
			}
			backend.windowUnfullscreen = func(ctx context.Context) {
				checkContext(ctx)
				calls = append(calls, "exit")
				fullscreen = false
			}
			backend.toggleFullscreen()
			backend.toggleFullscreen()
			want := []string{"is-fullscreen", "enter", "is-fullscreen", "exit"}
			if initiallyFullscreen {
				want = []string{"is-fullscreen", "exit", "is-fullscreen", "enter"}
			}
			if !slices.Equal(calls, want) || fullscreen != initiallyFullscreen {
				t.Fatalf("calls = %v, fullscreen = %t; want %v, %t", calls, fullscreen, want, initiallyFullscreen)
			}
		})
	}
}

func TestToolGUIToggleFullscreenNotReady(t *testing.T) {
	for _, phase := range []string{"before-startup", "cancelled", "after-shutdown"} {
		t.Run(phase, func(t *testing.T) {
			backend := NewToolGUI(tgframe.NewApp())
			backend.windowIsFullscreen = func(context.Context) bool {
				t.Fatal("window state queried without a live runtime context")
				return false
			}
			backend.windowFullscreen = func(context.Context) { t.Fatal("entered fullscreen without a live context") }
			backend.windowUnfullscreen = func(context.Context) { t.Fatal("exited fullscreen without a live context") }
			switch phase {
			case "cancelled":
				ctx, cancel := context.WithCancel(t.Context())
				backend.start(ctx)
				cancel()
			case "after-shutdown":
				backend.start(t.Context())
				backend.shutdown(t.Context())
			}
			backend.toggleFullscreen()
		})
	}
}

func TestExecutorFullscreenToggle(t *testing.T) {
	for _, enabled := range []bool{false, true} {
		name := "disabled"
		if enabled {
			name = "enabled"
		}
		t.Run(name, func(t *testing.T) {
			var runs atomic.Int32
			app := newTestApp(func(*tgframe.Params) error { runs.Add(1); return nil })
			backend, events := newTestToolGUI(t, app)
			defer backend.shutdown(t.Context())
			e := NewExecutor(app, &Conf{FullscreenToggle: enabled})
			if e.conf.FullscreenToggle != enabled {
				t.Fatal("withDefaults changed FullscreenToggle")
			}
			native := e.appOptions(backend).Menu
			var toggle func()
			if native != nil {
				for _, item := range native.Items {
					if item.Label == "View" {
						toggle = func() { click(item.SubMenu.Items[0]) }
					}
				}
			}
			if !enabled {
				if toggle != nil {
					t.Fatal("disabled config added a fullscreen command")
				}
				return
			}
			if toggle == nil {
				t.Fatal("enabled config did not wire up the native toggle")
			}
			// Before OnStartup the actual menu callback must safely do nothing.
			toggle()
			backend.start(t.Context())
			backend.emit = events.emit
			entered := 0
			backend.windowIsFullscreen = func(context.Context) bool { return false }
			backend.windowFullscreen = func(context.Context) { entered++ }
			if err := backend.Start(testPageName, ""); err != nil {
				t.Fatal(err)
			}
			events.waitResult(t)
			toggle()
			if entered != 1 || runs.Load() != 1 || len(events.packs) != 0 {
				t.Fatalf("entered = %d, page runs = %d, pending packs = %d", entered, runs.Load(), len(events.packs))
			}
			backend.menuLock.Lock()
			defer backend.menuLock.Unlock()
			if backend.menuBusy || len(backend.menuQueue) != 0 {
				t.Fatal("fullscreen toggle entered the session menu queue")
			}
		})
	}
	if DefaultConf().FullscreenToggle || NewExecutor(tgframe.NewApp(), nil).conf.FullscreenToggle {
		t.Fatal("fullscreen toggle must be opt-in")
	}
}

func TestToolGUIToggleFullscreenConcurrent(t *testing.T) {
	backend := NewToolGUI(tgframe.NewApp())
	var fullscreen bool
	queries := 0
	backend.windowIsFullscreen = func(context.Context) bool { queries++; return fullscreen }
	backend.windowFullscreen = func(context.Context) { fullscreen = true }
	backend.windowUnfullscreen = func(context.Context) { fullscreen = false }
	var wg sync.WaitGroup
	wg.Add(1)
	go func() { defer wg.Done(); backend.start(t.Context()) }()
	for range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range 10 {
				backend.toggleFullscreen()
			}
		}()
	}
	wg.Wait()
	// Some calls can precede startup and be ignored; all others are serialized.
	if fullscreen != (queries%2 != 0) {
		t.Fatalf("fullscreen = %t after %d queries", fullscreen, queries)
	}
}
