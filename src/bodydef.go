package rkt

import (
	"log"
)

type BodyDef struct {
	InertiaCoeff Vec3       `json:"inertia"`
	Tapers       []TaperDef `json:"tapers"`
	Planes       []PlaneDef `json:"planes"`
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

func (d *TaperDef) buildGeom() *Geom1 {
	rc := d.RingCount
	texture, ok := textureMap[d.TextureName]
	if !ok {
		log.Fatalf("build_geom: no such texture %s\n", d.TextureName)
	}

	g := NewGeom1(texture, rc*4)
	endCount := rc * 3

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
		g.Vertices[ui+0] = Vec3{0.0, 0.0, d.Upper.Offset}
		g.Vertices[ui+1] = upperVts[i]
		g.Vertices[ui+2] = upperVts[j]
		g.Vertices[li+0] = Vec3{0.0, 0.0, d.Lower.Offset}
		g.Vertices[li+1] = lowerVts[j]
		g.Vertices[li+2] = lowerVts[i]
		// end texcoord
		g.TexCoords[ui+0] = Vec2{0.25, 0.25}
		g.TexCoords[ui+1] = upperEndTCs[i]
		g.TexCoords[ui+2] = upperEndTCs[j]
		g.TexCoords[li+0] = Vec2{0.75, 0.25}
		g.TexCoords[li+1] = lowerEndTCs[j]
		g.TexCoords[li+2] = lowerEndTCs[i]
		// side vertex
		g.Vertices[si+0] = upperVts[j]
		g.Vertices[si+1] = upperVts[i]
		g.Vertices[si+2] = lowerVts[i]
		g.Vertices[si+3] = lowerVts[i]
		g.Vertices[si+4] = lowerVts[j]
		g.Vertices[si+5] = upperVts[j]
		// side texcoord
		g.TexCoords[si+0] = upperSideTCs[j]
		g.TexCoords[si+1] = upperSideTCs[i]
		g.TexCoords[si+2] = lowerSideTCs[i]
		g.TexCoords[si+3] = lowerSideTCs[i]
		g.TexCoords[si+4] = lowerSideTCs[j]
		g.TexCoords[si+5] = upperSideTCs[j]
	}

	return g
}

type PlaneDef struct {
	TextureName string `json:"texture"`
}

func (d *PlaneDef) buildGeom() *Geom1 {
	g := new(Geom1)
	g.Texture = textureMap[d.TextureName]
	g.Count = 0
	g.Vertices = nil
	g.TexCoords = nil
	return g
}
