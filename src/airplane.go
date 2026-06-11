package rkt

import (
	"math"
)

type RigidBody struct {
	forceAcc   Vec3
	torqueAcc  Vec3
	Pos        Vec3
	Rot        Quat
	Vel        Vec3
	Ang        Vec3
	Inertia    *Mat3
	InertiaInv *Mat3
	Mass       float32
}

func NewRigidBody(mass float32, inertia *Mat3) *RigidBody {
	r := new(RigidBody)
	r.Rot = ZeroQuat()
	r.Mass = mass
	r.Inertia = inertia
	r.InertiaInv = inertia.Invert()
	return r
}

func (r *RigidBody) IntoWorldSpace(v Vec3) Vec3 {
	return r.Rot.Rotate(v)
}
func (r *RigidBody) IntoBodySpace(v Vec3) Vec3 {
	return r.Rot.Conj().Rotate(v)
}
func (r *RigidBody) VelAt(pt Vec3) Vec3 {
	return r.IntoBodySpace(r.Vel).Add(r.Ang.Cross(pt))
}
func (r *RigidBody) AddBodyForceAt(force, pt Vec3) {
	r.forceAcc.AddSelf(r.IntoWorldSpace(force))
	r.torqueAcc.AddSelf(pt.Cross(force))
}
func (r *RigidBody) AddRelForce(force Vec3) {
	r.forceAcc.AddSelf(r.IntoWorldSpace(force))
}
func (r *RigidBody) Update(dt float32) {
	massInv := 1 / r.Mass
	accel := r.forceAcc.MulSca(massInv)
	accel.Z -= 9.81
	r.Vel.AddSelf(accel.MulSca(dt))
	r.Pos.AddSelf(r.Vel.MulSca(dt))
	localAng := r.Ang.Cross(r.Inertia.Apply(r.Ang))
	r.Ang.AddSelf(r.InertiaInv.Apply(r.torqueAcc.Sub(localAng)).MulSca(dt))
	w := Quat{0.0, r.Ang.X, r.Ang.Y, r.Ang.Z}
	r.Rot.AddSelf(r.Rot.Product(w).Scale(0.5 * dt))
	r.Rot = r.Rot.Norm()
	r.forceAcc.SetZero()
	r.torqueAcc.SetZero()
}

type Airplane struct {
	RigidBody
	Engine *Engine
	Wings  []Wing
}

func (a *Airplane) Update(dt float32) {
	a.Engine.Update(dt)
	a.Engine.ApplyForce(&a.RigidBody)
	for _, wing := range a.Wings {
		wing.ApplyForce(&a.RigidBody)
	}
	a.RigidBody.Update(dt)
}

type Engine struct {
	Throttle float32
	Thrust   float32
}

func (e *Engine) Update(dt float32) {
	e.Throttle = Clampf(e.Throttle, 0.0, 1.0)
}
func (e *Engine) ApplyForce(body *RigidBody) {
	body.AddRelForce(Vec3{e.Thrust * e.Throttle, 0, 0})
}

type Airfoil struct {
	minAlpha     float32
	maxAlpha     float32
	maxLiftCoeff float32
	curves       []Vec3
}

func NewAirfoil(curves []Vec3) *Airfoil {
	a := new(Airfoil)
	a.curves = make([]Vec3, len(curves))
	copy(a.curves, curves)
	a.minAlpha = +10.0
	a.maxAlpha = -10.0
	for i := range curves {
		alphaDeg := a.curves[i].X
		alphaRad := ToRad(alphaDeg)
		a.curves[i].X = alphaRad
		a.minAlpha = Minf(a.minAlpha, alphaRad)
		a.maxAlpha = Maxf(a.maxAlpha, alphaRad)
		liftCoeff := a.curves[i].Y
		a.maxLiftCoeff = Maxf(a.maxLiftCoeff, liftCoeff)
	}
	return a
}

