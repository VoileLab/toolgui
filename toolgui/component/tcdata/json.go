package tcdata

import (
	"encoding/json"

	"github.com/mudream4869/toolgui/toolgui/component/tcutil"
	"github.com/mudream4869/toolgui/toolgui/framework"
)

var _ framework.Component = &jsonComponent{}
var jsonComponentName = "json_component"

type jsonComponent struct {
	*framework.BaseComponent
	Value string `json:"value"`
}

func newJSONComponent(s string) *jsonComponent {
	return &jsonComponent{
		BaseComponent: &framework.BaseComponent{
			Name: jsonComponentName,
			ID:   tcutil.HashedID(jsonComponentName, []byte(s)),
		},
		Value: s,
	}
}

// JSON create a JSON viewer for v.
func JSON(c *framework.Container, v any) {
	bs, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}

	comp := newJSONComponent(string(bs))
	c.AddComponent(comp)
}
