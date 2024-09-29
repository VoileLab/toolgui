package tcinput

import (
	"github.com/VoileLab/toolgui/toolgui/tgcomp/tcutil"
	"github.com/VoileLab/toolgui/toolgui/tgframe"
)

var _ tgframe.Component = &buttonComponent{}
var buttonComponentName = "button_component"

type buttonComponent struct {
	*tgframe.BaseComponent
	Label    string       `json:"label"`
	Color    tcutil.Color `json:"color"`
	Disabled bool         `json:"disabled"`
}

func newButtonComponent(label string) *buttonComponent {
	return &buttonComponent{
		BaseComponent: &tgframe.BaseComponent{
			Name: buttonComponentName,
			ID:   tcutil.NormalID(buttonComponentName, label),
		},
		Label: label,
	}
}

type ButtonConf struct {
	Color    tcutil.Color
	Disabled bool

	ID string
}

// Button create a button and return true if it's clicked.
func Button(s *tgframe.State, c *tgframe.Container, label string) bool {
	return ButtonWithConf(s, c, label, nil)
}

// ButtonWithConf create a button and return true if it's clicked.
func ButtonWithConf(s *tgframe.State, c *tgframe.Container, label string, conf *ButtonConf) bool {
	if conf == nil {
		conf = &ButtonConf{}
	}

	comp := newButtonComponent(label)
	comp.Color = conf.Color
	comp.Disabled = conf.Disabled

	if comp.ID != "" {
		comp.SetID(comp.ID)
	}

	c.AddComponent(comp)
	return s.GetClickID() == comp.ID
}
