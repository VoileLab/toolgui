package tcinput_test

import (
	"github.com/voilelab/toolgui/toolgui/tgcomp/internal/tcinput"
	"github.com/voilelab/toolgui/toolgui/tgframe"
	"github.com/voilelab/toolgui/toolgui/tgjson"
	"testing"
)

func TestFolderPickerProps(t *testing.T) {
	s := tgframe.NewState()
	defer s.Destroy()
	var packs []tgframe.NotifyPack
	c := testContainer(s, &packs)
	conf := &tcinput.FolderPickerConf{ID: "output", Write: true, DefaultDir: "/tmp/output", Disabled: true}
	if f := tcinput.FolderPicker(c, "Output", conf); f != nil {
		t.Fatal("unselected folder must be nil")
	}
	if len(packs) != 1 {
		t.Fatalf("packs = %d", len(packs))
	}
	bs, err := tgjson.Marshal(packs[0])
	if err != nil {
		t.Fatal(err)
	}
	var pack struct {
		Props map[string]any `json:"component"`
	}
	if err := tgjson.Unmarshal(bs, &pack); err != nil {
		t.Fatal(err)
	}
	for k, want := range map[string]any{"name": "folder_picker_component", "id": "folder_picker_component_output", "label": "Output", "write": true, "default_dir": "/tmp/output", "disabled": true} {
		if got := pack.Props[k]; got != want {
			t.Errorf("%s = %v, want %v (pack %s)", k, got, want, bs)
		}
	}
}
