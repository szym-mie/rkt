package rkt

import (
	"encoding/json"
	"image"
	"image/draw"
	"log"
	"os"

	"github.com/go-gl/gl/v2.1/gl"
)

type Texture uint32

func (t Texture) bind() {
	gl.BindTexture(gl.TEXTURE_2D, uint32(t))
}

var textureMap = make(map[string]Texture, 64)
var partDefMap = make(map[string]PartDef, 64)

func loadTexture(filename string) Texture {
	if texture, ok := textureMap[filename]; ok {
		return texture
	}

	handle := uint32(0)
	fp, err := os.Open(filename)
	if err != nil {
		log.Fatalf("texture: %v\n", err)
	}

	img, _, err := image.Decode(fp)
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
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
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

func (d *TaperDef) buildGeom() *Geom {
	g := new(Geom)
	g.texture = loadTexture(d.TextureName)
	g.count = 24 * uint32(d.RingCount)
	g.vertices = make([]Vec3, g.count)
	g.texCoords = make([]Vec2, g.count)
	return g
}

type PlaneDef struct {
	TextureName string `json:"texture"`
}

func (d *PlaneDef) buildGeom() *Geom {
	g := new(Geom)
	g.texture = loadTexture(d.TextureName)
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
	p.texture = loadTexture(s.TextureName)
	p.offset = s.Offset
	p.size = s.Size
	return p
}

type EngineSpec struct {
	FuelDef     FuelDef  `json:"fuel"`
	PlumeDef    PlumeDef `json:"plume"`
	CanShutdown bool     `json:"can_shutdown"`
}

type DecoupSpec struct {
	Force float32 `json:"force"`
}

type PartDef struct {
	Name     string      `json:"name"`
	TypeName string      `json:"type"`
	Mass     float32     `json:"mass"`
	Body     BodyDef     `json:"body"`
	Attach   AttachPt    `json:"attach"`
	Ctrl     *CtrlSpec   `json:"ctrl"`
	Engine   *EngineSpec `json:"engine"`
	Decoup   *DecoupSpec `json:"decoup"`
}

func (d *PartDef) New() Part {
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
	for i, taper := range d.Body.Tapers {
		geom[i] = *taper.buildGeom()
	}

	for i, plane := range d.Body.Planes {
		j := taperLen + i
		geom[j] = *plane.buildGeom()
	}

	return geom
}

func LoadPartDef(filename string) *PartDef {
	if partDef, ok := partDefMap[filename]; ok {
		return &partDef
	}

	fp, err := os.Open(filename)
	if err != nil {
		log.Fatalf("part_conf: %v\n", err)
	}

	partDef := new(PartDef)
	dec := json.NewDecoder(fp)
	if dec.Decode(partDef) != nil {
		log.Fatalf("part_conf: %v\n", err)
	}

	return partDef
}

type Manifest struct {
	Version string            `json:"ver"`
	Preload map[string]string `json:"pre"`
}

const manifestFilename = "a.manifest.json"

func LoadPath(where string) uint {
	return 0
}
