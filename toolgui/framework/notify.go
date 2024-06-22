package framework

const (
	NOTIFY_TYPE_CREATE = 1
	NOTIFY_TYPE_UPDATE = 2
	NOTIFY_TYPE_DELETE = 3
)

type NotifyPack interface {
	GetType() int
}

type NotifyPackBase struct {
	Type int `json:"type"`
}

func (b *NotifyPackBase) GetType() int {
	return b.Type
}

type SendNotifyPackFunc func(NotifyPack)

type NotifyPackCreate struct {
	*NotifyPackBase
	ContainerID string    `json:"container_id"`
	Component   Component `json:"component"`
}

func NewNotifyPackCreate(containrID string, comp Component) NotifyPack {
	return &NotifyPackCreate{
		NotifyPackBase: &NotifyPackBase{
			Type: NOTIFY_TYPE_CREATE,
		},
		ContainerID: containrID,
		Component:   comp,
	}
}

type NotifyPackUpdate struct {
	*NotifyPackBase
	Component Component `json:"component"`
}

func NewNotifyPackUpdate(comp Component) NotifyPack {
	return &NotifyPackUpdate{
		NotifyPackBase: &NotifyPackBase{
			Type: NOTIFY_TYPE_UPDATE,
		},
		Component: comp,
	}
}

type NotifyPackDelete struct {
	*NotifyPackBase
	ComponentID string `json:"component_id"`
}

func NewNotifyPackDelete(compID string) NotifyPack {
	return &NotifyPackDelete{
		NotifyPackBase: &NotifyPackBase{
			Type: NOTIFY_TYPE_DELETE,
		},
		ComponentID: compID,
	}
}
