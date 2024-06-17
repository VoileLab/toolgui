package component

import "github.com/mudream4869/toolgui/toolgui/framework"

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
			ID:   label,
		},
		Label: label,
	}
}

func Checkbox(sess *framework.Session, c *framework.Container, label string) bool {
	comp := NewCheckboxComponent(label)
	c.AddComp(comp)

	v, ok := sess.Values[comp.ID]
	if !ok {
		return false
	}

	return v.(bool)
}