func (a *Airfoil) Sample(alpha float32) (float32, float32) {
	last := len(a.curves) - 1
	alphaStep := (a.maxAlpha - a.minAlpha) / float32(last)
	i := int((alpha - a.minAlpha) / alphaStep)
	j := i + 1
	if i < 0 {
		return a.curves[0].Y, a.curves[0].Z
	}
	if j > last {
		return a.curves[last].Y, a.curves[last].Z
	}
	w := (alpha - a.curves[i].X) / alphaStep
	lerp := a.curves[i].Lerp(a.curves[j], w)
	return lerp.Y, lerp.Z
}

type Wing struct {
	Airfoil     *Airfoil
	CoPres      Vec3
	Normal      Vec3
	aspectRatio float32
	flapRatio   float32
	area        float32
	span        float32
	chord       float32
	deflection  float32
	efficiency  float32
}

func NewWing(airfoil *Airfoil, coPres, normal Vec3, span, chord float32, flapRatio float32) *Wing {
	w := new(Wing)
	w.Airfoil = airfoil
	w.CoPres = coPres
	w.Normal = normal
	w.area = span * chord
	w.span = span
	w.chord = chord
	w.aspectRatio = span * span / w.area
	w.flapRatio = flapRatio
	w.efficiency = 1.0
	return w
}

func (w *Wing) SetDeflection(input float32) {
	w.deflection = Clampf(input, -1, +1)
}
func (w *Wing) ApplyForce(body *RigidBody) {
	localVel := body.VelAt(w.CoPres)
	speed := localVel.Len()
	if speed <= 1 {
		return
	}
	dragDir := localVel.Invert().Norm()
	liftDir := dragDir.Cross(w.Normal).Cross(dragDir).Norm()
	aoa := Asinf(dragDir.Dot(w.Normal))
	liftCoeff, dragCoeff := w.Airfoil.Sample(aoa)
	if w.flapRatio > 0 {
		deltaLiftCoeff := Sqrtf(w.flapRatio) * w.Airfoil.maxLiftCoeff * w.deflection * w.efficiency
		liftCoeff += deltaLiftCoeff
	}
	inducedDragCoeff := liftCoeff * liftCoeff / math.Pi * w.aspectRatio
	dragCoeff += inducedDragCoeff
	airDensity := float32(1.225)
	dynPres := 0.5 * speed * speed * airDensity * w.area
	lift := liftDir.MulSca(liftCoeff * dynPres)
	drag := dragDir.MulSca(dragCoeff * dynPres)
	total := lift.Add(drag)
	body.AddBodyForceAt(total, w.CoPres)
}

type L39 struct {
	Airplane
	meshGeoms   []*Geom1
	meshOffsets []Vec3
}

