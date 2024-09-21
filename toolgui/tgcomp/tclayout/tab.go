package tclayout

import (
	"strings"

	"github.com/mudream4869/toolgui/toolgui/tgcomp/tcutil"
	"github.com/mudream4869/toolgui/toolgui/tgframe"
)

var _ tgframe.Component = &tabComponent{}
var tabComponentName = "tab_component"

type tabComponent struct {
	*tgframe.BaseComponent
	Tabs []string `json:"tabs"`
}

func newTabComponent(tabs []string) *tabComponent {
	tabJoinedName := strings.Join(tabs, ",")

	return &tabComponent{
		BaseComponent: &tgframe.BaseComponent{
			Name: tabComponentName,
			ID:   tcutil.NormalID(tabComponentName, tabJoinedName),
		},
		Tabs: tabs,
	}
}

// Tab creates a new tab component with custom configuration
func Tab(c *tgframe.Container, tabs []string) []*tgframe.Container {
	comp := newTabComponent(tabs)
	c.AddComponent(comp)

	ret := make([]*tgframe.Container, len(tabs))
	for i, tab := range tabs {
		ret[i] = tgframe.NewContainer(comp.ID+"_"+tab, c.SendNotifyPack)
		c.SendNotifyPack(tgframe.NewNotifyPackCreate(comp.ID, ret[i]))
	}

	return ret
}

// Tab2 create 2 tabs.
func Tab2(c *tgframe.Container, tabs [2]string) (*tgframe.Container, *tgframe.Container) {
	retTabs := Tab(c, tabs[:])
	return retTabs[0], retTabs[1]
}

// Tab3 create 3 tabs.
func Tab3(c *tgframe.Container, tabs [3]string) (*tgframe.Container, *tgframe.Container, *tgframe.Container) {
	retTabs := Tab(c, tabs[:])
	return retTabs[0], retTabs[1], retTabs[2]
}
