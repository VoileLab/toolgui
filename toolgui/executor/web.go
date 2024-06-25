package executor

import (
	"encoding/json"
	"errors"
	"log"

	"net/http"

	toolguiweb "github.com/mudream4869/toolgui/toolgui-web"
	"github.com/mudream4869/toolgui/toolgui/framework"
	"github.com/mudream4869/toolgui/toolgui/sessions"

	"golang.org/x/net/websocket"
)

// ROOT_CONTAINER_ID is the id of root container.
// The creation of root container won't trigger NotifyAddCompFunc.
const ROOT_CONTAINER_ID string = "container_root"

// RunFunc is the type of a function handling page
type RunFunc func(*framework.Session, *framework.Container) error

// PageConfig stores basic setting of a page
type PageConfig struct {
	// Name should not duplicate to another page
	Name string `json:"name"`

	// Title will show as the title of page
	Title string `json:"title"`

	// Emoji will show as icon of a page
	Emoji string `json:"emoji"`
}

// WebExecutor is a web ui executor for ToolGUI
type WebExecutor struct {
	rootAssets map[string][]byte

	pageNames []string
	pageConfs map[string]*PageConfig
	pageFuncs map[string]RunFunc

	sessions sessions.Sessions[framework.Session]
}

type sessionValueChangeEvent struct {
	SessionID string `json:"session_id"`
	ID        string `json:"id"`
	Value     any    `json:"value"`
	IsTemp    bool   `json:"is_temp"`
}

type sessionPack struct {
	SessionID string `json:"session_id"`
}

type resultPack struct {
	Error   string `json:"error,omitempty"`
	Success bool   `json:"success"`
}

type pageData struct {
	PageNames []string               `json:"page_names"`
	PageConfs map[string]*PageConfig `json:"page_confs"`
}

// NewWebExecutor return a WebExecutor
func NewWebExecutor() *WebExecutor {
	return &WebExecutor{
		rootAssets: toolguiweb.GetRootAssets(),

		sessions:  sessions.NewSessions(framework.NewSession),
		pageConfs: make(map[string]*PageConfig),
		pageFuncs: make(map[string]RunFunc),
	}
}

// AddPage add a handled page by name, title, and runFunc
//
//	e := NewWebExecutor()
//	e.AddPage("index", "Index", func(s *framework.Session, c *framework.Container) error {
//		component.Text(c, "hello world")
//		return nil
//	})
func (e *WebExecutor) AddPage(name, title string, runFunc RunFunc) error {
	return e.AddPageByConfig(&PageConfig{
		Name:  name,
		Title: title,
	}, runFunc)
}

// AddPageByConfig add a handled page by name, title, icon, and runFunc
//
//	e := NewWebExecutor()
//	e.AddPage(e.AddPageByConfig(&executor.PageConfig{
//		Name:  "page1",
//		Title: "Page1",
//		Emoji: "🐱",
//	}, Page1)
func (e *WebExecutor) AddPageByConfig(conf *PageConfig, runFunc RunFunc) error {
	if conf == nil {
		return errors.New("nil config")
	}

	if conf.Name == "" || conf.Name == "api" || conf.Name == "static" {
		return errors.New("name should not be empty or 'api' or 'static'")
	}

	if _, exist := e.rootAssets[conf.Name]; exist {
		return errors.New("name duplicate with root assets")
	}

	if _, exist := e.pageConfs[conf.Name]; exist {
		return errors.New("name duplicate")
	}

	e.pageFuncs[conf.Name] = runFunc
	e.pageConfs[conf.Name] = conf
	e.pageNames = append(e.pageNames, conf.Name)

	return nil
}

func (e *WebExecutor) handleUpdate(ws *websocket.Conn) {
	pageName := ws.Request().PathValue("name")
	pageFunc, ok := e.pageFuncs[pageName]
	if !ok {
		log.Println("Not found", pageName)
		return
	}

	var event sessionValueChangeEvent
	websocket.JSON.Receive(ws, &event)

	sess := e.sessions.Get(event.SessionID)
	if sess == nil {
		sessID := e.sessions.New()
		sess = e.sessions.Get(sessID)
		websocket.JSON.Send(ws, sessionPack{
			SessionID: sessID,
		})
		event.Value = nil
	}

	newRoot := framework.NewContainer(ROOT_CONTAINER_ID,
		func(pack framework.NotifyPack) {
			err := websocket.JSON.Send(ws, pack)
			if err != nil {
				log.Println(err)
			}
		})

	if event.IsTemp {
		sess = sess.Copy()
	}

	if event.Value != nil {
		sess.Set(event.ID, event.Value)
	}

	err := pageFunc(sess, newRoot)
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

func (e *WebExecutor) handlePageData(resp http.ResponseWriter, req *http.Request) {
	bs, err := json.Marshal(pageData{
		PageNames: e.pageNames,
		PageConfs: e.pageConfs,
	})

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
	if len(e.pageConfs) == 0 {
		return nil, errors.New("no register page")
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{name}", e.handlePage)
	mux.Handle("/api/update/{name}", websocket.Handler(e.handleUpdate))
	mux.HandleFunc("GET /api/pages", e.handlePageData)

	mux.Handle("/", http.RedirectHandler(
		"/"+e.pageNames[0], http.StatusTemporaryRedirect))

	mux.Handle("/static/", http.FileServerFS(toolguiweb.GetStaticDir()))

	return mux, nil
}

// StartService start serving the app at addr
func (e *WebExecutor) StartService(addr string) error {
	mux, err := e.Mux()
	if err != nil {
		return err
	}

	err = http.ListenAndServe(addr, mux)
	if err != nil {
		return err
	}

	return nil
}
