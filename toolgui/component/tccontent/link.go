package tccontent

import (
	"github.com/mudream4869/toolgui/toolgui/component/tcutil"
	"github.com/mudream4869/toolgui/toolgui/framework"
)

var _ framework.Component = &linkComponent{}
var linkComponentName = "link_component"

type linkComponent struct {
	*framework.BaseComponent
	Text string `json:"text"`
	URL  string `json:"url"`
}

func newLinkComponent(text, url string) *linkComponent {
	return &linkComponent{
		BaseComponent: &framework.BaseComponent{
			Name: linkComponentName,
			ID:   tcutil.NormalID(linkComponentName, text),
		},
		Text: text,
		URL:  url,
	}
}

// Link create a link component.
func Link(c *framework.Container, text, url string) {
	c.AddComponent(newLinkComponent(text, url))
}

// LinkWithID create a link component with a user specific id.
func LinkWithID(c *framework.Container, text, url, id string) {
	comp := newLinkComponent(text, url)
	comp.SetID(id)
	c.AddComponent(comp)
}
