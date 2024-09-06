package tccontent

import (
	"github.com/mudream4869/toolgui/toolgui/framework"
	"github.com/mudream4869/toolgui/toolgui/tgcomp/tcutil"
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

// Subtitle create a subtitle.
func Subtitle(c *framework.Container, text string) {
	c.AddComponent(newSubtitleComponent(text))
}

// SubtitleWithID create a subtitle component with a user specific id.
func SubtitleWithID(c *framework.Container, text string, id string) {
	comp := newSubtitleComponent(text)
	comp.SetID(id)
	c.AddComponent(comp)
}
