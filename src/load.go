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
var geomDefMap = make(map[string]*GeomDef, 64)
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

type PlumeDef struct {
	TextureName string  `json:"texture"`
	Offset      Vec3    `json:"offset"`
	Size        float32 `json:"size"`
}

func (d *PlumeDef) create() *Plume {
	p := new(Plume)
	texture, ok := textureMap[d.TextureName]
	if !ok {
		log.Fatalf("build_geom: no such texture %s", d.TextureName)
	}

	g := NewGeom(texture, 12)
	ringVts := buildRingVertices(6, 1.4, -2.0)
	for i := range 6 {
		j := (i + 1) % 6
		ri := i * 6
		g.Vertices[ri+0] = Vec3{0.0, 0.0, 0.0}
		g.Vertices[ri+1] = ringVts[i]
		g.Vertices[ri+2] = ringVts[j]
		g.Vertices[ri+3] = Vec3{0.0, 0.0, -6.0}
		g.Vertices[ri+4] = ringVts[j]
		g.Vertices[ri+5] = ringVts[i]

		g.TexCoords[ri+0] = Vec2{0.5, 0.0}
		g.TexCoords[ri+1] = Vec2{0.0, 0.5}
		g.TexCoords[ri+2] = Vec2{1.0, 0.5}
		g.TexCoords[ri+3] = Vec2{0.5, 1.0}
		g.TexCoords[ri+4] = Vec2{1.0, 0.5}
		g.TexCoords[ri+5] = Vec2{0.0, 0.5}
	}

	p.geom = g
	p.offset = d.Offset
	p.size = d.Size
	return p
}

func loadPartDef(r io.Reader) *PartDef {
	def := new(PartDef)
	dec := json.NewDecoder(r)
	if err := dec.Decode(def); err != nil {
		log.Fatalf("part_def: %v", err)
	}

	if def.Body == nil {
		log.Fatalf("part_def: no body in JSON")
	}

	if def.Attach == nil {
		log.Fatalf("part_def: no attach in JSON")
	}

	if def.Ctrl == nil && def.Decoup == nil && def.Engine == nil {
		log.Fatalf("part_def: no spec field: ctrl/decoup/engine in JSON")
	}

	return def
}

func loadGeomDef(r io.Reader) *GeomDef {
	def := new(GeomDef)
	dec := json.NewDecoder(r)
	if err := dec.Decode(def); err != nil {
		log.Fatalf("load_geom_def: %v\n", err)
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
		case "geom.json":
			log.Printf("+geomdef %s", name)
			geomDefMap[name] = loadGeomDef(fp)
		}

		fp.Close()
	}

	return loadedCount
}
