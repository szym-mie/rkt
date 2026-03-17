package rkt

import "log"

type AttachPt struct {
	Upper Vec3 `json:"upper"`
	Lower Vec3 `json:"lower"`
}

type HullSpec struct {
}

type CtrlSpec struct {
}

type FuelDef struct {
	Mass    float32 `json:"mass"`
	Impulse float32 `json:"impulse"`
	Flow    float32 `json:"flow"`
}
type DecoupSpec struct {
	Force float32 `json:"force"`
}

type EngineSpec struct {
	FuelDef     FuelDef  `json:"fuel"`
	PlumeDef    PlumeDef `json:"plume"`
	CanShutdown bool     `json:"can_shutdown"`
}

type ChuteSpec struct {
	DeployTime float32 `json:"deploy_time"`
	Mass       float32 `json:"mass"`
	Height     float32 `json:"height"`
	Area       float32 `json:"area"`
	Drag       float32 `json:"drag"`
}

// Aero table - AOA, 0 deg, 10 deg,...,180 deg
// lift [ ] + drag [X] vals
// Ctrl src - player input, stage -> to deployed
// State - 0.0 to 1.0
// ChangeSpeed of State

type AeroDef struct {
	Body string `json:"body"`
	// Area []float32 `json:"area"`
	// Drag []float32 `json:"drag"`
}

type PartDef struct {
	Name     string      `json:"name"`
	TypeName string      `json:"type"`
	Mass     float32     `json:"mass"`
	Aero     *AeroDef    `json:"aero"`
	Body     *BodyDef    `json:"body"`
	Attach   *AttachPt   `json:"attach"`
	Ctrl     *CtrlSpec   `json:"ctrl"`
	Decoup   *DecoupSpec `json:"decoup"`
	Engine   *EngineSpec `json:"engine"`
	Chute    *ChuteSpec  `json:"chute"`
}

func (d *PartDef) create() Part {
	switch d.TypeName {
	case "hull":
		return d.createHull()
	case "ctrl":
		return d.createCtrl()
	case "decoup":
		return d.createDecoup()
	case "engine":
		return d.createEngine()
	case "chute":
		return d.createChute()
	}

	log.Fatalf("create: unknown type %v", d.TypeName)
	return nil
}
func (d *PartDef) createHull() *PartHull {
	p := new(PartHull)
	p.Def = d
	p.Geom = d.buildGeom()
	p.IsActive = false
	p.IsDead = false
	return p
}
func (d *PartDef) createCtrl() *PartCtrl {
	p := new(PartCtrl)
	p.Def = d
	p.Geom = d.buildGeom()
	p.IsActive = false
	p.IsDead = false
	return p
}
func (d *PartDef) createDecoup() *PartDecoup {
	p := new(PartDecoup)
	p.Def = d
	p.Geom = d.buildGeom()
	p.IsActive = false
	p.IsDead = false
	p.IsUsed = false
	return p
}
func (d *PartDef) createEngine() *PartEngine {
	p := new(PartEngine)
	p.Def = d
	p.Geom = d.buildGeom()
	p.IsActive = false
	p.IsDead = false
	p.FuelFlow = 0.0
	p.FuelMass = p.Def.Engine.FuelDef.Mass
	p.Plume = *p.Def.Engine.PlumeDef.create()
	return p
}
func (d *PartDef) createChute() *PartChute {
	p := new(PartChute)
	p.Def = d
	p.Geom = d.buildGeom()
	p.IsActive = false
	p.IsDead = false
	p.IsUsed = false
	p.ChuteGeom = geom1DefMap["base/geom/chute0"].create()
	return p
}
func (d *PartDef) buildGeom() []Geom1 {
	taperLen := len(d.Body.Tapers)
	geomLen := taperLen + len(d.Body.Planes)
	geom := make([]Geom1, geomLen)
	log.Printf("build_geom: %s - %d tapers %d planes",
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
