package tcinput

import (
	"github.com/mudream4869/toolgui/toolgui/framework"
	"github.com/mudream4869/toolgui/toolgui/tgcomp/tcutil"
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

// Datepicker create a datepicker and return its selected date.
func Datepicker(s *framework.State, c *framework.Container, label string) string {
	comp := newDatepickerComponent(label, "date")
	c.AddComponent(comp)
	return s.GetString(comp.ID)
}

// Timepicker create a timepicker and return its selected time.
func Timepicker(s *framework.State, c *framework.Container, label string) string {
	comp := newDatepickerComponent(label, "time")
	c.AddComponent(comp)
	return s.GetString(comp.ID)
}

// Datetimepicker create a datetimepicker and return its selected datetime.
func Datetimepicker(s *framework.State, c *framework.Container, label string) string {
	comp := newDatepickerComponent(label, "datetime-local")
	c.AddComponent(comp)
	return s.GetString(comp.ID)
}
