package rkt

import (
	"github.com/go-gl/gl/v3.3-core/gl"
)

type BufferAttrType int8

const (
	BufferAttrPos  BufferAttrType = 0
	BufferAttrNorm BufferAttrType = 1
	BufferAttrTang BufferAttrType = 2
	BufferAttrUV0  BufferAttrType = 4
	BufferAttrUV1  BufferAttrType = 5
)

type BufferAttr struct {
	Type  BufferAttrType // type of buffer attr
	Name  string         // shader location name
	Cocnt int32          // number of components
}

type Buffer struct {
	shader    Shader
	Attrs     []BufferAttr
	Stride1   int32 // stride in bytes
	Stride4   int32 // stride in floats (Stride1 / 4)
	Cnt       int32
	vaoHandle uint32
	vboHandle uint32
}

func NewBuffer(shader Shader, attrs []BufferAttr) *Buffer {
	b := new(Buffer)
	gl.GenVertexArrays(1, &b.vaoHandle)
	gl.GenBuffers(1, &b.vboHandle)
	b.shader = shader
	b.confAttrs(attrs)
	return b
}

func (b *Buffer) confAttrs(attrs []BufferAttr) {
	b.Attrs = attrs
	// gl.BindVertexArray(b.vaoHandle)
	// gl.BindBuffer(gl.ARRAY_BUFFER, b.vboHandle)
	// offset := uintptr(0)
	stride := int32(0)
	for _, attr := range attrs {
		// only floats supported
		stride += attr.Cocnt * 4
	}
	b.Stride1 = stride
	b.Stride4 = stride / 4
	// for _, attr := range attrs {
	// 	// TOOD: if only using fixed locations:
	// 	// unnecessary - the code is executed very rarely
	// 	// index := attr.Type
	// 	index := b.shader.getAttrib(attr.Name)
	// 	gl.VertexAttribPointerWithOffset(
	// 		uint32(index), attr.Size, gl.FLOAT,
	// 		false, stride, offset)
	// 	gl.EnableVertexAttribArray(uint32(index))
	// 	offset += uintptr(attr.Size * 4)
	// }
}
func (b *Buffer) newDataArray(count int) []float32 {
	return make([]float32, int32(count)*b.Stride4)
}
func (b *Buffer) data(data []float32) {
	size := len(data)
	gl.BindVertexArray(b.vaoHandle)
	gl.BindBuffer(gl.ARRAY_BUFFER, b.vboHandle)
	gl.BufferData(gl.ARRAY_BUFFER, size*4, gl.Ptr(data), gl.STATIC_DRAW)
	b.Cnt = int32(size) / b.Stride4
	offset := uintptr(0)
	stride := int32(0)
	for _, attr := range b.Attrs {
		// only floats supported
		stride += attr.Cocnt * 4
	}
	b.Stride1 = stride
	b.Stride4 = stride / 4
	for _, attr := range b.Attrs {
		// TODO: if only using fixed locations:
		// unnecessary - the code is executed very rarely
		// index := attr.Type
		index := b.shader.getAttrib(attr.Name)
		gl.VertexAttribPointerWithOffset(
			uint32(index), attr.Cocnt, gl.FLOAT,
			false, stride, offset)
		gl.EnableVertexAttribArray(uint32(index))
		offset += uintptr(attr.Cocnt * 4)
	}
}
func (b *Buffer) bind() {
	gl.BindVertexArray(b.vaoHandle)
}
func (b *Buffer) draw() {
	gl.DrawArrays(gl.TRIANGLES, 0, b.Cnt)
}
