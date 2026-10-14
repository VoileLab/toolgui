package main

import (
	"fmt"
	"io"

	"github.com/voilelab/toolgui/toolgui/tgcomp"
	"github.com/voilelab/toolgui/toolgui/tgframe"
)

// FolderPage works with the same grant across reruns without upload/download.
func FolderPage(p *tgframe.Params) error {
	tgcomp.Title(p.Main, "Local folder")
	write := tgcomp.Checkbox(p.Main, "Request write access")
	folder := tgcomp.FolderPicker(p.Main, "Choose local folder", &tgcomp.FolderPickerConf{Write: write})
	if folder == nil {
		return nil
	}
	access := "read-only"
	if folder.Writable() {
		access = "readwrite"
	}
	tgcomp.Text(p.Main, "Selected folder: "+access)
	entries, err := folder.ReadDir(".")
	if err != nil {
		return err
	}
	for _, entry := range entries {
		tgcomp.Text(p.Main, "Entry: "+entry.Name())
	}
	name := tgcomp.Textbox(p.Main, "File name", &tgcomp.TextboxConf{Default: "input.txt"})
	if tgcomp.Button(p.Main, "Read local file") {
		info, err := folder.Stat(name)
		if err != nil {
			return err
		}
		file, err := folder.Open(name)
		if err != nil {
			return err
		}
		defer file.Close()
		body, err := io.ReadAll(io.LimitReader(file, 1024))
		if err != nil {
			return err
		}
		tgcomp.Text(p.Main, fmt.Sprintf("File size: %d", info.Size()))
		tgcomp.Text(p.Main, "Local content: "+string(body))
	}
	if folder.Writable() {
		if tgcomp.Button(p.Main, "Create directories") {
			if err := folder.MkdirAll("nested/deep"); err != nil {
				return err
			}
			tgcomp.Text(p.Main, "Directories created")
		}
		if tgcomp.Button(p.Main, "Write local file") {
			file, err := folder.Create(name)
			if err != nil {
				return err
			}
			if _, err := io.WriteString(file, "written by Go\n"); err != nil {
				_ = file.Close()
				return err
			}
			if err := file.Close(); err != nil {
				return err
			}
			tgcomp.Text(p.Main, "Local file written")
		}
	}
	return nil
}
