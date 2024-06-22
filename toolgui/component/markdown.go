package component

import "github.com/mudream4869/toolgui/toolgui/framework"

var _ framework.Component = &MarkdownComponent{}
var MarkdownComponentName = "markdown_component"

type MarkdownComponent struct {
	*framework.BaseComponent
	Markdown string `json:"text"`
}

func NewMarkdownComponent(text string) *MarkdownComponent {
	return &MarkdownComponent{
		BaseComponent: &framework.BaseComponent{
			Name: MarkdownComponentName,
			ID:   text,
		},
		Markdown: text,
	}
}

func Markdown(c *framework.Container, text string) {
	comp := NewMarkdownComponent(text)
	c.AddComponent(comp)
}
