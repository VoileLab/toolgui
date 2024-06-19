package framework

type NotifyAddCompFunc func(containerID string, comp Component)

var _ Component = &Container{}
var ContainerComponentName = "container_component"

// Container contain list of components
type Container struct {
	*BaseComponent

	// Notify adding component under {containerID}
	NotifyAddComp NotifyAddCompFunc `json:"-"`
}

func NewContainer(id string, notifyAddComp NotifyAddCompFunc) *Container {
	return &Container{
		BaseComponent: &BaseComponent{
			Name: ContainerComponentName,
			ID:   id,
		},
		NotifyAddComp: notifyAddComp,
	}
}

func (c *Container) AddComp(comp Component) Component {
	c.NotifyAddComp(c.ID, comp)
	return comp
}

func (c *Container) AddContainer(id string) *Container {
	newContainer := NewContainer(id, c.NotifyAddComp)
	return c.AddComp(newContainer).(*Container)
}
