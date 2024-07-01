package tcinput

import (
	"github.com/mudream4869/toolgui/toolgui/component/tcutil"
	"github.com/mudream4869/toolgui/toolgui/framework"
)

var _ framework.Component = &textboxComponent{}
var textboxComponentName = "textbox_component"

type textboxComponent struct {
	*framework.BaseComponent
	Label string `json:"label"`
}

func newTextboxComponent(label string) *textboxComponent {
	return &textboxComponent{
		BaseComponent: &framework.BaseComponent{
			Name: textboxComponentName,
			ID:   tcutil.NormalID(textboxComponentName, label),
		},
		Label: label,
	}
}

func Textbox(sess *framework.Session, c *framework.Container, label string) string {
	comp := newTextboxComponent(label)
	c.AddComponent(comp)
	return sess.GetString(comp.ID)
}
