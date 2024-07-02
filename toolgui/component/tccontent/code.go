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

// CodeConf provide extra config for Code Component
type CodeConf struct {
	ID string
}

// Code create a code block with syntax highlight
func Code(c *framework.Container, code, lang string) {
	CodeWithConf(c, code, lang, nil)
}

// CodeWithConf create a code block with syntax highlight
func CodeWithConf(c *framework.Container, code, lang string, conf *CodeConf) {
	comp := newCodeComponent(code, lang)
	if conf != nil {
		if conf.ID != "" {
			comp.SetID(conf.ID)
		}
	}
	c.AddComponent(comp)
}
