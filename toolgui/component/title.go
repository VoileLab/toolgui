package component

import "github.com/mudream4869/toolgui/toolgui/framework"

var _ framework.Component = &TitleComponent{}
var TitleComponentName = "title_component"

type TitleComponent struct {
	*framework.BaseComponent
	Text string `json:"text"`
}

func NewTitleComponent(text string) *TitleComponent {
	return &TitleComponent{
		BaseComponent: &framework.BaseComponent{
			Name: TitleComponentName,
			ID:   text,
		},
		Text: text,
	}
}

func Title(c *framework.Container, text string) {
	c.AddComponent(NewTitleComponent(text))
}
