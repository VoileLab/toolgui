package tcdata

import (
	"encoding/json"

	"github.com/mudream4869/toolgui/toolgui/component/tcutil"
	"github.com/mudream4869/toolgui/toolgui/framework"
)

var _ framework.Component = &JSONComponent{}
var JSONComponentName = "json_component"

type JSONComponent struct {
	*framework.BaseComponent
	Value string `json:"value"`
}

func NewJSONComponent(s string) *JSONComponent {
	return &JSONComponent{
		BaseComponent: &framework.BaseComponent{
			Name: JSONComponentName,
			ID:   tcutil.HashedID(JSONComponentName, []byte(s)),
		},
		Value: s,
	}
}

func JSON(c *framework.Container, v any) {
	bs, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}

	comp := NewJSONComponent(string(bs))
	c.AddComponent(comp)
}
