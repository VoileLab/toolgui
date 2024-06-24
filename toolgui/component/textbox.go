package component

import "github.com/mudream4869/toolgui/toolgui/framework"

var _ framework.Component = &TextboxComponent{}
var TextboxComponentName = "textbox_component"

type TextboxComponent struct {
	*framework.BaseComponent
	Label string `json:"label"`
}

func NewTextboxComponent(label string) *TextboxComponent {
	return &TextboxComponent{
		BaseComponent: &framework.BaseComponent{
			Name: TextboxComponentName,
			ID:   normalID(TextboxComponentName, label),
		},
		Label: label,
	}
}

func Textbox(sess *framework.Session, c *framework.Container, label string) string {
	comp := NewTextboxComponent(label)
	c.AddComponent(comp)
	return sess.GetString(comp.ID)
}
