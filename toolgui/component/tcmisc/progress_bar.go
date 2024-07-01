package tcmisc

import (
	"github.com/mudream4869/toolgui/toolgui/component/tcutil"
	"github.com/mudream4869/toolgui/toolgui/framework"
)

var _ framework.Component = &progressBarComponent{}
var progressBarComponentName = "progress_bar_component"

type progressBarComponent struct {
	*framework.BaseComponent
	Value int    `json:"value"`
	Label string `json:"label"`

	SendNotifyPack framework.SendNotifyPackFunc `json:"-"`
}

func newProgressBarComponent(value int, label string, sendNotifyPack framework.SendNotifyPackFunc) *progressBarComponent {
	return &progressBarComponent{
		BaseComponent: &framework.BaseComponent{
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
	p.SendNotifyPack(framework.NewNotifyPackUpdate(p))
}

func (p *progressBarComponent) SetLabel(label string) {
	p.Label = label
	p.SendNotifyPack(framework.NewNotifyPackUpdate(p))
}

func (p *progressBarComponent) Remove() {
	p.SendNotifyPack(framework.NewNotifyPackDelete(p.ID))
}

func ProgressBar(c *framework.Container, value int, label string) *progressBarComponent {
	comp := newProgressBarComponent(value, label, c.SendNotifyPack)
	c.AddComponent(comp)
	return comp
}
