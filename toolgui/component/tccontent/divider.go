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

func Divider(c *framework.Container) {
	c.AddComponent(newDividerComponent())
}
