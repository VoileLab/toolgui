package component

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"image"
	"image/png"

	"github.com/mudream4869/toolgui/toolgui/framework"
)

var _ framework.Component = &ImageComponent{}
var ImageComponentName = "image_component"

type ImageComponent struct {
	*framework.BaseComponent
	Src string `json:"src"`
}

func NewImageComponent(src string) *ImageComponent {
	id := fmt.Sprintf("image_%x", md5.Sum([]byte(src)))
	return &ImageComponent{
		BaseComponent: &framework.BaseComponent{
			Name: ImageComponentName,
			ID:   id,
		},
		Src: src,
	}
}

func Image(c *framework.Container, img image.Image) {
	var imageBuf bytes.Buffer
	err := png.Encode(&imageBuf, img)
	if err != nil {
		panic(err)
	}
	bs := imageBuf.Bytes()
	b64 := base64.StdEncoding.EncodeToString(bs)
	src := fmt.Sprintf("data:image/png;base64,%s", b64)
	comp := NewImageComponent(src)
	c.AddComponent(comp)
}

func ImageByURL(c *framework.Container, url string) {
	comp := NewImageComponent(url)
	c.AddComponent(comp)
}
