package framework

const (
	NOTIFY_TYPE_CREATE = 1
	NOTIFY_TYPE_UPDATE = 2
	NOTIFY_TYPE_DELETE = 3
)

type NotifyPack interface {
	GetType() int
}

type notifyPackBase struct {
	Type int `json:"type"`
}

func (b *notifyPackBase) GetType() int {
	return b.Type
}

type SendNotifyPackFunc func(NotifyPack)

var _ NotifyPack = &notifyPackCreate{}

type notifyPackCreate struct {
	*notifyPackBase
	ContainerID string    `json:"container_id"`
	Component   Component `json:"component"`
}

func NewNotifyPackCreate(containerID string, comp Component) *notifyPackCreate {
	return &notifyPackCreate{
		notifyPackBase: &notifyPackBase{
			Type: NOTIFY_TYPE_CREATE,
		},
		ContainerID: containerID,
		Component:   comp,
	}
}

var _ NotifyPack = &notifyPackUpdate{}

type notifyPackUpdate struct {
	*notifyPackBase
	Component Component `json:"component"`
}

func NewNotifyPackUpdate(comp Component) *notifyPackUpdate {
	return &notifyPackUpdate{
		notifyPackBase: &notifyPackBase{
			Type: NOTIFY_TYPE_UPDATE,
		},
		Component: comp,
	}
}

var _ NotifyPack = &notifyPackDelete{}

type notifyPackDelete struct {
	*notifyPackBase
	ComponentID string `json:"component_id"`
}

func NewNotifyPackDelete(compID string) *notifyPackDelete {
	return &notifyPackDelete{
		notifyPackBase: &notifyPackBase{
			Type: NOTIFY_TYPE_DELETE,
		},
		ComponentID: compID,
	}
}
