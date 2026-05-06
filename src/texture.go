package rkt

import (
	"image"

	"github.com/go-gl/gl/v3.3-core/gl"
)

type Bitmap image.RGBA

type TextureFilter int

const (
	TextureFilterNearest = iota
	TextureFilterLinear
)

type Texture uint32

func (t Texture) setFilter(filterType TextureFilter) {
	t.bind2D()
	switch filterType {
	case TextureFilterNearest:
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	case TextureFilterLinear:
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	}
}
func (t Texture) setRepeat(repeatEnable bool) {
	t.bind2D()
	var param int32
	if repeatEnable {
		param = gl.REPEAT
	} else {
		param = gl.CLAMP_TO_BORDER
	}
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, param)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, param)
}
func (t Texture) bind2D() {
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, uint32(t))
}
func (t Texture) bindTo2D(unit uint32) {
	gl.ActiveTexture(gl.TEXTURE0 + unit)
	gl.BindTexture(gl.TEXTURE_2D, uint32(t))
}
func (t Texture) bindCube() {
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_CUBE_MAP, uint32(t))
}
func (t Texture) uniform(location int32, unit uint32) {
	t.bindTo2D(unit)
	gl.Uniform1i(location, int32(unit))
}

func InitTextureUnit(unit uint32) {
	gl.ActiveTexture(gl.TEXTURE0 + unit)
	gl.Enable(gl.TEXTURE_2D)
}

func createTexture2D(b *Bitmap) Texture {
	var handle uint32
	size := b.Rect.Size()
	gl.Enable(gl.TEXTURE_2D)
	gl.GenTextures(1, &handle)
	texture := Texture(handle)
	texture.bind2D()
	texture.setFilter(TextureFilterLinear)
	texture.setRepeat(true)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(size.X),
		int32(size.Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(b.Pix))
	gl.GenerateMipmap(gl.TEXTURE_2D)
	return texture
}

// func createTextureCube(c [6]*Bitmap) Texture {
// 	var handle uint32
// 	size := b.Rect.Size()
// 	gl.Enable(gl.TEXTURE_CUBE_MAP)
// 	gl.GenTextures(1, &handle)
// 	texture := Texture(handle)
// 	texture.bind()
// 	texture.setFilter(TextureFilterLinear)
// 	texture.setRepeat(true)
// 	gl.TexImage2D(
// 		gl.TEXTURE_2D,
// 		0,
// 		gl.RGBA,
// 		int32(size.X),
// 		int32(size.Y),
// 		0,
// 		gl.RGBA,
// 		gl.UNSIGNED_BYTE,
// 		gl.Ptr(b.Pix))
// 	gl.GenerateMipmap(gl.TEXTURE_2D)
// 	return texture
// }
