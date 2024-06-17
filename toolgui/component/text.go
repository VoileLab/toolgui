package component

import "github.com/mudream4869/toolgui/toolgui/framework"

var _ framework.Component = &TextComponent{}
var TextComponentName = "text_component"

type TextComponent struct {
	*framework.BaseComponent
	Text string `json:"text"`
}

func NewTextComponent(text string) *TextComponent {
	return &TextComponent{
		BaseComponent: &framework.BaseComponent{
			Name: TextComponentName,
			ID:   text,
		},
		Text: text,
	}
}

func Text(c *framework.Container, text string) {
	comp := NewTextComponent(text)
	c.AddComp(comp)
}
