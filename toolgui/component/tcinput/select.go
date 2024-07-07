package tcinput

import (
	"github.com/mudream4869/toolgui/toolgui/component/tcutil"
	"github.com/mudream4869/toolgui/toolgui/framework"
)

var _ framework.Component = &selectComponent{}
var selectComponentName = "select_component"

type selectComponent struct {
	*framework.BaseComponent
	Label string   `json:"label"`
	Items []string `json:"items"`
}

func newSelectComponent(label string, items []string) *selectComponent {
	return &selectComponent{
		BaseComponent: &framework.BaseComponent{
			Name: selectComponentName,
			ID:   tcutil.NormalID(selectComponentName, label),
		},
		Label: label,
		Items: items,
	}
}

func Select(sess *framework.State, c *framework.Container, label string, items []string) string {
	comp := newSelectComponent(label, items)
	c.AddComponent(comp)
	return sess.GetString(comp.ID)
}
