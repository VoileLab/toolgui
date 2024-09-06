package tcinput

import (
	"github.com/mudream4869/toolgui/toolgui/tgcomp/tcutil"
	"github.com/mudream4869/toolgui/toolgui/tgframe"
)

var _ tgframe.Component = &buttonComponent{}
var buttonComponentName = "button_component"

type buttonComponent struct {
	*tgframe.BaseComponent
	Label string `json:"label"`
}

func newButtonComponent(label string) *buttonComponent {
	return &buttonComponent{
		BaseComponent: &tgframe.BaseComponent{
			Name: buttonComponentName,
			ID:   tcutil.NormalID(buttonComponentName, label),
		},
		Label: label,
	}
}

// Button create a button and return true if it's clicked.
func Button(s *tgframe.State, c *tgframe.Container, label string) bool {
	comp := newButtonComponent(label)
	c.AddComponent(comp)
	return s.GetClickID() == comp.ID
}
