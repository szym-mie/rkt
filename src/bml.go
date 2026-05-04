package rkt

import (
	"encoding/binary"
	"fmt"
	"io"
)

type bmlInit struct {
	Ident     [4]byte
	Caps      [4]byte
	ExternPtr uint32
	AttribPtr uint32
	BufferPtr uint32
	EndPtr    uint32
}

type BMLExtern struct {
	What byte
	Name string
}

type BMLAttrib struct {
	Bindp int8
	Cocnt int8
	Name  string
}

type BMLHeader struct {
	Caps      string
	Externs   []BMLExtern
	Attribs   []BMLAttrib
	ElemCount uint32
}

type BML struct {
	Header *BMLHeader
	Buffer []float32
}

func (b *BMLHeader) ElemLen() uint32 {
	lenSum := int8(0)
	for _, attrib := range b.Attribs {
		lenSum += attrib.Cocnt
	}
	return uint32(lenSum)
}
func (b *BMLHeader) BufferLen() uint32 {
	return b.ElemLen() * b.ElemCount
}

var order = binary.BigEndian

func readAsciiz(r io.Reader, s *string) error {
	buf := make([]byte, 1)
	out := ""
	for {
		n, err := r.Read(buf)
		if err != nil {
			return fmt.Errorf("cannot read ASCIIZ: %w", err)
		}
		if n < 1 {
			return fmt.Errorf("cannot read ASCIIZ: EOF before NUL")
		}
		if buf[0] == 0 {
			*s = out
			return nil
		}
		out += string(buf)
	}
}

func readBMLInit(r io.Reader) (*bmlInit, error) {
	init := new(bmlInit)
	if err := binary.Read(r, order, init); err != nil {
		return nil, fmt.Errorf("cannot read BML init: %w", err)
	}
	return init, nil
}

func readBMLExterns(r io.Reader) ([]BMLExtern, error) {
	count := uint16(0)
	if err := binary.Read(r, order, &count); err != nil {
		return nil, fmt.Errorf("cannot read BML extern count: %w", err)
	}
	externs := make([]BMLExtern, count)
	for i := range count {
		what_name := ""
		if err := readAsciiz(r, &what_name); err != nil {
			return nil, fmt.Errorf("cannot read BML extern: %w", err)
		}
		externs[i].What = what_name[0]
		externs[i].Name = what_name[1:]
	}

	return externs, nil
}

func readBMLAttribs(r io.Reader) ([]BMLAttrib, error) {
	count := uint16(0)
	if err := binary.Read(r, order, &count); err != nil {
		return nil, fmt.Errorf("cannot read BML attrib count: %w", err)
	}
	attribs := make([]BMLAttrib, count)
	for i := range count {
		if err := binary.Read(r, order, &attribs[i].Bindp); err != nil {
			return nil, fmt.Errorf("cannot read BML attrib.bindp: %w", err)
		}
		if err := binary.Read(r, order, &attribs[i].Cocnt); err != nil {
			return nil, fmt.Errorf("cannot read BML attrib.cocnt: %w", err)
		}
		if err := readAsciiz(r, &attribs[i].Name); err != nil {
			return nil, fmt.Errorf("cannot read BML attrib.name: %w", err)
		}
	}

	return attribs, nil
}

func readBMLHeader(r io.Reader) (*BMLHeader, error) {
	init, err := readBMLInit(r)
	if err != nil {
		return nil, err
	}
	externs, err := readBMLExterns(r)
	if err != nil {
		return nil, err
	}
	attribs, err := readBMLAttribs(r)
	if err != nil {
		return nil, err
	}
	header := new(BMLHeader)
	header.Caps = string(init.Caps[:])
	header.Externs = externs
	header.Attribs = attribs
	if err := binary.Read(r, order, &header.ElemCount); err != nil {
		return nil, fmt.Errorf("cannot read BML elem count: %w", err)
	}
	return header, nil
}

func readBMLBuffer(r io.Reader, bufferLen uint32) ([]float32, error) {
	array := make([]float32, bufferLen)
	if err := binary.Read(r, order, &array); err != nil {
		return nil, fmt.Errorf("cannot read BML buffer: %w", err)
	}
	return array, nil
}

func ReadBML(r io.Reader) (*BML, error) {
	header, err := readBMLHeader(r)
	if err != nil {
		return nil, err
	}
	buffer, err := readBMLBuffer(r, header.BufferLen())
	if err != nil {
		return nil, err
	}
	b := new(BML)
	b.Header = header
	b.Buffer = buffer
	return b, nil
}
