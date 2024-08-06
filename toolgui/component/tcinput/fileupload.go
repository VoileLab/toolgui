package tcinput

import (
	"github.com/mudream4869/toolgui/toolgui/component/tcutil"
	"github.com/mudream4869/toolgui/toolgui/framework"
)

var _ framework.Component = &fileuploadComponent{}
var fileuploadComponentName = "fileupload_component"

type fileuploadComponent struct {
	*framework.BaseComponent
	Label  string `json:"label"`
	Accept string `json:"accept"`
}

func newFileuploadComponent(label, accept string) *fileuploadComponent {
	return &fileuploadComponent{
		BaseComponent: &framework.BaseComponent{
			Name: fileuploadComponentName,
			ID:   tcutil.NormalID(fileuploadComponentName, label),
		},
		Label:  label,
		Accept: accept,
	}
}

type FileObject struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Size int    `json:"size"`

	Bytes []byte `json:"_"`
}

// Fileupload create a fileupload and return its selected file.
func Fileupload(s *framework.State, c *framework.Container, label, accept string) FileObject {
	comp := newFileuploadComponent(label, accept)
	c.AddComponent(comp)

	var fileObj FileObject
	err := s.GetObject(comp.ID, &fileObj)
	if err != nil {
		panic(err)
	}

	fileObj.Bytes = s.GetFile(fileObj.Name)

	return fileObj
}
