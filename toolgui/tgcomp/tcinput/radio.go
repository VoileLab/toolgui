package tcinput

import (
	"github.com/VoileLab/toolgui/toolgui/tgcomp/tcutil"
	"github.com/VoileLab/toolgui/toolgui/tgframe"
)

var _ tgframe.Component = &radioComponent{}
var radioComponentName = "radio_component"

type radioComponent struct {
	*tgframe.BaseComponent
	Label string   `json:"label"`
	Items []string `json:"items"`
}

func newRadioComponent(label string, items []string) *radioComponent {
	return &radioComponent{
		BaseComponent: &tgframe.BaseComponent{
			Name: radioComponentName,
			ID:   tcutil.NormalID(radioComponentName, label),
		},
		Label: label,
		Items: items,
	}
}

// Radio create a group of radio items and return its selected value.
func Radio(s *tgframe.State, c *tgframe.Container, label string, items []string) string {
	comp := newRadioComponent(label, items)
	c.AddComponent(comp)
	return s.GetString(comp.ID, "")
}
