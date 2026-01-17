package rkt

import (
	"archive/zip"
	"encoding/json"
	"image"
	"image/draw"
	"io"
	"log"
	"math"
	"strings"

	"github.com/go-gl/gl/v2.1/gl"
)

type Texture uint32

func (t Texture) bind() {
	gl.BindTexture(gl.TEXTURE_2D, uint32(t))
}

var textureMap = make(map[string]Texture, 64)
var partDefMap = make(map[string]*PartDef, 64)

func loadTexture(r io.Reader) Texture {
	handle := uint32(0)
	img, _, err := image.Decode(r)
	if err != nil {
		log.Fatalf("texture: %v\n", err)
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		log.Fatalf("texture: bad stride %v\n", rgba.Stride)
	}

	draw.Draw(rgba, rgba.Bounds(), img, image.Pt(0, 0), draw.Src)
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
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix))

	return texture
}

type TaperEndDef struct {
	Offset float32 `json:"offset"`
	Radius float32 `json:"radius"`
}

type TaperDef struct {
	TextureName string      `json:"texture"`
	RingCount   uint16      `json:"rings"`
	Upper       TaperEndDef `json:"upper"`
	Lower       TaperEndDef `json:"lower"`
}

func buildRingVertices(ringCount uint16, radius, z float32) []Vec3 {
	ringVts := make([]Vec3, ringCount)
	for i := range ringCount {
		d := float64(i) / float64(ringCount)
		a := d * 2 * math.Pi
		x := radius * float32(math.Cos(a))
		y := radius * float32(math.Sin(a))
		ringVts[i] = Vec3{x, y, z}
	}

	return ringVts
}

func buildRingEndTexCoords(ringCount uint16, pageX uint8) []Vec2 {
	ringTCs := make([]Vec2, ringCount)
	for i := range ringCount {
		d := float64(i) / float64(ringCount)
		a := d * 2 * math.Pi
		x := float32(math.Cos(a))*0.25 + 0.25
		y := float32(math.Sin(a))*0.25 + 0.25
		ringTCs[i] = Vec2{x + float32(pageX)*0.5, y}
	}

	return ringTCs
}

func buildRingSideTexCoords(ringCount uint16, pageY uint8) []Vec2 {
	sideTCs := make([]Vec2, ringCount)
	for i := range ringCount {
		d := float64(i) / float64(ringCount)
		x := float32(math.Abs(d*2.0 - 1.0))
		y := float32(pageY) * 0.5
		sideTCs[i] = Vec2{x, y}
	}

	return sideTCs
}

func (d *TaperDef) buildGeom() *Geom {
	g := new(Geom)
	rc := d.RingCount
	texture, ok := textureMap[d.TextureName]
	if !ok {
		log.Fatalf("build_geom: no such texture %s\n", d.TextureName)
	}

	g.texture = texture
	endCount := rc * 3
	sideCount := rc * 6
	g.count = uint32(endCount)*2 + uint32(sideCount)
	g.vertices = make([]Vec3, g.count)
	g.texCoords = make([]Vec2, g.count)

	upperVts := buildRingVertices(rc, d.Upper.Radius, d.Upper.Offset)
	lowerVts := buildRingVertices(rc, d.Lower.Radius, d.Lower.Offset)
	upperEndTCs := buildRingEndTexCoords(rc, 0)
	lowerEndTCs := buildRingEndTexCoords(rc, 1)
	upperSideTCs := buildRingSideTexCoords(rc, 1)
	lowerSideTCs := buildRingSideTexCoords(rc, 2)
	for i := range rc {
		j := (i + 1) % rc
		ui := i * 3
		li := ui + endCount
		si := i*6 + 2*endCount

		// end vertex
		g.vertices[ui+0] = Vec3{0.0, 0.0, d.Upper.Offset}
		g.vertices[ui+1] = upperVts[i]
		g.vertices[ui+2] = upperVts[j]
		g.vertices[li+0] = Vec3{0.0, 0.0, d.Lower.Offset}
		g.vertices[li+1] = lowerVts[i]
		g.vertices[li+2] = lowerVts[j]
		// end texcoord
		g.texCoords[ui+0] = Vec2{0.25, 0.25}
		g.texCoords[ui+1] = upperEndTCs[i]
		g.texCoords[ui+2] = upperEndTCs[j]
		g.texCoords[li+0] = Vec2{0.75, 0.25}
		g.texCoords[li+1] = lowerEndTCs[i]
		g.texCoords[li+2] = lowerEndTCs[j]
		// side vertex
		g.vertices[si+0] = upperVts[i]
		g.vertices[si+1] = upperVts[j]
		g.vertices[si+2] = lowerVts[i]
		g.vertices[si+3] = lowerVts[j]
		g.vertices[si+4] = lowerVts[i]
		g.vertices[si+5] = upperVts[j]
		// side texcoord
		g.texCoords[si+0] = upperSideTCs[i]
		g.texCoords[si+1] = upperSideTCs[j]
		g.texCoords[si+2] = lowerSideTCs[i]
		g.texCoords[si+3] = lowerSideTCs[j]
		g.texCoords[si+4] = lowerSideTCs[i]
		g.texCoords[si+5] = upperSideTCs[j]
	}

	return g
}

