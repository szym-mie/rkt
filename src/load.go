package rkt

import (
	"archive/zip"
	"encoding/json"
	"image"
	"image/draw"
	"io"
	"log"
	"strings"
)

var bitmapMap = make(map[string]*Bitmap, 64)
var shaderMap = make(map[string]Shader, 64)
var textureMap = make(map[string]Texture, 64)
var geom1DefMap = make(map[string]*Geom1Def, 64)
var geom2DefMap = make(map[string]*Geom2Def, 64)
var patchDefMap = make(map[string]*PatchDef, 64)
var partDefMap = make(map[string]*PartDef, 64)

func loadBitmap(r io.Reader) *Bitmap {
	img, _, err := image.Decode(r)
	if err != nil {
		log.Fatalf("load_bitmap: %v", err)
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		log.Fatalf("load_bitmap: bad stride %v", rgba.Stride)
	}

	draw.Draw(rgba, rgba.Bounds(), img, image.Pt(0, 0), draw.Src)
	return (*Bitmap)(rgba)
}

func loadShader(r io.Reader) Shader {
	shader, err := DecodeShader(r)
	if err != nil {
		log.Fatalf("load_shader: %v", err)
	}

	return shader
}

func loadPartDef(r io.Reader) *PartDef {
	def := new(PartDef)
	dec := json.NewDecoder(r)
	if err := dec.Decode(def); err != nil {
		log.Fatalf("load_part_def: %v", err)
	}

	if def.Ctrl == nil && def.Decoup == nil &&
		def.Engine == nil && def.Chute == nil {
		log.Fatalf("load_part_def: no spec field in %s JSON", def.Name)
	}

	return def
}

const maxTextureCnt = 4

// TODO: BML supports up to 2 UV maps
func loadBMLGeom(r io.Reader) (*Geom1Def, *Geom2Def) {
	bml, err := ReadBML(r)
	if err != nil {
		log.Fatalf("load_bml_geom: %v\n", err)
	}

	shaderName := ""
	textureNames := make([]string, maxTextureCnt)
	textureCnt := 0
	for _, extern := range bml.Header.Externs {
		switch extern.What {
		case 'S': // Shader name
			shaderName = extern.Name
		case 'T': // Texture name
			if textureCnt < maxTextureCnt {
				textureNames[textureCnt] = extern.Name
				textureCnt++
			}
		}
	}

	bufferAttrCnt := len(bml.Header.Attribs)
	switch textureCnt {
	case 1:
		def1 := new(Geom1Def)
		def1.ShaderName = shaderName
		def1.TextureName = textureNames[0]
		def1.BufferAttrs = make([]BufferAttr, bufferAttrCnt)
		for i, bmlAttrib := range bml.Header.Attribs {
			def1.BufferAttrs[i].Type = BufferAttrType(bmlAttrib.Bindp)
			def1.BufferAttrs[i].Cocnt = int32(bmlAttrib.Cocnt)
			def1.BufferAttrs[i].Name = bmlAttrib.Name
		}
		def1.RawArray = bml.Buffer
		return def1, nil
	case 3:
		def2 := new(Geom2Def)
		def2.ShaderName = shaderName
		def2.Texture0Name = textureNames[0]
		def2.Texture1Name = textureNames[1]
		def2.Texture2Name = textureNames[2]
		def2.BufferAttrs = make([]BufferAttr, bufferAttrCnt)
		for i, bmlAttrib := range bml.Header.Attribs {
			def2.BufferAttrs[i].Type = BufferAttrType(bmlAttrib.Bindp)
			def2.BufferAttrs[i].Cocnt = int32(bmlAttrib.Cocnt)
			def2.BufferAttrs[i].Name = bmlAttrib.Name
		}
		def2.RawArray = bml.Buffer
		return nil, def2
	default:
		log.Fatalf("load_bml_geom: invalid number of textures %d\n", textureCnt)
	}

	return nil, nil
}

func loadPatchDef(r io.Reader) *PatchDef {
	def := new(PatchDef)
	dec := json.NewDecoder(r)
	if err := dec.Decode(def); err != nil {
		log.Fatalf("load_patch_def: %v\n", err)
	}

	return def
}

func LoadPkg(filename string) uint {
	r, err := zip.OpenReader(filename)
	if err != nil {
		log.Fatalf("load_pkg: %v", err)
	}

	loadedCount := uint(0)

	for _, zf := range r.File {
		path := strings.ReplaceAll(zf.Name, "\\", "/")
		name, suffix, found := strings.Cut(path, ".")
		if !found {
			continue
		}

		fp, err := zf.Open()
		if err != nil {
			log.Fatalf("load_pkg: %v", fp)
		}

		switch suffix {
		case "png", "jpg":
			log.Printf("+texture %s", name)
			bitmap := loadBitmap(fp)
			bitmapMap[name] = bitmap
			textureMap[name] = bitmap.createTexture()
		case "glsl":
			log.Printf("+shader %s", name)
			shaderMap[name] = loadShader(fp)
		case "bml":
			log.Printf("+bml %s", name)
			geom1Def, geom2Def := loadBMLGeom(fp)
			if geom1Def != nil {
				geom1DefMap[name] = geom1Def
			}
			if geom2Def != nil {
				geom2DefMap[name] = geom2Def
			}
		case "part.json":
			log.Printf("+partdef %s", name)
			partDefMap[name] = loadPartDef(fp)
		case "patch.json":
			log.Printf("+patchdef %s", name)
			patchDefMap[name] = loadPatchDef(fp)
		default:
			continue
		}

		loadedCount++
		fp.Close()
	}

	return loadedCount
}
