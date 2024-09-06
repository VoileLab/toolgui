package tclayout

import (
	"github.com/mudream4869/toolgui/toolgui/framework"
	"github.com/mudream4869/toolgui/toolgui/tgcomp/tcutil"
)

var _ framework.Component = &boxComponent{}
var boxComponentName = "box_component"

type boxComponent struct {
	*framework.BaseComponent
}

func newBoxComponent(id string) *boxComponent {
	return &boxComponent{
		BaseComponent: &framework.BaseComponent{
			Name: boxComponentName,
			ID:   tcutil.NormalID(boxComponentName, id),
		},
	}
}

// Box create a box container.
func Box(c *framework.Container, id string) *framework.Container {
	boxComp := newBoxComponent(id)
	c.AddComponent(boxComp)

	cont := framework.NewContainer(boxComp.ID+"_inner", c.SendNotifyPack)
	c.SendNotifyPack(framework.NewNotifyPackCreate(boxComp.ID, cont))

	return cont
}
