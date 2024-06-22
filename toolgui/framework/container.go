package framework

var _ Component = &Container{}
var ContainerComponentName = "container_component"

// Container contain list of components
type Container struct {
	*BaseComponent

	SendNotifyPack SendNotifyPackFunc `json:"-"`
}

func NewContainer(id string, notifyComp SendNotifyPackFunc) *Container {
	return &Container{
		BaseComponent: &BaseComponent{
			Name: ContainerComponentName,
			ID:   id,
		},
		SendNotifyPack: notifyComp,
	}
}

func (c *Container) AddComponent(comp Component) Component {
	c.SendNotifyPack(NewNotifyPackCreate(c.ID, comp))
	return comp
}

func (c *Container) AddContainer(id string) *Container {
	newContainer := NewContainer(id, c.SendNotifyPack)
	return c.AddComponent(newContainer).(*Container)
}
