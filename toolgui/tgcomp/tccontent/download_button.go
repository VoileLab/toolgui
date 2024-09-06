package tccontent

import (
	"encoding/base64"

	"github.com/mudream4869/toolgui/toolgui/tgcomp/tcutil"
	"github.com/mudream4869/toolgui/toolgui/tgframe"
)

var _ tgframe.Component = &downloadButtonComponent{}
var downloadButtonComponentName = "download_button_component"

type downloadButtonComponent struct {
	*tgframe.BaseComponent
	Text       string `json:"text"`
	Base64Body string `json:"base64_body"`
	Filename   string `json:"filename"`
}

func newDownloadButtonComponent(text, base64Body, filename string) *downloadButtonComponent {
	return &downloadButtonComponent{
		BaseComponent: &tgframe.BaseComponent{
			Name: downloadButtonComponentName,
			ID:   tcutil.NormalID(downloadButtonComponentName, text),
		},
		Text:       text,
		Base64Body: base64Body,
		Filename:   filename,
	}
}

// DownloadButton create a download button component.
func DownloadButton(c *tgframe.Container, text string, body []byte, filename string) {
	b64Body := base64.RawStdEncoding.EncodeToString(body)
	c.AddComponent(newDownloadButtonComponent(text, b64Body, filename))
}

// DownloadButtonWithID create a download button component with a user specific id.
func DownloadButtonWithID(c *tgframe.Container, text string, body []byte, filename, id string) {
	b64Body := base64.RawStdEncoding.EncodeToString(body)
	comp := newDownloadButtonComponent(text, b64Body, filename)
	comp.SetID(id)
	c.AddComponent(comp)
}
