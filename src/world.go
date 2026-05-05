package rkt

type PVMatrix struct {
	ProjMatrix, ViewMatrix Matrix4
}

var ActivePV *PVMatrix

type LightEnv struct {
	AmbColor  Vec3
	DirLights [2]DirLight
	PtLights  [8]PtLight
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
