package tcinput

import (
	"github.com/mudream4869/toolgui/toolgui/component/tcutil"
	"github.com/mudream4869/toolgui/toolgui/framework"
)

var _ framework.Component = &textareaComponent{}
var textareaComponentName = "textarea_component"

type textareaComponent struct {
	*framework.BaseComponent
	Label  string `json:"label"`
	Height int    `json:"height"`
}

func newTextareaComponent(label string, height int) *textareaComponent {
	return &textareaComponent{
		BaseComponent: &framework.BaseComponent{
			Name: textareaComponentName,
			ID:   tcutil.NormalID(textareaComponentName, label),
		},
		Label:  label,
		Height: height,
	}
}

// Textarea create a textarea and return its value.
func Textarea(s *framework.State, c *framework.Container, label string, height int) string {
	comp := newTextareaComponent(label, height)
	c.AddComponent(comp)
	return s.GetString(comp.ID)
}
