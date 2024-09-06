package tgexec

import (
	"encoding/json"
	"io"
	"log/slog"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"net/http"

	toolguiweb "github.com/mudream4869/toolgui/toolgui-web"
	"github.com/mudream4869/toolgui/toolgui/tgframe"
	"github.com/mudream4869/toolgui/toolgui/tgutil"

	"golang.org/x/net/websocket"
)

// TODO: Let it be configurable

// MaxUploadSize limit the size of file uploading form.
const MaxUploadSize int64 = 1024 * 1024 * 1024

// ErrUpdateInterrupt is raise at panic when current state is going to interrupt
var ErrUpdateInterrupt = tgutil.NewError("update interrupt")

// FIXME: use errors.Is or errors.As
const forceClosedByRemoteStr = "An existing connection was forcibly closed by the remote host."

// WebExecutor is a web ui executor for ToolGUI.
type WebExecutor struct {
	rootAssets map[string][]byte

	stateMap tgutil.UUIDMap[tgframe.State]

	app *tgframe.App
}

type stateValueChangeEvent struct {
	ID     string `json:"id"`
	Value  any    `json:"value"`
	IsTemp bool   `json:"is_temp"`
}

type stateIDPack struct {
	StateID string `json:"state_id"`
}

type resultPack struct {
	Error   string `json:"error,omitempty"`
	Success bool   `json:"success"`
}

type readyPack struct {
	Ready bool `json:"ready"`
}

// NewWebExecutor return a WebExecutor.
func NewWebExecutor(app *tgframe.App) *WebExecutor {
	return &WebExecutor{
		rootAssets: toolguiweb.GetRootAssets(),

		stateMap: tgutil.NewUUIDMap(
			tgframe.NewState, func(t *tgframe.State) { t.Destroy() },
			5*time.Minute),

		app: app,
	}
}

// Destory release all resource.
func (e *WebExecutor) Destroy() {
	e.stateMap.Destroy()
}

func (e *WebExecutor) handleUpdate(ws *websocket.Conn) {
	pageName := ws.Request().PathValue("name")
	if !e.app.HasPage(pageName) {
		websocket.JSON.Send(ws, &resultPack{
			Error:   "page not found",
			Success: false,
		})
		slog.Error("page not found", "page", pageName)
		return
	}

	var pack stateIDPack
	err := websocket.JSON.Receive(ws, &pack)
	if err != nil {
		websocket.JSON.Send(ws, &resultPack{
			Error:   err.Error(),
			Success: false,
		})
		slog.Error("state id", "error", err)
		return
	}

	var stateID string
	stateID = pack.StateID

	state, alive := e.stateMap.Get(stateID)
	if state == nil {
		stateID = e.stateMap.New()
		state, _ = e.stateMap.Get(stateID)
		websocket.JSON.Send(ws, stateIDPack{
			StateID: stateID,
		})
	} else {
		if alive {
			websocket.JSON.Send(ws, &resultPack{
				Error:   "state id already alive",
				Success: false,
			})
			slog.Error("state id already alive", "state_id", stateID)
			return
		}

		e.stateMap.SetAlive(stateID, true)
	}

	var stopUpdating atomic.Bool
	var running sync.Mutex

	for {
		var event stateValueChangeEvent
		err := websocket.JSON.Receive(ws, &event)
		if err != nil {
			if err == io.EOF || strings.Contains(err.Error(), forceClosedByRemoteStr) {
				// Connection closed

				// sending stop signal and wait for runner
				stopUpdating.Store(true)
				running.Lock()

				e.stateMap.SetAlive(stateID, false)
				break
			}

			websocket.JSON.Send(ws, &resultPack{
				Error:   err.Error(),
				Success: false,
			})
			slog.Error("state value change", "error", err)
			continue
		}

		// sending stop signal and wait for runner
		stopUpdating.Store(true)
		running.Lock()

		// tell client we cut the previous runner
		err = websocket.JSON.Send(ws, &readyPack{
			Ready: true,
		})
		if err != nil {
			slog.Error("send ready pack", "error", err)
		}

		// Clear temp state
		state.SetClickID("")

		if event.Value != nil {
			if event.IsTemp {
				// Only button click will send is_temp currently
				state.SetClickID(event.ID)
			} else {
				state.Set(event.ID, event.Value)
			}
		}

		sendNotifyPack := func(pack tgframe.NotifyPack) {
			if stopUpdating.Load() {
				panic(ErrUpdateInterrupt)
			}

			err := websocket.JSON.Send(ws, pack)
			if err != nil {
				panic(err)
			}
		}

		stopUpdating.Store(false)
		go func() {
			defer running.Unlock()
			err := e.app.RunWithHandlingPanic(pageName, state, sendNotifyPack)
			if err != nil {
				websocket.JSON.Send(ws, &resultPack{
					Error:   err.Error(),
					Success: false,
				})
				slog.Error("run err", "error", err)
			} else {
				websocket.JSON.Send(ws, &resultPack{
					Success: true,
				})
			}
		}()
	}
}

