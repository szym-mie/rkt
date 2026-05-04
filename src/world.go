package rkt

type PVMatrixPair struct {
	ProjMatrix, ViewMatrix Matrix4
}

var ActivePV *PVMatrixPair
