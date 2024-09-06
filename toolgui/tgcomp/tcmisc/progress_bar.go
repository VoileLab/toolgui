package tcmisc

import (
	"github.com/mudream4869/toolgui/toolgui/tgcomp/tcutil"
	"github.com/mudream4869/toolgui/toolgui/tgframe"
)

var _ tgframe.Component = &progressBarComponent{}
var progressBarComponentName = "progress_bar_component"

type progressBarComponent struct {
	*tgframe.BaseComponent
	Value int    `json:"value"`
	Label string `json:"label"`

	SendNotifyPack tgframe.SendNotifyPackFunc `json:"-"`
}

func newProgressBarComponent(value int, label string, sendNotifyPack tgframe.SendNotifyPackFunc) *progressBarComponent {
	return &progressBarComponent{
		BaseComponent: &tgframe.BaseComponent{
			Name: progressBarComponentName,
			ID:   tcutil.HashedID(progressBarComponentName, []byte(label)),
		},
		Value: value,
		Label: label,

		SendNotifyPack: sendNotifyPack,
	}
}

func (p *progressBarComponent) SetValue(value int) {
	p.Value = value
	p.SendNotifyPack(tgframe.NewNotifyPackUpdate(p))
}

func (p *progressBarComponent) SetLabel(label string) {
	p.Label = label
	p.SendNotifyPack(tgframe.NewNotifyPackUpdate(p))
}

func (p *progressBarComponent) Remove() {
	p.SendNotifyPack(tgframe.NewNotifyPackDelete(p.ID))
}

func ProgressBar(c *tgframe.Container, value int, label string) *progressBarComponent {
	comp := newProgressBarComponent(value, label, c.SendNotifyPack)
	c.AddComponent(comp)
	return comp
}
