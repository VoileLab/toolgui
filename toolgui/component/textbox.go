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
			ID:   label,
		},
		Label: label,
	}
}

func Textbox(sess *framework.Session, c *framework.Container, label string) string {
	comp := NewTextboxComponent(label)
	c.AddComp(comp)

	v, ok := sess.Values[comp.ID]
	if !ok {
		return ""
	}

	return v.(string)
}
