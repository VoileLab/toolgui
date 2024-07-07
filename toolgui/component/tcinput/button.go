package tcinput

import (
	"github.com/mudream4869/toolgui/toolgui/component/tcutil"
	"github.com/mudream4869/toolgui/toolgui/framework"
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

func Button(sess *framework.State, c *framework.Container, label string) bool {
	comp := newButtonComponent(label)
	c.AddComponent(comp)
	return sess.GetBool(comp.ID)
}
