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
	Type  string `json:"type"`
}

func newDatepickerComponent(label string, typ string) *datepickerComponent {
	return &datepickerComponent{
		BaseComponent: &framework.BaseComponent{
			Name: datepickerComponentName,
			ID:   tcutil.NormalID(datepickerComponentName, label),
		},
		Label: label,
		Type:  typ,
	}
}

func Datepicker(sess *framework.Session, c *framework.Container, label string) string {
	comp := newDatepickerComponent(label, "date")
	c.AddComponent(comp)
	return sess.GetString(comp.ID)
}

func Timepicker(sess *framework.Session, c *framework.Container, label string) string {
	comp := newDatepickerComponent(label, "time")
	c.AddComponent(comp)
	return sess.GetString(comp.ID)
}

func Datetimepicker(sess *framework.Session, c *framework.Container, label string) string {
	comp := newDatepickerComponent(label, "datetime-local")
	c.AddComponent(comp)
	return sess.GetString(comp.ID)
}