type PlaneDef struct {
	TextureName string `json:"texture"`
}

func (d *PlaneDef) buildGeom() *Geom {
	g := new(Geom)
	g.texture = textureMap[d.TextureName]
	g.count = 0
	g.vertices = nil
	g.texCoords = nil
	return g
}

type BodyDef struct {
	Tapers []TaperDef `json:"tapers"`
	Planes []PlaneDef `json:"planes"`
}

type AttachPt struct {
	Upper Vec3 `json:"upper"`
	Lower Vec3 `json:"lower"`
}

type CtrlSpec struct {
}

type FuelDef struct {
	Mass    float32 `json:"mass"`
	Impulse float32 `json:"impulse"`
	Flow    float32 `json:"flow"`
}

type PlumeDef struct {
	TextureName string  `json:"texture"`
	Offset      Vec3    `json:"offset"`
	Size        float32 `json:"size"`
}

func (s *PlumeDef) newPlume() *Plume {
	p := new(Plume)
	p.texture = textureMap[s.TextureName]
	p.offset = s.Offset
	p.size = s.Size
	return p
}

type DecoupSpec struct {
	Force float32 `json:"force"`
}

type EngineSpec struct {
	FuelDef     FuelDef  `json:"fuel"`
	PlumeDef    PlumeDef `json:"plume"`
	CanShutdown bool     `json:"can_shutdown"`
}

type PartDef struct {
	Name     string      `json:"name"`
	TypeName string      `json:"type"`
	Mass     float32     `json:"mass"`
	Body     *BodyDef    `json:"body"`
	Attach   *AttachPt   `json:"attach"`
	Ctrl     *CtrlSpec   `json:"ctrl"`
	Decoup   *DecoupSpec `json:"decoup"`
	Engine   *EngineSpec `json:"engine"`
}

func (d *PartDef) newPart() Part {
	switch d.TypeName {
	case "ctrl":
		return d.newCtrl()
	case "engine":
		return d.newEngine()
	case "decoup":
		return d.newDecoup()
	}

	log.Fatalf("new_part_conf: uknown type %v\n", d.TypeName)
	return nil
}
func (d *PartDef) newCtrl() *PartCtrl {
	p := new(PartCtrl)
	p.Def = d
	p.Geom = d.buildGeom()
	p.IsActive = false
	p.IsDead = false
	return p
}
func (d *PartDef) newEngine() *PartEngine {
	p := new(PartEngine)
	p.Def = d
	p.Geom = d.buildGeom()
	p.IsActive = false
	p.IsDead = false
	p.FuelFlow = 0.0
	p.FuelMass = p.Def.Engine.FuelDef.Mass
	p.Plume = *p.Def.Engine.PlumeDef.newPlume()
	return p
}
func (d *PartDef) newDecoup() *PartDecoup {
	p := new(PartDecoup)
	p.Def = d
	p.Geom = d.buildGeom()
	p.IsActive = false
	p.IsDead = false
	p.IsUsed = false
	return p
}
func (d *PartDef) buildGeom() []Geom {
	taperLen := len(d.Body.Tapers)
	geomLen := taperLen + len(d.Body.Planes)
	geom := make([]Geom, geomLen)
	log.Printf("build_geom: %s - %d tapers %d planes\n",
		d.Name, taperLen, len(d.Body.Planes))
	for i, taper := range d.Body.Tapers {
		geom[i] = *taper.buildGeom()
	}

	for i, plane := range d.Body.Planes {
		j := taperLen + i
		geom[j] = *plane.buildGeom()
	}

	return geom
}

func loadPartDef(r io.Reader) *PartDef {
	partDef := new(PartDef)
	dec := json.NewDecoder(r)
	if err := dec.Decode(partDef); err != nil {
		log.Fatalf("part_def: %v\n", err)
	}

	if partDef.Body == nil {
		log.Fatalf("part_def: no body in JSON\n")
	}

	if partDef.Attach == nil {
		log.Fatalf("part_def: no attach in JSON\n")
	}

	if partDef.Ctrl == nil && partDef.Decoup == nil && partDef.Engine == nil {
		log.Fatalf("part_def: no spec field: ctrl/decoup/engine in JSON\n")
	}

	return partDef
}

func NewPart(name string) Part {
	partDef, ok := partDefMap[name]
	if !ok {
		log.Fatalf("new_part: no partdef %s\n", name)
	}

	return partDef.newPart()
}

func LoadPkg(filename string) uint {
	r, err := zip.OpenReader(filename)
	if err != nil {
		log.Fatalf("load_pkg: %v\n", err)
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
			log.Fatalf("load_pkg: %v\n", fp)
		}

		switch suffix {
		case "png", "jpg":
			log.Printf("+texture %s\n", name)
			textureMap[name] = loadTexture(fp)
			loadedCount++
		case "part.json":
			log.Printf("+partdef %s\n", name)
			partDefMap[name] = loadPartDef(fp)
			loadedCount++
		}

		fp.Close()
	}

	return loadedCount
}
