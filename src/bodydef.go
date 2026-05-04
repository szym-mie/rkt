package rkt

type BodyDef struct {
	InertiaCoeff Vec3     `json:"inertia"`
	GeomDefs     []string `json:"geoms"`
}

// TODO: merge bodyDef with colliders and calcluate axis inertia from defs
