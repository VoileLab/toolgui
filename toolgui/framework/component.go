package framework

type Component interface {
	GetID() string
}

var _ Component = &BaseComponent{}

type BaseComponent struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

func (c *BaseComponent) GetID() string {
	return c.ID
}
