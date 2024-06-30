package tcinput

import (
	"github.com/mudream4869/toolgui/toolgui/component/tcutil"
	"github.com/mudream4869/toolgui/toolgui/framework"
)

var _ framework.Component = &FileuploadComponent{}
var FileuploadComponentName = "fileupload_component"

type FileuploadComponent struct {
	*framework.BaseComponent
	Label string `json:"label"`
}

func NewFileuploadComponent(label string) *FileuploadComponent {
	return &FileuploadComponent{
		BaseComponent: &framework.BaseComponent{
			Name: FileuploadComponentName,
			ID:   tcutil.NormalID(FileuploadComponentName, label),
		},
		Label: label,
	}
}

type FileObject struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Size int    `json:"size"`
	Body string `json:"body"`
}

func Fileupload(sess *framework.Session, c *framework.Container, label string) FileObject {
	comp := NewFileuploadComponent(label)
	c.AddComponent(comp)

	var fileObj FileObject
	sess.GetObject(comp.ID, &fileObj)
	return fileObj
}
