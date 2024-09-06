package tcmisc

import (
	"github.com/mudream4869/toolgui/toolgui/tgcomp/tcutil"
	"github.com/mudream4869/toolgui/toolgui/tgframe"
)

var _ tgframe.Component = &messageComponent{}
var messageComponentName = "message_component"

type messageComponent struct {
	*tgframe.BaseComponent
	Type  string `json:"type"`
	Title string `json:"title"`
	Text  string `json:"text"`
}

func newMessageComponent(typ, title, text string) *messageComponent {
	return &messageComponent{
		BaseComponent: &tgframe.BaseComponent{
			Name: messageComponentName,
			ID:   tcutil.NormalID(messageComponentName, title+text),
		},
		Type:  typ,
		Title: title,
		Text:  text,
	}
}

func Info(c *tgframe.Container, title, text string) {
	c.AddComponent(newMessageComponent("info", title, text))
}

func Success(c *tgframe.Container, title, text string) {
	c.AddComponent(newMessageComponent("success", title, text))
}

func Warning(c *tgframe.Container, title, text string) {
	c.AddComponent(newMessageComponent("warning", title, text))
}

func Error(c *tgframe.Container, title, text string) {
	c.AddComponent(newMessageComponent("error", title, text))
}

func Danger(c *tgframe.Container, title, text string) {
	c.AddComponent(newMessageComponent("danger", title, text))
}
