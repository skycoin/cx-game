package particles

type ParticleList struct {
	Particles []Particle
}

func NewParticleList() ParticleList {
	return ParticleList{}
}
