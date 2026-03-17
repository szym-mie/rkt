package rkt

// A cylinder based collider, tapered towards one of the ends.
// Both end circles have their centers and normals lie on the same axis.
type Taper struct {
	Pos         Vec3
	height      float32
	upperRadius float32
	lowerRadius float32
}

func (t *Taper) DistTo(p Vec3) float32 {
	// TODO
	return 1.0
}
func (t *Taper) TestPoint(p Vec3) bool {
	lp := p.Sub(t.Pos)
	// Early Z checks
	if lp.Z > t.height || lp.Z < 0.0 {
		return false
	}

	ratio := lp.Z / t.height
	// Local Z radius, local point length squared
	lzRadius := ratio*t.upperRadius + (1.0-ratio)*t.lowerRadius
	lpLenSq := lp.X*lp.X + lp.Y*lp.Y
	return lpLenSq < lzRadius*lzRadius
}

// A sphere collider with a radius. Computationally cheapest to use -
// best choice for a proximity precheck.
type Sphere struct {
	Pos      Vec3
	radiusSq float32
}

func (s *Sphere) DistTo(p Vec3) float32 {
	return p.Sub(s.Pos).LenSq() - s.radiusSq
}
func (s *Sphere) TestPoint(p Vec3) bool {
	return p.Sub(s.Pos).LenSq() < s.radiusSq
}

// A plane collider, forms a quad.
type Quad struct {
	Pos Vec3
}

func (q *Quad) TestPoint(p Vec3) bool {
	return false
}
