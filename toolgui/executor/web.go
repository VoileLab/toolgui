package executor

import (
	"encoding/json"
	"io"
	"log"
	"time"

	"net/http"

	toolguiweb "github.com/mudream4869/toolgui/toolgui-web"
	"github.com/mudream4869/toolgui/toolgui/framework"
	"github.com/mudream4869/toolgui/toolgui/tgutil"

	"golang.org/x/net/websocket"
)

// WebExecutor is a web ui executor for ToolGUI
type WebExecutor struct {
	rootAssets map[string][]byte

	stateMap tgutil.UUIDMap[framework.State]

	app *framework.App
}

type stateValueChangeEvent struct {
	StateID string `json:"state_id"`
	ID      string `json:"id"`
	Value   any    `json:"value"`
	IsTemp  bool   `json:"is_temp"`
}

type healthEvent struct {
	Stop    bool   `json:"stop"`
	StateID string `json:"state_id"`
}

type statePack struct {
	StateID string `json:"state_id"`
}

type resultPack struct {
	Error   string `json:"error,omitempty"`
	Success bool   `json:"success"`
}

// NewWebExecutor return a WebExecutor
func NewWebExecutor(app *framework.App) *WebExecutor {
	return &WebExecutor{
		rootAssets: toolguiweb.GetRootAssets(),

		stateMap: tgutil.NewUUIDMap(
			framework.NewState, func(t *framework.State) { t.Destroy() },
			5*time.Minute),

		app: app,
	}
}

// Destory release all resource
func (e *WebExecutor) Destroy() {
	e.stateMap.Destroy()
}

func (e *WebExecutor) handleHealth(ws *websocket.Conn) {
	pageName := ws.Request().PathValue("name")
	if !e.app.HasPage(pageName) {
		websocket.JSON.Send(ws, &resultPack{
			Error:   "page not found",
			Success: false,
		})
		log.Println("Not found", pageName)
		return
	}

	for {
		var event healthEvent
		err := websocket.JSON.Receive(ws, &event)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Println(err)
			continue
		}

		if event.Stop {
			break
		}

		log.Println("State: ", event.StateID)

		e.stateMap.Get(event.StateID)
	}
}

func (e *WebExecutor) handleUpdate(ws *websocket.Conn) {
	pageName := ws.Request().PathValue("name")
	if !e.app.HasPage(pageName) {
		websocket.JSON.Send(ws, &resultPack{
			Error:   "page not found",
			Success: false,
		})
		log.Println("Not found", pageName)
		return
	}

	var event stateValueChangeEvent
	err := websocket.JSON.Receive(ws, &event)
	if err != nil {
		websocket.JSON.Send(ws, &resultPack{
			Error:   err.Error(),
			Success: false,
		})
		log.Println(err)
		return
	}

	state := e.stateMap.Get(event.StateID)
	if state == nil {
		stateID := e.stateMap.New()
		state = e.stateMap.Get(stateID)
		websocket.JSON.Send(ws, statePack{
			StateID: stateID,
		})
		event.Value = nil
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

	sendNotifyPack := func(pack framework.NotifyPack) {
		err := websocket.JSON.Send(ws, pack)
		if err != nil {
			log.Println(err)
		}
	}

	err = e.app.Run(pageName, state, sendNotifyPack)
	if err != nil {
		websocket.JSON.Send(ws, &resultPack{
			Error:   err.Error(),
			Success: false,
		})
		log.Println(err)
		return
	}

	websocket.JSON.Send(ws, &resultPack{
		Success: true,
	})
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

func (e *WebExecutor) handleAppConf(resp http.ResponseWriter, req *http.Request) {
	bs, err := json.Marshal(e.app.AppConf())
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp.Write(bs)
}

// Mux return a http mux to handle whole app
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
	mux.Handle("GET /api/health/{name}", websocket.Handler(e.handleHealth))
	mux.HandleFunc("GET /api/app", e.handleAppConf)

	mux.Handle("GET /static/", http.FileServerFS(toolguiweb.GetStaticDir()))

	return mux, nil
}

// StartService start serving the app at addr
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
