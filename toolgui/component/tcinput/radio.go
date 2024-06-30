package tcinput

import (
	"github.com/mudream4869/toolgui/toolgui/component/tcutil"
	"github.com/mudream4869/toolgui/toolgui/framework"
)

var _ framework.Component = &RadioComponent{}
var RadioComponentName = "radio_component"

type RadioComponent struct {
	*framework.BaseComponent
	Label string   `json:"label"`
	Items []string `json:"items"`
}

func NewRadioComponent(label string, items []string) *RadioComponent {
	return &RadioComponent{
		BaseComponent: &framework.BaseComponent{
			Name: RadioComponentName,
			ID:   tcutil.NormalID(RadioComponentName, label),
		},
		Label: label,
		Items: items,
	}
}

func Radio(sess *framework.Session, c *framework.Container, label string, items []string) string {
	comp := NewRadioComponent(label, items)
	c.AddComponent(comp)
	return sess.GetString(comp.ID)
}