func NewL39() *L39 {
	mass := float32(2500.0)
	thrust := float32(14000.0)
	wingOffset := float32(+0.3)
	tailOffset := float32(-4.8)
	naca0012 := NewAirfoil(naca0012Curves)
	naca2412 := NewAirfoil(naca2412Curves)
	up := Vec3{0.0, 0.0, +1.0}
	// wingL := Vec3{0.0, -0.03, +0.97}.Norm()
	// wingR := Vec3{0.0, +0.03, +0.97}.Norm()
	wingL := up
	wingR := up
	// TODO: move to NACA64A012
	// TODO: more wing sections
	// TODO: real flaps
	// TODO: extrapolate NACA data for higher AoA
	// TODO: body drag (simple), body lift (oh my)
	right := Vec3{0.0, -1.0, 0.0}
	wings := []Wing{
		*NewWing(naca2412, Vec3{wingOffset, +1.5, -0.4}, wingL, 2.0, 1.7, 0.2), // left inner wing
		*NewWing(naca2412, Vec3{wingOffset, -1.5, -0.4}, wingR, 2.0, 1.7, 0.2), // right inner wing
		*NewWing(naca0012, Vec3{wingOffset, +3.5, -0.3}, wingL, 2.0, 1.3, 0.3), // left outer wing
		*NewWing(naca0012, Vec3{wingOffset, -3.5, -0.3}, wingR, 2.0, 1.3, 0.3), // right outer wing
		*NewWing(naca0012, Vec3{tailOffset, +1.2, +0.3}, up, 1.7, 0.8, 0.5),    // left elevator
		*NewWing(naca0012, Vec3{tailOffset, -1.2, +0.3}, up, 1.7, 0.8, 0.5),    // right elevator tab
		*NewWing(naca0012, Vec3{tailOffset, 0.0, +0.1}, right, 2.0, 3.0, 0.2),  // rudder
	}
	inertia := Mat3{
		+8531.0, -1320.0, 0.0,
		-1320.0, +56608.0, 0.0,
		0.0, 0.0, +11333.0,
	}
	a := new(L39)
	a.Airplane.RigidBody = *NewRigidBody(mass, &inertia)
	a.Airplane.Engine = &Engine{1.0, thrust}
	a.Airplane.Wings = wings

	meshGeoms := []*Geom1{
		geom1DefMap["base/geom/l39_hull"].create(),
		geom1DefMap["base/geom/l39_wing_l"].create(),
		geom1DefMap["base/geom/l39_wing_r"].create(),
		geom1DefMap["base/geom/l39_stab_l"].create(),
		geom1DefMap["base/geom/l39_stab_r"].create(),
		geom1DefMap["base/geom/l39_stab_v"].create(),
		geom1DefMap["base/geom/apu60i"].create(),
		geom1DefMap["base/geom/r60a"].create(),
	}
	meshOffsets := []Vec3{
		{0.0, 0.0, 0.0},
		{-0.931, +2.141, -0.571},
		{-0.931, -2.141, -0.571},
		{-5.325, +1.090, +0.321},
		{-5.325, -1.090, +0.321},
		{-4.765, 0.0, +1.258},
		{-0.020, -2.367, -1.107},
		{-0.020, -2.367, -1.107},
	}
	a.meshGeoms = meshGeoms
	a.meshOffsets = meshOffsets
	return a
}

func (a *L39) Draw() {
	model := a.Rot.ToMat4()
	model.SetPos(a.Pos)
	for i := range a.meshGeoms {
		meshModel := NewMat4Pos(a.meshOffsets[i])
		a.meshGeoms[i].draw(model.Mul(meshModel))
	}
}

