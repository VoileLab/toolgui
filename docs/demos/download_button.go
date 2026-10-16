package demos

import (
	"github.com/voilelab/toolgui/toolgui/tgcomp"
	"github.com/voilelab/toolgui/toolgui/tgframe"
)

func downloadButtonDemo(p *tgframe.Params) error {
	// ANCHOR: demo
	if tgcomp.DownloadButton(
		p.Main, "Download", []byte("123"),
		&tgcomp.DownloadButtonConf{
			Filename: "123.txt",
			Color:    tgcomp.ColorInfo,
		}) {
		tgcomp.Text(p.Main, "Downloaded!")
	}
	// ANCHOR_END: demo
	return nil
}
