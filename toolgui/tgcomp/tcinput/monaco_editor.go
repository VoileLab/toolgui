package tcinput

import (
	"github.com/VoileLab/toolgui/toolgui/tgcomp/tcutil"
	"github.com/VoileLab/toolgui/toolgui/tgframe"
)

var _ tgframe.Component = &monacoEditorComponent{}
var monacoEditorComponentName = "monaco_editor_component"

type monacoEditorComponent struct {
	*tgframe.BaseComponent
}

func newMonacoEditorComponent(id string) *monacoEditorComponent {
	return &monacoEditorComponent{
		BaseComponent: &tgframe.BaseComponent{
			Name: monacoEditorComponentName,
			ID:   tcutil.NormalID(monacoEditorComponentName, id),
		},
	}
}

func MonacoEditor(s *tgframe.State, c *tgframe.Container, id string) string {
	comp := newMonacoEditorComponent(id)
	c.AddComponent(comp)
	return s.GetString(comp.ID, "")
}