var naca0012Curves = []Vec3{
	{-18.50, -1.225, 0.1023}, {-18.25, -1.245, 0.0950}, {-18.00, -1.265, 0.0878},
	{-17.75, -1.285, 0.0808}, {-17.50, -1.303, 0.0742}, {-17.25, -1.319, 0.0681},
	{-17.00, -1.332, 0.0625}, {-16.75, -1.342, 0.0574}, {-16.50, -1.351, 0.0526},
	{-16.25, -1.368, 0.0470}, {-16.00, -1.380, 0.0423}, {-15.75, -1.386, 0.0386},
	{-15.50, -1.387, 0.0357}, {-15.25, -1.385, 0.0334}, {-15.00, -1.380, 0.0316},
	{-14.75, -1.373, 0.0299}, {-14.50, -1.366, 0.0284}, {-14.25, -1.359, 0.0271},
	{-14.00, -1.350, 0.0260}, {-13.75, -1.339, 0.0250}, {-13.50, -1.328, 0.0241},
	{-13.25, -1.314, 0.0234}, {-13.00, -1.313, 0.0220}, {-12.75, -1.300, 0.0211},
	{-12.50, -1.283, 0.0204}, {-12.25, -1.264, 0.0198}, {-12.00, -1.245, 0.0193},
	{-11.75, -1.225, 0.0188}, {-11.50, -1.203, 0.0183}, {-11.25, -1.187, 0.0174},
	{-11.00, -1.166, 0.0169}, {-10.75, -1.144, 0.0164}, {-10.50, -1.122, 0.0160},
	{-10.25, -1.100, 0.0156}, {-10.00, -1.080, 0.0149}, {-9.75, -1.059, 0.0145},
	{-9.50, -1.036, 0.0142}, {-9.25, -1.014, 0.0139}, {-9.00, -0.994, 0.0134},
	{-8.75, -0.973, 0.0130}, {-8.50, -0.951, 0.0127}, {-8.25, -0.931, 0.0124},
	{-8.00, -0.910, 0.0120}, {-7.75, -0.888, 0.0118}, {-7.50, -0.868, 0.0114},
	{-7.25, -0.847, 0.0112}, {-7.00, -0.826, 0.0109}, {-6.75, -0.797, 0.0106},
	{-6.50, -0.763, 0.0103}, {-6.25, -0.728, 0.0100}, {-6.00, -0.693, 0.0097},
	{-5.75, -0.660, 0.0094}, {-5.50, -0.626, 0.0091}, {-5.00, -0.557, 0.0084},
	{-4.75, -0.521, 0.0081}, {-4.50, -0.490, 0.0078}, {-4.25, -0.458, 0.0075},
	{-4.00, -0.427, 0.0072}, {-3.75, -0.400, 0.0070}, {-3.50, -0.373, 0.0068},
	{-3.25, -0.346, 0.0065}, {-3.00, -0.320, 0.0064}, {-2.75, -0.293, 0.0062},
	{-2.50, -0.267, 0.0060}, {-2.25, -0.241, 0.0059}, {-2.00, -0.214, 0.0058},
	{-1.75, -0.187, 0.0057}, {-1.50, -0.161, 0.0056}, {-1.25, -0.134, 0.0055},
	{-1.00, -0.107, 0.0054}, {-0.75, -0.080, 0.0054}, {-0.50, -0.053, 0.0054},
	{-0.25, -0.026, 0.0054}, {0.00, 0.000, 0.0054}, {0.25, 0.026, 0.0054},
	{0.50, 0.053, 0.0054}, {0.75, 0.080, 0.0054}, {1.00, 0.107, 0.0054},
	{1.25, 0.134, 0.0055}, {1.50, 0.161, 0.0056}, {1.75, 0.187, 0.0057},
	{2.00, 0.214, 0.0058}, {2.25, 0.241, 0.0059}, {2.50, 0.267, 0.0060},
	{2.75, 0.293, 0.0062}, {3.00, 0.320, 0.0064}, {3.25, 0.346, 0.0065},
	{3.50, 0.373, 0.0068}, {3.75, 0.400, 0.0070}, {4.00, 0.427, 0.0072},
	{4.25, 0.458, 0.0075}, {4.50, 0.490, 0.0078}, {4.75, 0.521, 0.0081},
	{5.00, 0.557, 0.0084}, {5.50, 0.626, 0.0091}, {5.75, 0.660, 0.0094},
	{6.00, 0.694, 0.0097}, {6.25, 0.729, 0.0100}, {6.50, 0.763, 0.0103},
	{6.75, 0.797, 0.0106}, {7.00, 0.826, 0.0109}, {7.25, 0.847, 0.0112},
	{7.50, 0.868, 0.0114}, {7.75, 0.888, 0.0118}, {8.00, 0.910, 0.0120},
	{8.25, 0.931, 0.0124}, {8.50, 0.951, 0.0127}, {8.75, 0.973, 0.0130},
	{9.00, 0.994, 0.0134}, {9.25, 1.014, 0.0139}, {9.50, 1.036, 0.0142},
	{9.75, 1.059, 0.0145}, {10.00, 1.080, 0.0149}, {10.25, 1.100, 0.0156},
	{10.50, 1.122, 0.0160}, {10.75, 1.145, 0.0164}, {11.00, 1.166, 0.0169},
	{11.25, 1.187, 0.0174}, {11.50, 1.203, 0.0184}, {11.75, 1.225, 0.0188},
	{12.00, 1.245, 0.0193}, {12.25, 1.265, 0.0198}, {12.50, 1.283, 0.0204},
	{12.75, 1.300, 0.0211}, {13.00, 1.313, 0.0220}, {13.25, 1.315, 0.0234},
	{13.50, 1.329, 0.0241}, {13.75, 1.340, 0.0250}, {14.00, 1.350, 0.0260},
	{14.25, 1.360, 0.0271}, {14.50, 1.367, 0.0284}, {14.75, 1.375, 0.0299},
	{15.00, 1.381, 0.0315}, {15.25, 1.386, 0.0334}, {15.50, 1.389, 0.0357},
	{15.75, 1.388, 0.0386}, {16.00, 1.381, 0.0423}, {16.25, 1.370, 0.0470},
	{16.50, 1.353, 0.0526}, {16.75, 1.345, 0.0573}, {17.00, 1.335, 0.0623},
	{17.25, 1.321, 0.0680}, {17.50, 1.305, 0.0741}, {17.75, 1.288, 0.0807},
	{18.00, 1.268, 0.0877}, {18.25, 1.248, 0.0949}, {18.50, 1.228, 0.1022},
}

