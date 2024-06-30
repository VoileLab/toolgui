package tcdata

import (
	"github.com/mudream4869/toolgui/toolgui/component/tcutil"
	"github.com/mudream4869/toolgui/toolgui/framework"
)

var _ framework.Component = &TableComponent{}
var TableComponentName = "table_component"

type TableComponent struct {
	*framework.BaseComponent
	Head  []string   `json:"head"`
	Table [][]string `json:"table"`
}

func NewTableComponent(head []string, table [][]string) *TableComponent {
	return &TableComponent{
		BaseComponent: &framework.BaseComponent{
			Name: TableComponentName,
			ID:   tcutil.RandID(TableComponentName),
		},
		Head:  head,
		Table: table,
	}
}

func Table(c *framework.Container, head []string, table [][]string) {
	if len(table) == 0 {
		return
	}

	if len(table[0]) != len(head) {
		panic("len of head should equal to len of table[0]")
	}

	comp := NewTableComponent(head, table)
	c.AddComponent(comp)
}
