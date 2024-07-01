package tccontent

import (
	"github.com/mudream4869/toolgui/toolgui/component/tcutil"
	"github.com/mudream4869/toolgui/toolgui/framework"
)

var _ framework.Component = &textComponent{}
var textComponentName = "text_component"

type textComponent struct {
	*framework.BaseComponent
	Text string `json:"text"`
}

func newTextComponent(text string) *textComponent {
	return &textComponent{
		BaseComponent: &framework.BaseComponent{
			Name: textComponentName,
			ID:   tcutil.NormalID(textComponentName, text),
		},
		Text: text,
	}
}

func Text(c *framework.Container, text string) {
	comp := newTextComponent(text)
	c.AddComponent(comp)
}
