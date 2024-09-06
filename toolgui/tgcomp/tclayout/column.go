package tclayout

import (
	"fmt"

	"github.com/mudream4869/toolgui/toolgui/tgcomp/tcutil"
	"github.com/mudream4869/toolgui/toolgui/tgframe"
)

var _ tgframe.Component = &columnComponent{}
var columnComponentName = "column_component"

type columnComponent struct {
	*tgframe.BaseComponent
}

func newColumnComponent(id string) *columnComponent {
	return &columnComponent{
		BaseComponent: &tgframe.BaseComponent{
			Name: columnComponentName,
			ID:   tcutil.NormalID(columnComponentName, id),
		},
	}
}

// Column create N columns.
func Column(c *tgframe.Container, id string, n uint) []*tgframe.Container {
	if n == 0 {
		panic("number of columns should > 0")
	}

	colsComp := newColumnComponent(id)
	c.AddComponent(colsComp)

	cols := make([]*tgframe.Container, n)
	for i := range n {
		cols[i] = tgframe.NewContainer(fmt.Sprintf("%s_%d", colsComp.ID, i), c.SendNotifyPack)
		c.SendNotifyPack(tgframe.NewNotifyPackCreate(colsComp.ID, cols[i]))
	}

	return cols
}

// Column2 create 2 columns.
func Column2(c *tgframe.Container, id string) (*tgframe.Container, *tgframe.Container) {
	cols := Column(c, id, 2)
	return cols[0], cols[1]
}

// Column3 create 3 columns.
func Column3(c *tgframe.Container, id string) (*tgframe.Container, *tgframe.Container, *tgframe.Container) {
	cols := Column(c, id, 3)
	return cols[0], cols[1], cols[2]
}
