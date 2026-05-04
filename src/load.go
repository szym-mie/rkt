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

func loadGeom2Def(r io.Reader) *Geom2Def {
	def := new(Geom2Def)
	dec := json.NewDecoder(r)
	if err := dec.Decode(def); err != nil {
		log.Fatalf("load_geom2_def: %v\n", err)
	}

	return def
}

// TODO: BML supports up to 2 UV maps
func loadBMLGeom(r io.Reader) *Geom1Def {
	bml, err := ReadBML(r)
	if err != nil {
		log.Fatalf("load_bml_geom: %v\n", err)
	}
	def := new(Geom1Def)
	for _, extern := range bml.Header.Externs {
		switch extern.What {
		case 'S': // Shader name
			def.ShaderName = extern.Name
		case 'T': // Texture name
			def.TextureName = extern.Name
		}
	}
	bufferAttrCnt := len(bml.Header.Attribs)
	def.BufferAttrs = make([]BufferAttr, bufferAttrCnt)
	log.Printf("%s+%s\n", def.ShaderName, def.TextureName)
	for i, bmlAttrib := range bml.Header.Attribs {
		log.Printf("%v\n", bmlAttrib)
		def.BufferAttrs[i].Type = BufferAttrType(bmlAttrib.Bindp)
		def.BufferAttrs[i].Cocnt = int32(bmlAttrib.Cocnt)
		def.BufferAttrs[i].Name = bmlAttrib.Name
	}
	def.RawArray = bml.Buffer
	return def
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
			geom1DefMap[name] = loadBMLGeom(fp)
		case "geom2.json":
			log.Printf("+geom2def %s", name)
			geom2DefMap[name] = loadGeom2Def(fp)
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
