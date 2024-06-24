package component

import "github.com/mudream4869/toolgui/toolgui/framework"

var _ framework.Component = &MessageComponent{}
var MessageComponentName = "message_component"

type MessageComponent struct {
	*framework.BaseComponent
	Type  string `json:"type"`
	Title string `json:"title"`
	Text  string `json:"text"`
}

func NewMessageComponent(typ, title, text string) *MessageComponent {
	return &MessageComponent{
		BaseComponent: &framework.BaseComponent{
			Name: MessageComponentName,
			ID:   normalID(MessageComponentName, title+text),
		},
		Type:  typ,
		Title: title,
		Text:  text,
	}
}

func Info(c *framework.Container, title, text string) {
	c.AddComponent(NewMessageComponent("info", title, text))
}

func Success(c *framework.Container, title, text string) {
	c.AddComponent(NewMessageComponent("success", title, text))
}

func Warning(c *framework.Container, title, text string) {
	c.AddComponent(NewMessageComponent("warning", title, text))
}

func Error(c *framework.Container, title, text string) {
	c.AddComponent(NewMessageComponent("error", title, text))
}

func Danger(c *framework.Container, title, text string) {
	c.AddComponent(NewMessageComponent("danger", title, text))
}
