package tcinput

import (
	"github.com/mudream4869/toolgui/toolgui/tgcomp/tcutil"
	"github.com/mudream4869/toolgui/toolgui/tgframe"
)

var _ tgframe.Component = &textareaComponent{}
var textareaComponentName = "textarea_component"

type textareaComponent struct {
	*tgframe.BaseComponent
	Label  string `json:"label"`
	Height int    `json:"height"`
}

func newTextareaComponent(label string, height int) *textareaComponent {
	return &textareaComponent{
		BaseComponent: &tgframe.BaseComponent{
			Name: textareaComponentName,
			ID:   tcutil.NormalID(textareaComponentName, label),
		},
		Label:  label,
		Height: height,
	}
}

// Textarea create a textarea and return its value.
func Textarea(s *tgframe.State, c *tgframe.Container, label string, height int) string {
	comp := newTextareaComponent(label, height)
	c.AddComponent(comp)
	return s.GetString(comp.ID)
}
