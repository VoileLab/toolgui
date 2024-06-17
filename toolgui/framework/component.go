package framework

type Component interface {
	GetID() string

	// iterate all child-component
	IterComp(yield func(Component))
}

var _ Component = &BaseComponent{}

type BaseComponent struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

func (c *BaseComponent) GetID() string {
	return c.ID
}

func (c *BaseComponent) IterComp(func(Component)) {
}
