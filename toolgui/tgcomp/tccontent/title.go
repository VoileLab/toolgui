package tccontent

import (
	"github.com/mudream4869/toolgui/toolgui/framework"
	"github.com/mudream4869/toolgui/toolgui/tgcomp/tcutil"
)

var _ framework.Component = &titleComponent{}
var titleComponentName = "title_component"

type titleComponent struct {
	*framework.BaseComponent
	Text string `json:"text"`
}

func newTitleComponent(text string) *titleComponent {
	return &titleComponent{
		BaseComponent: &framework.BaseComponent{
			Name: titleComponentName,
			ID:   tcutil.NormalID(titleComponentName, text),
		},
		Text: text,
	}
}

// Title show a title.
func Title(c *framework.Container, text string) {
	c.AddComponent(newTitleComponent(text))
}

// TitleWithID create a title component with a user specific id.
func TitleWithID(c *framework.Container, text string, id string) {
	comp := newTitleComponent(text)
	comp.SetID(id)
	c.AddComponent(comp)
}
