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

// Text show a text.
func Text(c *framework.Container, text string) {
	comp := newTextComponent(text)
	c.AddComponent(comp)
}

// TextWithID create a text component with a user specific id.
func TextWithID(c *framework.Container, text string, id string) {
	comp := newTextComponent(text)
	comp.SetID(id)
	c.AddComponent(comp)
}
