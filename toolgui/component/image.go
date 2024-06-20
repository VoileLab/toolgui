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
	Base64Image string `json:"base64_image"`
	Format      string `json:"format"`
}

func NewImageComponent(img image.Image) (*ImageComponent, error) {
	var imageBuf bytes.Buffer
	err := png.Encode(&imageBuf, img)
	if err != nil {
		return nil, err
	}

	bs := imageBuf.Bytes()
	id := fmt.Sprintf("image_%x", md5.Sum(bs))
	b64 := base64.StdEncoding.EncodeToString(bs)

	return &ImageComponent{
		BaseComponent: &framework.BaseComponent{
			Name: ImageComponentName,
			ID:   id,
		},
		Base64Image: b64,
		Format:      "png",
	}, nil
}

func Image(c *framework.Container, img image.Image) {
	comp, err := NewImageComponent(img)
	if err != nil {
		panic(err)
	}
	c.AddComp(comp)
}
