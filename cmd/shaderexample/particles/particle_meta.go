package particles

type ParticleMeta struct {
	MetaID                 int
	GravityEffect          bool
	Bounces                bool
	CollisionWithTerrain   bool
	CollsionWithPlayer     bool
	ParticleDamagesPlayer  bool
	ParticleDamagesTerrain bool
}
