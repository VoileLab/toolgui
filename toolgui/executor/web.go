package executor

import (
	"encoding/json"
	"fmt"
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

type RunFunc func(*framework.Session, *framework.Container) error

type PageConfig struct {
	Name  string `json:"name"`
	Title string `json:"title"`
	Emoji string `json:"emoji"`
}

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

type componentPack struct {
	ContainerID string              `json:"container_id"`
	Component   framework.Component `json:"component"`
}

type pageData struct {
	PageNames []string               `json:"page_names"`
	PageConfs map[string]*PageConfig `json:"page_confs"`
}

func NewWebExecutor() *WebExecutor {
	return &WebExecutor{
		rootAssets: toolguiweb.GetRootAssets(),

		sessions:  sessions.NewSessions(framework.NewSession),
		pageConfs: make(map[string]*PageConfig),
		pageFuncs: make(map[string]RunFunc),
	}
}

func (e *WebExecutor) AddPage(name, title string, runFunc RunFunc) error {
	return e.AddPageByConfig(&PageConfig{
		Name:  name,
		Title: title,
	}, runFunc)
}

func (e *WebExecutor) AddPageByConfig(conf *PageConfig, runFunc RunFunc) error {
	if conf == nil {
		return fmt.Errorf("nil config")
	}

	if conf.Name == "" || conf.Name == "api" || conf.Name == "static" {
		return fmt.Errorf("name should not be empty or 'api' or 'static'")
	}

	if _, exist := e.rootAssets[conf.Name]; exist {
		return fmt.Errorf("name duplicate with root assets")
	}

	if _, exist := e.pageConfs[conf.Name]; exist {
		return fmt.Errorf("name duplicate")
	}

	e.pageFuncs[conf.Name] = runFunc
	e.pageConfs[conf.Name] = conf
	e.pageNames = append(e.pageNames, conf.Name)

	return nil
}

func (e *WebExecutor) HandleUpdate(ws *websocket.Conn) {
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
		func(containerID string, comp framework.Component) {
			websocket.JSON.Send(ws, componentPack{
				ContainerID: containerID,
				Component:   comp,
			})
		})

	if event.IsTemp {
		sess = sess.Copy()
	}

	if event.Value != nil {
		sess.Set(event.ID, event.Value)
	}

	err := pageFunc(sess, newRoot)
	if err != nil {
		log.Println(err)
	}
}

func (e *WebExecutor) HandlePage(resp http.ResponseWriter, req *http.Request) {
	pageName := req.PathValue("name")
	body, isRootAssets := e.rootAssets[pageName]
	if isRootAssets {
		resp.Write(body)
		return
	}

	resp.Write([]byte(toolguiweb.IndexBody))
}

func (e *WebExecutor) HandlePageData(resp http.ResponseWriter, req *http.Request) {
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

func (e *WebExecutor) StartService() error {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{name}", e.HandlePage)
	mux.Handle("/api/update/{name}", websocket.Handler(e.HandleUpdate))
	mux.HandleFunc("GET /api/pages", e.HandlePageData)

	if len(e.pageConfs) != 0 {
		mux.Handle("/", http.RedirectHandler(
			"/"+e.pageNames[0], http.StatusTemporaryRedirect))
	}

	mux.Handle("/static/", http.FileServerFS(toolguiweb.GetStaticDir()))

	err := http.ListenAndServe(":3000", mux)
	if err != nil {
		return err
	}

	return nil
}
