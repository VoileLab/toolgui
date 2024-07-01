package tccontent

import (
	"github.com/mudream4869/toolgui/toolgui/component/tcutil"
	"github.com/mudream4869/toolgui/toolgui/framework"
)

var _ framework.Component = &subtitleComponent{}
var subtitleComponentName = "subtitle_component"

type subtitleComponent struct {
	*framework.BaseComponent
	Text string `json:"text"`
}

func newSubtitleComponent(text string) *subtitleComponent {
	return &subtitleComponent{
		BaseComponent: &framework.BaseComponent{
			Name: subtitleComponentName,
			ID:   tcutil.NormalID(subtitleComponentName, text),
		},
		Text: text,
	}
}

func Subtitle(c *framework.Container, text string) {
	c.AddComponent(newSubtitleComponent(text))
}
