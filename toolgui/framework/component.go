package framework

import "fmt"

// Component is the interface of a component
type Component interface {
	GetID() string
}

var _ Component = &BaseComponent{}

// BaseComponent stores the basic info of a component
type BaseComponent struct {
	// Name is the typename of the component
	Name string `json:"name"`

	// ID shouldn't be duplicated
	ID string `json:"id"`
}

// GetID return component's ID
func (c *BaseComponent) GetID() string {
	return c.ID
}

// SetID set component's ID
func (c *BaseComponent) SetID(id string) {
	c.ID = fmt.Sprintf("%s_%s", c.Name, id)
}
