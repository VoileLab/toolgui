package component

import (
	"fmt"

	"github.com/mudream4869/toolgui/toolgui/framework"
)

var _ framework.Component = &TextComponent{}
var ColumnComponentName = "column_component"

type ColumnComponent struct {
	*framework.BaseComponent
}

func NewColumnComponent(id string) *ColumnComponent {
	return &ColumnComponent{
		BaseComponent: &framework.BaseComponent{
			Name: ColumnComponentName,
			ID:   id,
		},
	}
}

func Column(c *framework.Container, id string, n uint) []*framework.Container {
	if n == 0 {
		panic("number of columns should > 0")
	}

	colsComp := NewColumnComponent(id)
	c.AddComponent(colsComp)

	cols := make([]*framework.Container, n)
	for i := range n {
		cols[i] = framework.NewContainer(fmt.Sprintf("%s_%d", id, i), c.SendNotifyPack)
		c.SendNotifyPack(framework.NewNotifyPackCreate(id, cols[i]))
	}

	return cols
}

func Column2(c *framework.Container, id string) (*framework.Container, *framework.Container) {
	cols := Column(c, id, 2)
	return cols[0], cols[1]
}

func Column3(c *framework.Container, id string) (*framework.Container, *framework.Container, *framework.Container) {
	cols := Column(c, id, 3)
	return cols[0], cols[1], cols[2]
}
