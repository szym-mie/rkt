package rkt

import (
	"log"
)

type BodyDef struct {
	InertiaCoeff Vec3     `json:"inertia"`
	GeomDefs     []string `json:"geoms"`
}

// TODO: merge bodyDef with colliders and calcluate axis inertia from defs

type TaperEndDef struct {
	Offset float32 `json:"offset"`
	Radius float32 `json:"radius"`
}

type TaperDef struct {
	TextureName string      `json:"texture"`
	RingCount   int         `json:"rings"`
	Upper       TaperEndDef `json:"upper"`
	Lower       TaperEndDef `json:"lower"`
}

const bodyShaderName = "base/glsl/phong1"

func setVec2At(arr []float32, i *int, v Vec2) {
	arr[*i+0] = v.X
	arr[*i+1] = v.Y
	*i += 2
}
func setVec3At(arr []float32, i *int, v Vec3) {
	arr[*i+0] = v.X
	arr[*i+1] = v.Y
	arr[*i+2] = v.Z
	*i += 3
}

func (d *TaperDef) buildGeom() *Geom1 {
	rc := d.RingCount
	shader, ok := shaderMap[bodyShaderName]
	if !ok {
		log.Fatalf("build_geom: no such shader %s\n", bodyShaderName)
	}

	texture, ok := textureMap[d.TextureName]
	if !ok {
		log.Fatalf("build_geom: no such texture %s\n", d.TextureName)
	}

	g := NewGeom1(shader, texture)
	bufferData := g.Buffer.newDataArray(rc * 4 * 3)

	upperVts := buildRingVertices(rc, d.Upper.Radius, d.Upper.Offset)
	lowerVts := buildRingVertices(rc, d.Lower.Radius, d.Lower.Offset)
	upperEndTCs := buildRingEndTexCoords(rc, 0)
	lowerEndTCs := buildRingEndTexCoords(rc, 1)
	upperSideTCs := buildRingSideTexCoords(rc, 1)
	lowerSideTCs := buildRingSideTexCoords(rc, 2)
	i := 0
	for j := range rc {
		k := (j + 1) % rc

		// end vertex
		setVec3At(bufferData, &i, Vec3{0.0, 0.0, d.Upper.Offset})
		setVec2At(bufferData, &i, Vec2{0.25, 0.25})
		setVec3At(bufferData, &i, upperVts[j])
		setVec2At(bufferData, &i, upperEndTCs[j])
		setVec3At(bufferData, &i, upperVts[k])
		setVec2At(bufferData, &i, upperEndTCs[k])
		setVec3At(bufferData, &i, Vec3{0.0, 0.0, d.Lower.Offset})
		setVec2At(bufferData, &i, Vec2{0.75, 0.25})
		setVec3At(bufferData, &i, lowerVts[k])
		setVec2At(bufferData, &i, lowerEndTCs[k])
		setVec3At(bufferData, &i, lowerVts[j])
		setVec2At(bufferData, &i, lowerEndTCs[j])
		// side
		setVec3At(bufferData, &i, upperVts[k])
		setVec2At(bufferData, &i, upperSideTCs[k])
		setVec3At(bufferData, &i, upperVts[j])
		setVec2At(bufferData, &i, upperSideTCs[j])
		setVec3At(bufferData, &i, lowerVts[j])
		setVec2At(bufferData, &i, lowerSideTCs[j])
		setVec3At(bufferData, &i, lowerVts[j])
		setVec2At(bufferData, &i, lowerSideTCs[j])
		setVec3At(bufferData, &i, lowerVts[k])
		setVec2At(bufferData, &i, lowerSideTCs[k])
		setVec3At(bufferData, &i, upperVts[k])
		setVec2At(bufferData, &i, upperSideTCs[k])
	}
	g.Buffer.data(bufferData)

	return g
}

type PlaneDef struct {
	TextureName string `json:"texture"`
}

func (d *PlaneDef) buildGeom() *Geom1 {
	g := new(Geom1)
	g.Texture = textureMap[d.TextureName]
	g.Count = 0
	return g
}
