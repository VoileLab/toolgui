package component

import "github.com/mudream4869/toolgui/toolgui/framework"

var _ framework.Component = &ButtonComponent{}
var ButtonComponentName = "button_component"

type ButtonComponent struct {
	*framework.BaseComponent
	Label string `json:"label"`
}

func NewButtonComponent(label string) *ButtonComponent {
	return &ButtonComponent{
		BaseComponent: &framework.BaseComponent{
			Name: ButtonComponentName,
			ID:   label,
		},
		Label: label,
	}
}

func Button(sess *framework.Session, c *framework.Container, label string) bool {
	comp := NewButtonComponent(label)
	c.AddComp(comp)
	return sess.GetBool(comp.ID)
}
