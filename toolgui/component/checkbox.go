package component

import (
	"github.com/mudream4869/toolgui/toolgui/framework"
)

var _ framework.Component = &CheckboxComponent{}
var CheckboxComponentName = "checkbox_component"

type CheckboxComponent struct {
	*framework.BaseComponent
	Label string `json:"label"`
}

func NewCheckboxComponent(label string) *CheckboxComponent {
	return &CheckboxComponent{
		BaseComponent: &framework.BaseComponent{
			Name: CheckboxComponentName,
			ID:   normalID(CheckboxComponentName, label),
		},
		Label: label,
	}
}

func Checkbox(sess *framework.Session, c *framework.Container, label string) bool {
	comp := NewCheckboxComponent(label)
	c.AddComponent(comp)
	return sess.GetBool(comp.ID)
}
