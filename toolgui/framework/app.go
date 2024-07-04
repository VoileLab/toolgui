package framework

import (
	"errors"
	"fmt"
)

var ErrPageNotFound = errors.New("page not found")

// RootContainerID is the id of root container.
// The creation of root container won't trigger SendNotifyPackFunc.
const RootContainerID string = "container_root"

func realRootContainerID() string {
	return NewContainer(RootContainerID, nil).ID
}

// SidebarContainerID is the id of sidebar container.
// The creation of root container won't trigger SendNotifyPackFunc.
const SidebarContainerID string = "container_sidebar"

func realSidebarContainerID() string {
	return NewContainer(SidebarContainerID, nil).ID
}

// RunFunc is the type of a function handling page
type RunFunc func(*Session, *Container, *Container) error

// PageConfig stores basic setting of a page
type PageConfig struct {
	// Name should not duplicate to another page
	Name string `json:"name"`

	// Title will show as the title of page
	Title string `json:"title"`

	// Emoji will show as icon of a page
	Emoji string `json:"emoji"`
}

// App is an app
type App struct {
	pageNames []string
	pageConfs map[string]*PageConfig
	pageFuncs map[string]RunFunc

	hashPageNameMode bool
}

// AppConf
type AppConf struct {
	PageNames []string               `json:"page_names"`
	PageConfs map[string]*PageConfig `json:"page_confs"`

	HashPageNameMode bool `json:"hash_page_name_mode"`

	RootContainerID    string `json:"root_container_id"`
	SidebarContainerID string `json:"sidebar_container_id"`
}

// NewApp return App
func NewApp() *App {
	return &App{
		pageNames: make([]string, 0),
		pageConfs: make(map[string]*PageConfig),
		pageFuncs: make(map[string]RunFunc),
	}
}

// SetHashPageMode set value of hash page name mode flag.
func (app *App) SetHashPageNameMode(v bool) {
	app.hashPageNameMode = v
}

// AddPage add a handled page by name, title, and runFunc
//
//	app.AddPage("index", "Index", f})
func (app *App) AddPage(name, title string, runFunc RunFunc) {
	err := app.addPageByConfig(&PageConfig{
		Name:  name,
		Title: title,
	}, runFunc)
	if err != nil {
		panic(err)
	}
}

// AddPageByConfig add a handled page by name, title, icon, and runFunc
//
//	app.AddPageByConfig(&framework.PageConfig{
//		Name:  "page1",
//		Title: "Page1",
//		Emoji: "🐱",
//	}, Page1)
func (app *App) AddPageByConfig(conf *PageConfig, runFunc RunFunc) {
	err := app.addPageByConfig(conf, runFunc)
	if err != nil {
		panic(err)
	}
}

func (app *App) addPageByConfig(conf *PageConfig, runFunc RunFunc) error {
	if conf == nil {
		return errors.New("nil config")
	}

	if conf.Name == "" || conf.Name == "api" || conf.Name == "static" {
		return errors.New("name should not be empty or 'api' or 'static'")
	}

	if _, exist := app.pageConfs[conf.Name]; exist {
		return errors.New("name duplicate")
	}

	app.pageFuncs[conf.Name] = runFunc
	app.pageConfs[conf.Name] = conf
	app.pageNames = append(app.pageNames, conf.Name)

	return nil
}

// AppConf return App Config
func (app *App) AppConf() *AppConf {
	return &AppConf{
		PageNames: app.pageNames,
		PageConfs: app.pageConfs,

		RootContainerID:    realRootContainerID(),
		SidebarContainerID: realSidebarContainerID(),

		HashPageNameMode: app.hashPageNameMode,
	}
}

// Run run a page which named `name` with sess
func (app *App) Run(name string, sess *Session, notifyFunc SendNotifyPackFunc) error {
	pageFunc, ok := app.pageFuncs[name]
	if !ok {
		return fmt.Errorf("%w: `%s`", ErrPageNotFound, name)
	}

	newRoot := NewContainer(RootContainerID, notifyFunc)
	newSidebar := NewContainer(SidebarContainerID, notifyFunc)

	err := pageFunc(sess, newRoot, newSidebar)
	if err != nil {
		return err
	}

	return nil
}

// HasPage return existence of page which named `name`
func (app *App) HasPage(name string) bool {
	_, ok := app.pageFuncs[name]
	return ok
}

// FirstPage return the first page in the app
func (app *App) FirstPage() (string, bool) {
	if len(app.pageNames) == 0 {
		return "", false
	}
	return app.pageNames[0], true
}
