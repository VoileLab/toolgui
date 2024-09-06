package tclayout

import (
	"fmt"

	"github.com/mudream4869/toolgui/toolgui/framework"
	"github.com/mudream4869/toolgui/toolgui/tgcomp/tcutil"
)

var _ framework.Component = &columnComponent{}
var columnComponentName = "column_component"

type columnComponent struct {
	*framework.BaseComponent
}

func newColumnComponent(id string) *columnComponent {
	return &columnComponent{
		BaseComponent: &framework.BaseComponent{
			Name: columnComponentName,
			ID:   tcutil.NormalID(columnComponentName, id),
		},
	}
}

// Column create N columns.
func Column(c *framework.Container, id string, n uint) []*framework.Container {
	if n == 0 {
		panic("number of columns should > 0")
	}

	colsComp := newColumnComponent(id)
	c.AddComponent(colsComp)

	cols := make([]*framework.Container, n)
	for i := range n {
		cols[i] = framework.NewContainer(fmt.Sprintf("%s_%d", colsComp.ID, i), c.SendNotifyPack)
		c.SendNotifyPack(framework.NewNotifyPackCreate(colsComp.ID, cols[i]))
	}

	return cols
}

// Column2 create 2 columns.
func Column2(c *framework.Container, id string) (*framework.Container, *framework.Container) {
	cols := Column(c, id, 2)
	return cols[0], cols[1]
}

// Column3 create 3 columns.
func Column3(c *framework.Container, id string) (*framework.Container, *framework.Container, *framework.Container) {
	cols := Column(c, id, 3)
	return cols[0], cols[1], cols[2]
}
