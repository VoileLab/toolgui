package tccontent

import (
	"github.com/mudream4869/toolgui/toolgui/component/tcutil"
	"github.com/mudream4869/toolgui/toolgui/framework"
)

var _ framework.Component = &markdownComponent{}
var markdownComponentName = "markdown_component"

type markdownComponent struct {
	*framework.BaseComponent
	Markdown string `json:"text"`
}

func newMarkdownComponent(text string) *markdownComponent {
	return &markdownComponent{
		BaseComponent: &framework.BaseComponent{
			Name: markdownComponentName,
			ID:   tcutil.HashedID(markdownComponentName, []byte(text)),
		},
		Markdown: text,
	}
}

// Markdown render markdown to html
func Markdown(c *framework.Container, markdown string) {
	comp := newMarkdownComponent(markdown)
	c.AddComponent(comp)
}

// Markdown create a markdown rendering space with a user-specific id
func MarkdownWithID(c *framework.Container, markdown string, id string) {
	comp := newMarkdownComponent(markdown)
	comp.SetID(id)
	c.AddComponent(comp)
}
