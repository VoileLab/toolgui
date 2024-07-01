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

func Checkbox(sess *framework.Session, c *framework.Container, label string) bool {
	comp := newCheckboxComponent(label)
	c.AddComponent(comp)
	return sess.GetBool(comp.ID)
}
