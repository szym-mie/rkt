package rkt

import "github.com/go-gl/gl/v2.1/gl"

type BufferAttrType int8

const (
	BufferAttrPos BufferAttrType = iota
	BufferAttrNorm
	BufferAttrTang
	BufferAttrTexCoord0
	BufferAttrTexCoord1
)

type BufferAttr struct {
	Type BufferAttrType
	Name string
	Size int32
}

type Buffer struct {
	shader    Shader
	Attrs     []BufferAttr
	Stride    int32
	Cnt       int32
	vaoHandle uint32
	vboHandle uint32
}

func NewBuffer(shader Shader, attrs []BufferAttr) *Buffer {
	b := new(Buffer)
	gl.GenBuffers(1, &b.vaoHandle)
	gl.GenBuffers(1, &b.vboHandle)
	b.shader = shader
	b.confAttrs(attrs)
	return b
}

func (b *Buffer) confAttrs(attrs []BufferAttr) {
	b.Attrs = attrs
	gl.BindVertexArray(b.vaoHandle)
	offset := uintptr(0)
	stride := int32(0)
	for _, attr := range attrs {
		stride += attr.Size
	}
	b.Stride = stride
	for _, attr := range attrs {
		index := b.shader.getAttrib(attr.Name)
		gl.VertexAttribPointerWithOffset(
			index, attr.Size, gl.FLOAT,
			false, stride*4, offset*4)

		offset += uintptr(attr.Size)
	}
}
func (b *Buffer) Data(data []float32) {
	size := len(data)
	gl.BindBuffer(gl.ARRAY_BUFFER, b.vboHandle)
	gl.BufferData(gl.ARRAY_BUFFER, size*4, gl.Ptr(data), gl.STATIC_DRAW)
	b.Cnt = int32(size) / b.Stride
}
func (b *Buffer) Bind() {
	gl.BindVertexArray(b.vaoHandle)
}
func (b *Buffer) Draw() {
	gl.DrawArrays(gl.TRIANGLES, 0, b.Cnt)
}
