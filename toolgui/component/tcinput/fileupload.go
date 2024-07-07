package tcinput

import (
	"github.com/mudream4869/toolgui/toolgui/component/tcutil"
	"github.com/mudream4869/toolgui/toolgui/framework"
)

var _ framework.Component = &fileuploadComponent{}
var fileuploadComponentName = "fileupload_component"

type fileuploadComponent struct {
	*framework.BaseComponent
	Label string `json:"label"`
}

func newFileuploadComponent(label string) *fileuploadComponent {
	return &fileuploadComponent{
		BaseComponent: &framework.BaseComponent{
			Name: fileuploadComponentName,
			ID:   tcutil.NormalID(fileuploadComponentName, label),
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
	comp := newFileuploadComponent(label)
	c.AddComponent(comp)

	var fileObj FileObject
	err := sess.GetObject(comp.ID, &fileObj)
	if err != nil {
		panic(err)
	}

	return fileObj
}
