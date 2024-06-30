package tccontent

import (
	"github.com/mudream4869/toolgui/toolgui/component/tcutil"
	"github.com/mudream4869/toolgui/toolgui/framework"
)

var _ framework.Component = &TextComponent{}
var DividerComponentName = "divider_component"

type DividerComponent struct {
	*framework.BaseComponent
}

func NewDividerComponent() *DividerComponent {
	return &DividerComponent{
		BaseComponent: &framework.BaseComponent{
			Name: DividerComponentName,
			ID:   tcutil.RandID(DividerComponentName),
		},
	}
}

func Divider(c *framework.Container) {
	c.AddComponent(NewDividerComponent())
}
