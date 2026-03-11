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
var textureMap = make(map[string]Texture, 64)
var geom1DefMap = make(map[string]*Geom1Def, 64)
var geom2DefMap = make(map[string]*Geom2Def, 64)
var patchDefMap = make(map[string]*PatchDef, 64)
var partDefMap = make(map[string]*PartDef, 64)

func loadBitmap(r io.Reader) *Bitmap {
	img, _, err := image.Decode(r)
	if err != nil {
		log.Fatalf("load_image: %v", err)
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		log.Fatalf("bitmap: bad stride %v", rgba.Stride)
	}

	draw.Draw(rgba, rgba.Bounds(), img, image.Pt(0, 0), draw.Src)
	return (*Bitmap)(rgba)
}

func loadPartDef(r io.Reader) *PartDef {
	def := new(PartDef)
	dec := json.NewDecoder(r)
	if err := dec.Decode(def); err != nil {
		log.Fatalf("part_def: %v", err)
	}

	if def.Body == nil {
		log.Fatalf("part_def: no body in %s JSON", def.Name)
	}

	if def.Attach == nil {
		log.Fatalf("part_def: no attach in %s JSON", def.Name)
	}

	if def.Ctrl == nil && def.Decoup == nil &&
		def.Engine == nil && def.Chute == nil {
		log.Fatalf("part_def: no spec field in %s JSON", def.Name)
	}

	return def
}

func loadGeom1Def(r io.Reader) *Geom1Def {
	def := new(Geom1Def)
	dec := json.NewDecoder(r)
	if err := dec.Decode(def); err != nil {
		log.Fatalf("load_geom1_def: %v\n", err)
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
			loadedCount++
		case "part.json":
			log.Printf("+partdef %s", name)
			partDefMap[name] = loadPartDef(fp)
			loadedCount++
		case "geom1.json":
			log.Printf("+geom1def %s", name)
			geom1DefMap[name] = loadGeom1Def(fp)
		case "geom2.json":
			log.Printf("+geom2def %s", name)
			geom2DefMap[name] = loadGeom2Def(fp)
		case "patch.json":
			log.Printf("+patchdef %s", name)
			patchDefMap[name] = loadPatchDef(fp)
		}

		fp.Close()
	}

	return loadedCount
}
