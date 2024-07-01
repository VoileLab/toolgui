package tccontent

import (
	"github.com/mudream4869/toolgui/toolgui/component/tcutil"
	"github.com/mudream4869/toolgui/toolgui/framework"
)

var _ framework.Component = &codeComponent{}
var codeComponentName = "code_component"

type codeComponent struct {
	*framework.BaseComponent
	Code string `json:"code"`
	Lang string `json:"lang"`
}

func newCodeComponent(code, lang string) *codeComponent {
	return &codeComponent{
		BaseComponent: &framework.BaseComponent{
			Name: codeComponentName,
			ID:   tcutil.HashedID(codeComponentName, []byte(code)),
		},
		Code: code,
		Lang: lang,
	}
}

func Code(c *framework.Container, text, lang string) {
	comp := newCodeComponent(text, lang)
	c.AddComponent(comp)
}
