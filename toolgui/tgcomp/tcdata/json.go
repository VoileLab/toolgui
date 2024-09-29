package tcdata

import (
	"encoding/json"

	"github.com/VoileLab/toolgui/toolgui/tgcomp/tcutil"
	"github.com/VoileLab/toolgui/toolgui/tgframe"
)

var _ tgframe.Component = &jsonComponent{}
var jsonComponentName = "json_component"

type jsonComponent struct {
	*tgframe.BaseComponent
	Value string `json:"value"`
}

func newJSONComponent(s string) *jsonComponent {
	return &jsonComponent{
		BaseComponent: &tgframe.BaseComponent{
			Name: jsonComponentName,
			ID:   tcutil.HashedID(jsonComponentName, []byte(s)),
		},
		Value: s,
	}
}

// JSON create a JSON viewer for v.
func JSON(c *tgframe.Container, v any) {
	bs, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}

	comp := newJSONComponent(string(bs))
	c.AddComponent(comp)
}
