package rkt

import (
	"image"

	"github.com/go-gl/gl/v2.1/gl"
)

type Bitmap image.RGBA

type Texture uint32

func (t Texture) bind() {
	gl.BindTexture(gl.TEXTURE_2D, uint32(t))
}

func (b *Bitmap) createTexture() Texture {
	var handle uint32
	size := b.Rect.Size()
	gl.Enable(gl.TEXTURE_2D)
	gl.GenTextures(1, &handle)
	texture := Texture(handle)
	texture.bind()
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.MIRRORED_REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.MIRRORED_REPEAT)
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

	return texture
}