var naca2412Curves = []Vec3{
	{-17.50, -1.111, 0.0860}, {-17.25, -1.173, 0.0723}, {-17.00, -1.229, 0.0592},
	{-16.75, -1.262, 0.0493}, {-16.50, -1.279, 0.0425}, {-16.25, -1.285, 0.0379},
	{-16.00, -1.286, 0.0345}, {-15.75, -1.285, 0.0320}, {-15.50, -1.281, 0.0301},
	{-15.25, -1.275, 0.0286}, {-15.00, -1.267, 0.0275}, {-14.75, -1.266, 0.0260},
	{-14.50, -1.258, 0.0249}, {-14.25, -1.242, 0.0241}, {-14.00, -1.225, 0.0235},
	{-13.75, -1.206, 0.0230}, {-13.50, -1.188, 0.0225}, {-13.25, -1.169, 0.0221},
	{-13.00, -1.156, 0.0212}, {-12.75, -1.141, 0.0205}, {-12.50, -1.123, 0.0200},
	{-12.25, -1.105, 0.0197}, {-12.00, -1.088, 0.0193}, {-11.75, -1.070, 0.0189},
	{-11.50, -1.051, 0.0186}, {-11.25, -1.032, 0.0184}, {-11.00, -1.022, 0.0173},
	{-10.75, -1.004, 0.0169}, {-10.50, -0.976, 0.0165}, {-10.25, -0.944, 0.0161},
	{-10.00, -0.912, 0.0158}, {-9.75, -0.879, 0.0155}, {-9.50, -0.848, 0.0147},
	{-9.25, -0.817, 0.0141}, {-9.00, -0.786, 0.0137}, {-8.75, -0.754, 0.0133},
	{-8.50, -0.721, 0.0130}, {-8.25, -0.687, 0.0127}, {-8.00, -0.652, 0.0122},
	{-7.75, -0.616, 0.0116}, {-7.50, -0.585, 0.0112}, {-7.25, -0.553, 0.0108},
	{-7.00, -0.520, 0.0105}, {-6.75, -0.488, 0.0103}, {-6.50, -0.462, 0.0099},
	{-6.25, -0.434, 0.0096}, {-6.00, -0.407, 0.0093}, {-5.75, -0.380, 0.0091},
	{-5.50, -0.353, 0.0088}, {-5.25, -0.327, 0.0086}, {-5.00, -0.300, 0.0083},
	{-4.75, -0.273, 0.0081}, {-4.50, -0.246, 0.0079}, {-4.00, -0.191, 0.0076},
	{-3.75, -0.164, 0.0075}, {-3.50, -0.137, 0.0073}, {-3.25, -0.110, 0.0071},
	{-3.00, -0.082, 0.0070}, {-2.75, -0.055, 0.0068}, {-2.50, -0.027, 0.0067},
	{-2.25, -0.000, 0.0066}, {-2.00, 0.027, 0.0065}, {-1.75, 0.054, 0.0064},
	{-1.50, 0.081, 0.0062}, {-1.25, 0.109, 0.0061}, {-1.00, 0.136, 0.0060},
	{-0.75, 0.163, 0.0058}, {-0.50, 0.190, 0.0058}, {-0.25, 0.217, 0.0057},
	{0.00, 0.244, 0.0056}, {0.25, 0.270, 0.0056}, {0.50, 0.296, 0.0055},
	{0.75, 0.321, 0.0054}, {1.00, 0.346, 0.0054}, {1.25, 0.372, 0.0055},
	{1.50, 0.397, 0.0055}, {1.75, 0.425, 0.0056}, {2.00, 0.454, 0.0058},
	{2.75, 0.558, 0.0062}, {3.00, 0.594, 0.0063}, {3.25, 0.631, 0.0065},
	{3.50, 0.668, 0.0067}, {3.75, 0.691, 0.0069}, {4.00, 0.715, 0.0071},
	{4.25, 0.738, 0.0073}, {4.50, 0.762, 0.0075}, {4.75, 0.785, 0.0077},
	{5.00, 0.808, 0.0080}, {5.25, 0.831, 0.0083}, {5.50, 0.855, 0.0086},
	{5.75, 0.878, 0.0090}, {6.00, 0.901, 0.0094}, {6.25, 0.925, 0.0098},
	{6.50, 0.948, 0.0102}, {6.75, 0.971, 0.0107}, {7.00, 0.994, 0.0111},
	{7.25, 1.017, 0.0115}, {7.50, 1.041, 0.0119}, {7.75, 1.064, 0.0123},
	{8.00, 1.088, 0.0127}, {8.25, 1.111, 0.0131}, {8.50, 1.135, 0.0134},
	{8.75, 1.158, 0.0138}, {9.00, 1.180, 0.0143}, {9.25, 1.203, 0.0147},
	{9.50, 1.226, 0.0150}, {9.75, 1.248, 0.0154}, {10.00, 1.269, 0.0159},
	{10.25, 1.288, 0.0165}, {10.50, 1.309, 0.0169}, {10.75, 1.329, 0.0173},
	{11.00, 1.350, 0.0178}, {11.25, 1.368, 0.0182}, {11.50, 1.383, 0.0188},
	{11.75, 1.393, 0.0197}, {12.00, 1.411, 0.0201}, {12.25, 1.428, 0.0206},
	{12.50, 1.444, 0.0212}, {12.75, 1.459, 0.0218}, {13.00, 1.469, 0.0228},
	{13.25, 1.481, 0.0237}, {13.50, 1.496, 0.0244}, {13.75, 1.510, 0.0252},
	{14.00, 1.522, 0.0262}, {14.25, 1.531, 0.0275}, {14.50, 1.538, 0.0289},
	{14.75, 1.549, 0.0300}, {15.00, 1.559, 0.0314}, {15.25, 1.567, 0.0329},
	{15.50, 1.571, 0.0349}, {15.75, 1.572, 0.0373}, {16.00, 1.577, 0.0393},
	{16.25, 1.580, 0.0416}, {16.50, 1.582, 0.0442}, {16.75, 1.581, 0.0471},
	{17.00, 1.578, 0.0504}, {17.25, 1.571, 0.0542}, {17.50, 1.560, 0.0589},
	{17.75, 1.548, 0.0638}, {18.00, 1.541, 0.0680}, {18.25, 1.532, 0.0726},
	{18.50, 1.521, 0.0777}, {18.75, 1.508, 0.0832}, {19.00, 1.494, 0.0889},
	{19.25, 1.478, 0.0950},
}
