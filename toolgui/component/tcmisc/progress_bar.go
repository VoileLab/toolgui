package tcmisc

import (
	"github.com/mudream4869/toolgui/toolgui/component/tcutil"
	"github.com/mudream4869/toolgui/toolgui/framework"
)

var _ framework.Component = &ProgressBarComponent{}
var ProgressBarComponentName = "progress_bar_component"

type ProgressBarComponent struct {
	*framework.BaseComponent
	Value int    `json:"value"`
	Label string `json:"label"`

	SendNotifyPack framework.SendNotifyPackFunc `json:"-"`
}

func NewProgressBarComponent(value int, label string, sendNotifyPack framework.SendNotifyPackFunc) *ProgressBarComponent {
	return &ProgressBarComponent{
		BaseComponent: &framework.BaseComponent{
			Name: ProgressBarComponentName,
			ID:   tcutil.HashedID(ProgressBarComponentName, []byte(label)),
		},
		Value: value,
		Label: label,

		SendNotifyPack: sendNotifyPack,
	}
}

func (p *ProgressBarComponent) SetValue(value int) {
	p.Value = value
	p.SendNotifyPack(framework.NewNotifyPackUpdate(p))
}

func (p *ProgressBarComponent) SetLabel(label string) {
	p.Label = label
	p.SendNotifyPack(framework.NewNotifyPackUpdate(p))
}

func (p *ProgressBarComponent) Remove() {
	p.SendNotifyPack(framework.NewNotifyPackDelete(p.ID))
}

func ProgressBar(c *framework.Container, value int, label string) *ProgressBarComponent {
	comp := NewProgressBarComponent(value, label, c.SendNotifyPack)
	c.AddComponent(comp)
	return comp
}
