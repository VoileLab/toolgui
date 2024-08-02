package tcinput

import (
	"github.com/mudream4869/toolgui/toolgui/component/tcutil"
	"github.com/mudream4869/toolgui/toolgui/framework"
)

var _ framework.Component = &checkboxComponent{}
var checkboxComponentName = "checkbox_component"

type checkboxComponent struct {
	*framework.BaseComponent
	Label string `json:"label"`
}

func newCheckboxComponent(label string) *checkboxComponent {
	return &checkboxComponent{
		BaseComponent: &framework.BaseComponent{
			Name: checkboxComponentName,
			ID:   tcutil.NormalID(checkboxComponentName, label),
		},
		Label: label,
	}
}

// Checkbox create a checkbox and return true if it's clicked.
func Checkbox(s *framework.State, c *framework.Container, label string) bool {
	comp := newCheckboxComponent(label)
	c.AddComponent(comp)
	return s.GetBool(comp.ID)
}
