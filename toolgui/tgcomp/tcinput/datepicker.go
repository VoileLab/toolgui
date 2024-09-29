package tcinput

import (
	"github.com/VoileLab/toolgui/toolgui/tgcomp/tcutil"
	"github.com/VoileLab/toolgui/toolgui/tgframe"
)

var _ tgframe.Component = &datepickerComponent{}
var datepickerComponentName = "datepicker_component"

type datepickerComponent struct {
	*tgframe.BaseComponent
	Label string `json:"label"`
	Type  string `json:"type"`
}

func newDatepickerComponent(label string, typ string) *datepickerComponent {
	return &datepickerComponent{
		BaseComponent: &tgframe.BaseComponent{
			Name: datepickerComponentName,
			ID:   tcutil.NormalID(datepickerComponentName, label),
		},
		Label: label,
		Type:  typ,
	}
}

// Datepicker create a datepicker and return its selected date.
func Datepicker(s *tgframe.State, c *tgframe.Container, label string) string {
	comp := newDatepickerComponent(label, "date")
	c.AddComponent(comp)
	return s.GetString(comp.ID)
}

// Timepicker create a timepicker and return its selected time.
func Timepicker(s *tgframe.State, c *tgframe.Container, label string) string {
	comp := newDatepickerComponent(label, "time")
	c.AddComponent(comp)
	return s.GetString(comp.ID)
}

// Datetimepicker create a datetimepicker and return its selected datetime.
func Datetimepicker(s *tgframe.State, c *tgframe.Container, label string) string {
	comp := newDatepickerComponent(label, "datetime-local")
	c.AddComponent(comp)
	return s.GetString(comp.ID)
}
