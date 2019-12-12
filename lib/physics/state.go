package physics

import "github.com/go-gl/mathgl/mgl32"

// https://gafferongames.com/post/physics_in_3d/

type State struct {

	// primary
	momentum        mgl32.Vec3
	angularMomentum mgl32.Vec3
	position        mgl32.Vec3
	orientation     mgl32.Quat

	// secondary
	spin            mgl32.Quat
	angularVelocity mgl32.Vec3
	velocity        mgl32.Vec3

	// constant
	//inertia        float32
	inverseInertia float32
	//mass           float32
	inverseMass float32
}

func (s *State) SetMass(m float32) {
	s.inverseMass = 1 / m
}

func (s *State) recalculate() {
	s.velocity = s.momentum.Mul(s.inverseMass)
	s.angularVelocity = s.angularMomentum.Mul(s.inverseInertia)
	s.orientation.Normalize()

	q := mgl32.Quat{
		W: 0,
		V: s.angularVelocity,
	}

	//s.spin = 0.5 * q. * s.orientation
	s.spin = q.Mul(s.orientation).Scale(0.5)
}

/**
dp/dt = F
v = p/m
dx/dt = v
*/

type Derivative struct {
	force    mgl32.Vec3
	velocity mgl32.Vec3
	spin     mgl32.Quat
	torque   mgl32.Quat
}

// This function returns an acceleration torque to induce a spin around the x axis,
// but also applies a damping over time so that at a certain speed the accelerating
//and damping will cancel each other out. This is done so that the rotation will
//reach a certain rate and stay constant instead of getting faster and faster over time.
func torque(state *State, t float32) mgl32.Vec3 {
	//return mgl32.Vec3{1,0,0}.Sub() - state.angularVelocity * 0.1f;
	return mgl32.Vec3{1, 0, 0}.Sub(state.angularVelocity).Mul(0.1)
}
