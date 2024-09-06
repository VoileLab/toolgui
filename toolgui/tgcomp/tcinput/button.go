package tcinput

import (
	"github.com/mudream4869/toolgui/toolgui/framework"
	"github.com/mudream4869/toolgui/toolgui/tgcomp/tcutil"
)

var _ framework.Component = &buttonComponent{}
var buttonComponentName = "button_component"

type buttonComponent struct {
	*framework.BaseComponent
	Label string `json:"label"`
}

func newButtonComponent(label string) *buttonComponent {
	return &buttonComponent{
		BaseComponent: &framework.BaseComponent{
			Name: buttonComponentName,
			ID:   tcutil.NormalID(buttonComponentName, label),
		},
		Label: label,
	}
}

// Button create a button and return true if it's clicked.
func Button(s *framework.State, c *framework.Container, label string) bool {
	comp := newButtonComponent(label)
	c.AddComponent(comp)
	return s.GetClickID() == comp.ID
}
