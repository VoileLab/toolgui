package component

import "github.com/mudream4869/toolgui/toolgui/framework"

var _ framework.Component = &CodeComponent{}
var CodeComponentName = "code_component"

type CodeComponent struct {
	*framework.BaseComponent
	Code string `json:"code"`
	Lang string `json:"lang"`
}

func NewCodeComponent(code, lang string) *CodeComponent {
	return &CodeComponent{
		BaseComponent: &framework.BaseComponent{
			Name: CodeComponentName,
			ID:   code,
		},
		Code: code,
		Lang: lang,
	}
}

func Code(c *framework.Container, text, lang string) {
	comp := NewCodeComponent(text, lang)
	c.AddComponent(comp)
}
