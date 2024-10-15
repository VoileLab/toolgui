package tccontent

import (
	"github.com/VoileLab/toolgui/toolgui/tgcomp/tcutil"
	"github.com/VoileLab/toolgui/toolgui/tgframe"
)

var _ tgframe.Component = &htmlComponent{}
var htmlComponentName = "html_component"

type htmlComponent struct {
	*tgframe.BaseComponent
	Html   string `json:"html"`
	Script bool   `json:"script"`
}

func newHtmlComponent(html string, script bool) *htmlComponent {
	return &htmlComponent{
		BaseComponent: &tgframe.BaseComponent{
			Name: htmlComponentName,
			ID:   tcutil.HashedID(htmlComponentName, []byte(html)),
		},
		Html:   html,
		Script: script,
	}
}

// Html show a html.
// script is used to allow the iframe to run javascript. (notice that this is not secure)
func Html(c *tgframe.Container, html string, script bool) {
	comp := newHtmlComponent(html, script)
	c.AddComponent(comp)
}

// HtmlWithID create a html component with a user specific id.
// script is used to allow the iframe to run javascript. (notice that this is not secure)
func HtmlWithID(c *tgframe.Container, html string, id string, script bool) {
	comp := newHtmlComponent(html, script)
	comp.SetID(id)
	c.AddComponent(comp)
}
