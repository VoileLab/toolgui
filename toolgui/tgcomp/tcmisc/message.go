package tcmisc

import (
	"github.com/VoileLab/toolgui/toolgui/tgcomp/tcutil"
	"github.com/VoileLab/toolgui/toolgui/tgframe"
)

var _ tgframe.Component = &messageComponent{}
var messageComponentName = "message_component"

type messageComponent struct {
	*tgframe.BaseComponent
	Title string       `json:"title"`
	Body  string       `json:"body"`
	Color tcutil.Color `json:"color"`
}

func newMessageComponent(body string) *messageComponent {
	return &messageComponent{
		BaseComponent: &tgframe.BaseComponent{
			Name: messageComponentName,
			ID:   tcutil.NormalID(messageComponentName, body),
		},
		Body: body,
	}
}

type MessageConf struct {
	Title string
	Color tcutil.Color

	ID string
}

func Message(c *tgframe.Container, text string) {
	MessageWithConf(c, text, nil)
}

func MessageWithConf(c *tgframe.Container, text string, conf *MessageConf) {
	if conf == nil {
		conf = &MessageConf{}
	}

	comp := newMessageComponent(text)
	comp.Color = conf.Color
	comp.Title = conf.Title

	if conf.ID != "" {
		comp.SetID(conf.ID)
	}

	c.AddComponent(comp)
}