func (e *WebExecutor) handleUpload(w http.ResponseWriter, req *http.Request) {
	stateID := req.Header.Get("STATE_ID")
	state, alive := e.stateMap.Get(stateID)
	if state == nil || !alive {
		http.Error(w, "State ID is invalid or not alive", http.StatusForbidden)
		return
	}

	err := req.ParseMultipartForm(MaxUploadSize)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		slog.Error("Parse form", "error", err)
		return
	}

	file, handler, err := req.FormFile("file")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		slog.Error("Get formfile", "error", err)
		return
	}
	defer file.Close()

	bs, err := io.ReadAll(file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		slog.Error("Open file", "error", err)
		return
	}

	// TODO: Remove old file
	state.SetFile(handler.Filename, bs)
}

func (e *WebExecutor) handlePage(resp http.ResponseWriter, req *http.Request) {
	pageName := req.PathValue("name")
	body, isRootAssets := e.rootAssets[pageName]
	if isRootAssets {
		resp.Write(body)
		return
	}

	resp.Write([]byte(toolguiweb.IndexBody))
}

func (e *WebExecutor) handleAssets(resp http.ResponseWriter, req *http.Request) {
	pageName := req.PathValue("name")
	body, isRootAssets := e.rootAssets[pageName]
	if !isRootAssets {
		resp.WriteHeader(http.StatusNotFound)
		return
	}

	resp.Write(body)
}

func (e *WebExecutor) handleIndex(resp http.ResponseWriter, req *http.Request) {
	resp.Write([]byte(toolguiweb.IndexBody))
}

func (e *WebExecutor) handleHealth(resp http.ResponseWriter, req *http.Request) {
	resp.WriteHeader(http.StatusOK)
}

func (e *WebExecutor) handleAppConf(resp http.ResponseWriter, req *http.Request) {
	bs, err := json.Marshal(e.app.AppConf())
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp.Write(bs)
}

// Mux return a http mux to handle whole app.
//
//	mux, _ := e.Mux()
//	http.ListenAndServe(":8080", mux)
func (e *WebExecutor) Mux() (*http.ServeMux, error) {
	mux := http.NewServeMux()
	if e.app.AppConf().HashPageNameMode {
		mux.HandleFunc("GET /{name}", e.handleAssets)
		mux.HandleFunc("GET /", e.handleIndex)
	} else {
		mux.HandleFunc("GET /{name}", e.handlePage)
		if firstPage, ok := e.app.FirstPage(); ok {
			mux.Handle("GET /", http.RedirectHandler("/"+firstPage,
				http.StatusTemporaryRedirect))
		}
	}

	mux.Handle("GET /api/update/{name}", websocket.Handler(e.handleUpdate))
	mux.HandleFunc("POST /api/files", e.handleUpload)
	mux.HandleFunc("GET /api/app", e.handleAppConf)
	mux.HandleFunc("GET /api/health", e.handleHealth)

	mux.Handle("GET /static/", http.FileServerFS(toolguiweb.GetStaticDir()))

	return mux, nil
}

// StartService start serving the app at addr.
func (e *WebExecutor) StartService(addr string) error {
	mux, err := e.Mux()
	if err != nil {
		return tgutil.Errorf("%w", err)
	}

	err = http.ListenAndServe(addr, mux)
	if err != nil {
		return tgutil.Errorf("%w", err)
	}

	return nil
}
