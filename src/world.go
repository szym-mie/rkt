package rkt

import (
	"math"
	"time"
)

type PVMatrix struct {
	ProjMatrix, ViewMatrix Mat4
}

var ActivePV *PVMatrix

type LightEnv struct {
	AmbColor  Vec3
	DirLights [2]DirLight
	PtLights  [8]PtLight
}

const SECONDS = 24.0 * 60.0 * 60.0

func (l *LightEnv) SetFromExtern(t time.Time, geoPos Vec3) {
	amb := &l.AmbColor
	sun := &l.DirLights[0]
	alb := &l.DirLights[1]
	h, m, s := t.Clock()
	theta := float64(h*3600.0+m*60.0+s) / SECONDS * 2.0 * math.Pi
	z := float32(-math.Cos(theta))
	x := float32(math.Sin(theta))

	ambBlack := Vec3{0.1, 0.1, 0.2}
	ambWhite := Vec3{0.4, 0.4, 0.5}

	sunBlack := Vec3{0.0, 0.0, 0.0}
	sunOrange := Vec3{1.2, 0.4, 0.1}
	sunWhite := Vec3{0.9, 0.9, 1.0}

	albBlack := Vec3{0.0, 0.0, 0.0}
	albWhite := Vec3{0.0, 0.1, 0.1}

	sun.Dir.X = float32(x)
	sun.Dir.Y = 0.0
	sun.Dir.Z = float32(z)
	if z < -0.1 {
		*amb = ambBlack
		alb.Color = albBlack
	} else if z > 0.2 {
		*amb = ambWhite
		alb.Color = albWhite
	} else {
		weight := (z + 0.1) * 3.33
		*amb = ambBlack.Lerp(ambWhite, weight)
		alb.Color = albBlack.Lerp(albWhite, weight)
	}
	if z < -0.1 {
		sun.Color = sunBlack

	} else if z > 0.5 {
		sun.Color = sunWhite

	} else if z > 0.1 {
		weight := (z - 0.1) * 2.5
		sun.Color = sunOrange.Lerp(sunWhite, weight)
	} else {
		weight := (z + 0.1) * 5.0
		sun.Color = sunBlack.Lerp(sunOrange, weight)
	}
}

func (l *LightEnv) uniform(ambLocation, dirDirLocation, dirColorLocation int32) {
	l.AmbColor.uniform(ambLocation)
	for i, dirLight := range l.DirLights {
		offset := int32(i)
		dirLight.Dir.uniform(dirDirLocation + offset)
		dirLight.Color.uniform(dirColorLocation + offset)
	}
}

var ActiveLightEnv LightEnv
