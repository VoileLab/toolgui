package tcinput

import (
	"github.com/mudream4869/toolgui/toolgui/tgcomp/tcutil"
	"github.com/mudream4869/toolgui/toolgui/tgframe"
)

var _ tgframe.Component = &checkboxComponent{}
var checkboxComponentName = "checkbox_component"

type checkboxComponent struct {
	*tgframe.BaseComponent
	Label string `json:"label"`
}

func newCheckboxComponent(label string) *checkboxComponent {
	return &checkboxComponent{
		BaseComponent: &tgframe.BaseComponent{
			Name: checkboxComponentName,
			ID:   tcutil.NormalID(checkboxComponentName, label),
		},
		Label: label,
	}
}

// Checkbox create a checkbox and return true if it's clicked.
func Checkbox(s *tgframe.State, c *tgframe.Container, label string) bool {
	comp := newCheckboxComponent(label)
	c.AddComponent(comp)
	return s.GetBool(comp.ID)
}
