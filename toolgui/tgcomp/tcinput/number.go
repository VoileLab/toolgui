package tcinput

import (
	"github.com/VoileLab/toolgui/toolgui/tgcomp/tcutil"
	"github.com/VoileLab/toolgui/toolgui/tgframe"
)

var _ tgframe.Component = &numberComponent{}
var numberComponentName = "number_component"

type numberComponent struct {
	*tgframe.BaseComponent

	Label       string       `json:"label"`
	Default     *float64     `json:"default,omitempty"`
	Min         *float64     `json:"min,omitempty"`
	Max         *float64     `json:"max,omitempty"`
	Step        *float64     `json:"step,omitempty"`
	Color       tcutil.Color `json:"color"`
	Placeholder string       `json:"placeholder"`
	Disabled    bool         `json:"disabled"`
}

func newNumberComponent(label string) *numberComponent {
	return &numberComponent{
		BaseComponent: &tgframe.BaseComponent{
			Name: numberComponentName,
			ID:   tcutil.NormalID(numberComponentName, label),
		},
		Label: label,
	}
}

// NumberConf is the configuration for a number component.
type NumberConf struct {
	// Default is the default value of the number component.
	Default *float64

	// Min is the minimum value of the number component.
	Min *float64

	// Max is the maximum value of the number component.
	Max *float64

	// Step is the step of the number component.
	Step *float64

	// Color is the color of the number component.
	Color tcutil.Color

	// Placeholder is the placeholder of the number component.
	Placeholder string

	// Disabled is the disabled state of the number component.
	Disabled bool

	// ID is the ID of the number component.
	ID string
}

func (c *NumberConf) SetMin(min float64) *NumberConf {
	c.Min = &min
	return c
}

func (c *NumberConf) SetMax(max float64) *NumberConf {
	c.Max = &max
	return c
}

func (c *NumberConf) SetStep(step float64) *NumberConf {
	c.Step = &step
	return c
}

func Number(c *tgframe.Container, s *tgframe.State, label string) *float64 {
	return NumberWithConf(c, s, label, nil)
}

func NumberWithConf(c *tgframe.Container, s *tgframe.State, label string, conf *NumberConf) *float64 {
	if conf == nil {
		conf = &NumberConf{}
	}

	comp := newNumberComponent(label)
	comp.Placeholder = conf.Placeholder
	comp.Color = conf.Color
	comp.Default = conf.Default
	comp.Min = conf.Min
	comp.Max = conf.Max
	comp.Step = conf.Step
	comp.Disabled = conf.Disabled

	if conf.ID != "" {
		comp.SetID(conf.ID)
	}

	c.AddComponent(comp)

	val := s.GetFloat(comp.ID)
	if val == nil {
		return conf.Default
	}

	return val
}
