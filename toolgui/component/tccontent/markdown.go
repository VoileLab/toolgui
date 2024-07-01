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

func Markdown(c *framework.Container, text string) {
	comp := newMarkdownComponent(text)
	c.AddComponent(comp)
}
