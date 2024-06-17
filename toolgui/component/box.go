package component

import "github.com/mudream4869/toolgui/toolgui/framework"

var _ framework.Component = &TextComponent{}
var BoxComponentName = "box_component"

type BoxComponent struct {
	*framework.BaseComponent

	Container *framework.Container
}

func NewBoxComponent(id string) *BoxComponent {
	return &BoxComponent{
		BaseComponent: &framework.BaseComponent{
			Name: BoxComponentName,
			ID:   id,
		},
	}
}

func Box(c *framework.Container, id string) *framework.Container {
	boxComp := NewBoxComponent(id)
	c.AddComp(boxComp)

	cont := framework.NewContainer(id+"_inner", c.NotifyAddComp)
	boxComp.Container = cont
	c.NotifyAddComp(id, cont)

	return cont
}
