package tccontent

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/png"

	"github.com/mudream4869/toolgui/toolgui/component/tcutil"
	"github.com/mudream4869/toolgui/toolgui/framework"
)

var _ framework.Component = &imageComponent{}
var imageComponentName = "image_component"

type imageComponent struct {
	*framework.BaseComponent
	Src string `json:"src"`
}

func newImageComponent(src string) *imageComponent {
	return &imageComponent{
		BaseComponent: &framework.BaseComponent{
			Name: imageComponentName,
			ID:   tcutil.HashedID(imageComponentName, []byte(src)),
		},
		Src: src,
	}
}

// Image show a image.
func Image(c *framework.Container, img image.Image) {
	var imageBuf bytes.Buffer
	err := png.Encode(&imageBuf, img)
	if err != nil {
		panic(err)
	}

	bs := imageBuf.Bytes()
	b64 := base64.StdEncoding.EncodeToString(bs)
	src := fmt.Sprintf("data:image/png;base64,%s", b64)
	comp := newImageComponent(src)
	c.AddComponent(comp)
}

// ImageByURI show an image by URI.
func ImageByURI(c *framework.Container, uri string) {
	comp := newImageComponent(uri)
	c.AddComponent(comp)
}
