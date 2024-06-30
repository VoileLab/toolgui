package tccontent

import (
	"github.com/mudream4869/toolgui/toolgui/component/tcutil"
	"github.com/mudream4869/toolgui/toolgui/framework"
)

var _ framework.Component = &SubtitleComponent{}
var SubtitleComponentName = "subtitle_component"

type SubtitleComponent struct {
	*framework.BaseComponent
	Text string `json:"text"`
}

func NewSubtitleComponent(text string) *SubtitleComponent {
	return &SubtitleComponent{
		BaseComponent: &framework.BaseComponent{
			Name: SubtitleComponentName,
			ID:   tcutil.NormalID(SubtitleComponentName, text),
		},
		Text: text,
	}
}

func Subtitle(c *framework.Container, text string) {
	c.AddComponent(NewSubtitleComponent(text))
}
