package tcinput

import (
	"github.com/mudream4869/toolgui/toolgui/component/tcutil"
	"github.com/mudream4869/toolgui/toolgui/framework"
)

var _ framework.Component = &radioComponent{}
var radioComponentName = "radio_component"

type radioComponent struct {
	*framework.BaseComponent
	Label string   `json:"label"`
	Items []string `json:"items"`
}

func newRadioComponent(label string, items []string) *radioComponent {
	return &radioComponent{
		BaseComponent: &framework.BaseComponent{
			Name: radioComponentName,
			ID:   tcutil.NormalID(radioComponentName, label),
		},
		Label: label,
		Items: items,
	}
}

func Radio(sess *framework.State, c *framework.Container, label string, items []string) string {
	comp := newRadioComponent(label, items)
	c.AddComponent(comp)
	return sess.GetString(comp.ID)
}
