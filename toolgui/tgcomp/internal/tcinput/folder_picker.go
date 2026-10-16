package tcinput

import (
	"github.com/voilelab/toolgui/toolgui/tgcomp/internal/tcutil"
	"github.com/voilelab/toolgui/toolgui/tgframe"
)

// FolderObject is a user-authorized folder.
//
//tgcomp:export
type FolderObject = tgframe.Folder

// FolderPickerConf configures FolderPicker.
//
//tgcomp:export
type FolderPickerConf struct {
	tgframe.Base
	// Write requests write access. Desktop grants it unconditionally.
	Write bool
	// DefaultDir is the desktop dialog's initial location.
	DefaultDir string
	Disabled   bool
}
type folderPickerComponent struct {
	*tgframe.BaseComponent
	Label      string `json:"label"`
	Write      bool   `json:"write"`
	DefaultDir string `json:"default_dir"`
	Disabled   bool   `json:"disabled"`
}

// FolderPicker draws a native folder selection button; nil until authorized.
//
//tgcomp:export
func FolderPicker(c *tgframe.Container, label string, conf ...*FolderPickerConf) FolderObject {
	cf := tgframe.OneConf("FolderPicker", conf)
	const name = "folder_picker_component"
	comp := &folderPickerComponent{
		BaseComponent: &tgframe.BaseComponent{
			Name: name,
			ID:   tcutil.NormalID(name, label),
		},
		Label:      label,
		Write:      cf.Write,
		DefaultDir: cf.DefaultDir,
		Disabled:   cf.Disabled,
	}
	tgframe.SetConfIDIn(c, comp, cf)
	c.AddComponent(comp)
	return c.State.GetFolder(comp.ID)
}
