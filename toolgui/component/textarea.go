package component

import "github.com/mudream4869/toolgui/toolgui/framework"

var _ framework.Component = &TextareaComponent{}
var TextareaComponentName = "textarea_component"

type TextareaComponent struct {
	*framework.BaseComponent
	Label string `json:"label"`
}

func NewTextareaComponent(label string) *TextareaComponent {
	return &TextareaComponent{
		BaseComponent: &framework.BaseComponent{
			Name: TextareaComponentName,
			ID:   normalID(TextareaComponentName, label),
		},
		Label: label,
	}
}

func Textarea(sess *framework.Session, c *framework.Container, label string) string {
	comp := NewTextareaComponent(label)
	c.AddComponent(comp)
	return sess.GetString(comp.ID)
}
