package tcinput

import (
	"github.com/mudream4869/toolgui/toolgui/tgcomp/tcutil"
	"github.com/mudream4869/toolgui/toolgui/tgframe"
)

var _ tgframe.Component = &textboxComponent{}
var textboxComponentName = "textbox_component"

type textboxComponent struct {
	*tgframe.BaseComponent
	Label string `json:"label"`
}

func newTextboxComponent(label string) *textboxComponent {
	return &textboxComponent{
		BaseComponent: &tgframe.BaseComponent{
			Name: textboxComponentName,
			ID:   tcutil.NormalID(textboxComponentName, label),
		},
		Label: label,
	}
}

// Textbox create a textbox and return its value.
func Textbox(s *tgframe.State, c *tgframe.Container, label string) string {
	comp := newTextboxComponent(label)
	c.AddComponent(comp)
	return s.GetString(comp.ID)
}
