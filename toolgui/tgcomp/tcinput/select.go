package tcinput

import (
	"github.com/VoileLab/toolgui/toolgui/tgcomp/tcutil"
	"github.com/VoileLab/toolgui/toolgui/tgframe"
)

var _ tgframe.Component = &selectComponent{}
var selectComponentName = "select_component"

type selectComponent struct {
	*tgframe.BaseComponent
	Label string   `json:"label"`
	Items []string `json:"items"`
}

func newSelectComponent(label string, items []string) *selectComponent {
	return &selectComponent{
		BaseComponent: &tgframe.BaseComponent{
			Name: selectComponentName,
			ID:   tcutil.NormalID(selectComponentName, label),
		},
		Label: label,
		Items: items,
	}
}

// Select create a select dropdown list and return its selected value.
// 1-indexed, return 0 if no item is selected.
func Select(s *tgframe.State, c *tgframe.Container, label string, items []string) int {
	comp := newSelectComponent(label, items)
	c.AddComponent(comp)
	return int(s.GetFloat(comp.ID))
}
