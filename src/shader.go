package rkt

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/go-gl/gl/v2.1/gl"
)

type ShaderUnitType uint

const (
	shaderUnitVert ShaderUnitType = iota + 1
	shaderUnitFrag
)

func (ut ShaderUnitType) String() string {
	switch ut {
	case shaderUnitVert:
		return "vertex"
	case shaderUnitFrag:
		return "fragment"
	}

	return "unknown"
}

type ShaderUnit uint32

func compileShader(unitSrc string, unitType ShaderUnitType) (ShaderUnit, error) {
	var glType uint32
	switch unitType {
	case shaderUnitVert:
		glType = gl.VERTEX_SHADER
	case shaderUnitFrag:
		glType = gl.FRAGMENT_SHADER
	}

	if glType == 0 {
		return 0, fmt.Errorf("compile_shader: bad type")
	}

	h := gl.CreateShader(glType)
	u := ShaderUnit(h)
	unitSrcC, free := gl.Strs(unitSrc + "\x00")
	defer free()

	gl.ShaderSource(h, 1, unitSrcC, nil)
	gl.CompileShader(h)

	var compileStatus int32
	gl.GetShaderiv(h, gl.COMPILE_STATUS, &compileStatus)
	if compileStatus == gl.FALSE {
		var logLen int32
		gl.GetShaderiv(h, gl.INFO_LOG_LENGTH, &logLen)

		log := strings.Repeat("\x00", int(logLen+1))
		gl.GetShaderInfoLog(h, logLen, nil, gl.Str(log))

		return 0, fmt.Errorf(
			"compile_shader: %v compile failed:\n%v\n///\n%v",
			unitType, log, unitSrc)
	}

	return u, nil
}

func (u ShaderUnit) attachTo(s Shader) {
	gl.AttachShader(uint32(s), uint32(u))
}

func (u ShaderUnit) delete() {
	gl.DeleteShader(uint32(u))
}

type Shader uint32

func linkShader(vertUnit, fragUnit ShaderUnit) (Shader, error) {
	h := gl.CreateProgram()
	s := Shader(h)
	vertUnit.attachTo(s)
	fragUnit.attachTo(s)
	gl.LinkProgram(h)
	var linkStatus int32
	gl.GetProgramiv(h, gl.LINK_STATUS, &linkStatus)
	if linkStatus == gl.FALSE {
		var logLen int32
		gl.GetProgramiv(h, gl.INFO_LOG_LENGTH, &logLen)

		log := strings.Repeat("\x00", int(logLen+1))
		gl.GetProgramInfoLog(h, logLen, nil, gl.Str(log))

		return 0, fmt.Errorf("create_program: link failed:\n%v", log)
	}

	return s, nil
}

func NewShader(vertSrc, fragSrc string) (Shader, error) {
	vert, err := compileShader(vertSrc, shaderUnitVert)
	if err != nil {
		return 0, fmt.Errorf("new_shader: %v", err)
	}
	defer vert.delete()

	frag, err := compileShader(fragSrc, shaderUnitFrag)
	if err != nil {
		return 0, fmt.Errorf("new_shader: %v", err)
	}
	defer frag.delete()

	s, err := linkShader(vert, frag)
	if err != nil {
		return 0, fmt.Errorf("new_shader: %v", err)
	}

	return s, nil
}

func (s Shader) getAttrib(name string) int32 {
	return gl.GetAttribLocation(uint32(s), gl.Str(name+"\x00"))
}
func (s Shader) getUniform(name string) int32 {
	return gl.GetUniformLocation(uint32(s), gl.Str(name+"\x00"))
}
func (s Shader) active() {
	gl.UseProgram(uint32(s))
}

type shaderLoadState uint

const (
	shaderLoadNoneState shaderLoadState = iota
	shaderLoadVertState
	shaderLoadFragState
)

func DecodeShader(r io.Reader) (Shader, error) {
	sc := bufio.NewScanner(r)
	vertSrc, fragSrc := "", ""
	loadState := shaderLoadNoneState
	for sc.Scan() {
		line := sc.Text()
		if strings.HasPrefix(line, "//+vert") {
			loadState = shaderLoadVertState
			continue
		}
		if strings.HasPrefix(line, "//+frag") {
			loadState = shaderLoadFragState
			continue
		}

		switch loadState {
		case shaderLoadVertState:
			vertSrc += line + "\n"
		case shaderLoadFragState:
			fragSrc += line + "\n"
		}
	}

	s, err := NewShader(vertSrc, fragSrc)
	if err != nil {
		return 0, err
	}

	return s, nil
}
