package tcinput

import (
	"github.com/mudream4869/toolgui/toolgui/component/tcutil"
	"github.com/mudream4869/toolgui/toolgui/framework"
)

var _ framework.Component = &textareaComponent{}
var textareaComponentName = "textarea_component"

type textareaComponent struct {
	*framework.BaseComponent
	Label string `json:"label"`
}

func newTextareaComponent(label string) *textareaComponent {
	return &textareaComponent{
		BaseComponent: &framework.BaseComponent{
			Name: textareaComponentName,
			ID:   tcutil.NormalID(textareaComponentName, label),
		},
		Label: label,
	}
}

func Textarea(s *framework.State, c *framework.Container, label string) string {
	comp := newTextareaComponent(label)
	c.AddComponent(comp)
	return s.GetString(comp.ID)
}
