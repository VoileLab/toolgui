package tgframe

const (
	NotifyTypeCreate = 1
	NotifyTypeUpdate = 2
	NotifyTypeDelete = 3
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
			Type: NotifyTypeCreate,
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
			Type: NotifyTypeUpdate,
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
			Type: NotifyTypeDelete,
		},
		ComponentID: compID,
	}
}
