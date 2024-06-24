package component

import (
	"github.com/mudream4869/toolgui/toolgui/framework"
)

var _ framework.Component = &TextComponent{}
var BoxComponentName = "box_component"

type BoxComponent struct {
	*framework.BaseComponent
}

func NewBoxComponent(id string) *BoxComponent {
	return &BoxComponent{
		BaseComponent: &framework.BaseComponent{
			Name: BoxComponentName,
			ID:   normalID(BoxComponentName, id),
		},
	}
}

func Box(c *framework.Container, id string) *framework.Container {
	boxComp := NewBoxComponent(id)
	c.AddComponent(boxComp)

	cont := framework.NewContainer(boxComp.ID+"_inner", c.SendNotifyPack)
	c.SendNotifyPack(framework.NewNotifyPackCreate(boxComp.ID, cont))

	return cont
}
