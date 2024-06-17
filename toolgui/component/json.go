package component

import (
	"crypto/md5"
	"encoding/json"
	"fmt"

	"github.com/mudream4869/toolgui/toolgui/framework"
)

var _ framework.Component = &JSONComponent{}
var JSONComponentName = "json_component"

type JSONComponent struct {
	*framework.BaseComponent
	Value string `json:"value"`
}

func NewJSONComponent(s string) *JSONComponent {
	id := fmt.Sprintf("json_%x", md5.Sum([]byte(s)))
	return &JSONComponent{
		BaseComponent: &framework.BaseComponent{
			Name: JSONComponentName,
			ID:   id,
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
	c.AddComp(comp)
}
