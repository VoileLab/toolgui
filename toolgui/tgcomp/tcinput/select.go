package tcinput

import (
	"github.com/mudream4869/toolgui/toolgui/framework"
	"github.com/mudream4869/toolgui/toolgui/tgcomp/tcutil"
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

// Select create a select dropdown list and return its selected value.
func Select(s *framework.State, c *framework.Container, label string, items []string) string {
	comp := newSelectComponent(label, items)
	c.AddComponent(comp)
	return s.GetString(comp.ID)
}
