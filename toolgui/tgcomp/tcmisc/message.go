package tcmisc

import (
	"github.com/mudream4869/toolgui/toolgui/framework"
	"github.com/mudream4869/toolgui/toolgui/tgcomp/tcutil"
)

var _ framework.Component = &messageComponent{}
var messageComponentName = "message_component"

type messageComponent struct {
	*framework.BaseComponent
	Type  string `json:"type"`
	Title string `json:"title"`
	Text  string `json:"text"`
}

func newMessageComponent(typ, title, text string) *messageComponent {
	return &messageComponent{
		BaseComponent: &framework.BaseComponent{
			Name: messageComponentName,
			ID:   tcutil.NormalID(messageComponentName, title+text),
		},
		Type:  typ,
		Title: title,
		Text:  text,
	}
}

func Info(c *framework.Container, title, text string) {
	c.AddComponent(newMessageComponent("info", title, text))
}

func Success(c *framework.Container, title, text string) {
	c.AddComponent(newMessageComponent("success", title, text))
}

func Warning(c *framework.Container, title, text string) {
	c.AddComponent(newMessageComponent("warning", title, text))
}

func Error(c *framework.Container, title, text string) {
	c.AddComponent(newMessageComponent("error", title, text))
}

func Danger(c *framework.Container, title, text string) {
	c.AddComponent(newMessageComponent("danger", title, text))
}
