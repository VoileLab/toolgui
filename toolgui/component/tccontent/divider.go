package tccontent

import (
	"github.com/mudream4869/toolgui/toolgui/component/tcutil"
	"github.com/mudream4869/toolgui/toolgui/framework"
)

var _ framework.Component = &textComponent{}
var dividerComponentName = "divider_component"

type dividerComponent struct {
	*framework.BaseComponent
}

func newDividerComponent() *dividerComponent {
	return &dividerComponent{
		BaseComponent: &framework.BaseComponent{
			Name: dividerComponentName,
			ID:   tcutil.RandID(dividerComponentName),
		},
	}
}

// Divider create a horizontal line.
func Divider(c *framework.Container) {
	c.AddComponent(newDividerComponent())
}

// DividerWithID create a horizontal line with ID.
func DividerWithID(c *framework.Container, id string) {
	comp := newDividerComponent()
	comp.SetID(id)
	c.AddComponent(comp)
}
