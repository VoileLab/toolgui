package component

import "github.com/mudream4869/toolgui/toolgui/framework"

var _ framework.Component = &SelectComponent{}
var SelectComponentName = "select_component"

type SelectComponent struct {
	*framework.BaseComponent
	Label string   `json:"label"`
	Items []string `json:"items"`
}

func NewSelectComponent(label string, items []string) *SelectComponent {
	return &SelectComponent{
		BaseComponent: &framework.BaseComponent{
			Name: SelectComponentName,
			ID:   label,
		},
		Label: label,
		Items: items,
	}
}

func Select(sess *framework.Session, c *framework.Container, label string, items []string) string {
	comp := NewSelectComponent(label, items)
	c.AddComponent(comp)
	return sess.GetString(comp.ID)
}
