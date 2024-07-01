package tcinput

import (
	"github.com/mudream4869/toolgui/toolgui/component/tcutil"
	"github.com/mudream4869/toolgui/toolgui/framework"
)

var _ framework.Component = &datepickerComponent{}
var datepickerComponentName = "datepicker_component"

type datepickerComponent struct {
	*framework.BaseComponent
	Label string `json:"label"`
}

func newDatepickerComponent(label string) *datepickerComponent {
	return &datepickerComponent{
		BaseComponent: &framework.BaseComponent{
			Name: datepickerComponentName,
			ID:   tcutil.NormalID(datepickerComponentName, label),
		},
		Label: label,
	}
}

func Datepicker(sess *framework.Session, c *framework.Container, label string) string {
	comp := newDatepickerComponent(label)
	c.AddComponent(comp)
	return sess.GetString(comp.ID)
}
