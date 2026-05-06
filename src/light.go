package rkt

type DirLight struct {
	Dir, Color Vec3
}

type PtLight struct {
	Pos, Color, Power Vec3
}
