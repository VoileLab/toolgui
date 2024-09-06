package tcdata

import (
	"github.com/mudream4869/toolgui/toolgui/framework"
	"github.com/mudream4869/toolgui/toolgui/tgcomp/tcutil"
)

var _ framework.Component = &tableComponent{}
var tableComponentName = "table_component"

type tableComponent struct {
	*framework.BaseComponent
	Head  []string   `json:"head"`
	Table [][]string `json:"table"`
}

func newTableComponent(head []string, table [][]string) *tableComponent {
	return &tableComponent{
		BaseComponent: &framework.BaseComponent{
			Name: tableComponentName,
			ID:   tcutil.RandID(tableComponentName),
		},
		Head:  head,
		Table: table,
	}
}

// Table create a table by heading(head) and values(table).
func Table(c *framework.Container, head []string, table [][]string) {
	if len(table) == 0 {
		return
	}

	if len(table[0]) != len(head) {
		panic("len of head should equal to len of table[0]")
	}

	comp := newTableComponent(head, table)
	c.AddComponent(comp)
}
