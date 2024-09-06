package tcinput

import (
	"github.com/mudream4869/toolgui/toolgui/framework"
	"github.com/mudream4869/toolgui/toolgui/tgcomp/tcutil"
)

var _ framework.Component = &textboxComponent{}
var textboxComponentName = "textbox_component"

type textboxComponent struct {
	*framework.BaseComponent
	Label string `json:"label"`
}

func newTextboxComponent(label string) *textboxComponent {
	return &textboxComponent{
		BaseComponent: &framework.BaseComponent{
			Name: textboxComponentName,
			ID:   tcutil.NormalID(textboxComponentName, label),
		},
		Label: label,
	}
}

// Textbox create a textbox and return its value.
func Textbox(s *framework.State, c *framework.Container, label string) string {
	comp := newTextboxComponent(label)
	c.AddComponent(comp)
	return s.GetString(comp.ID)
}
