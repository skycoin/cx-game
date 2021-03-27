package particles

type Particle struct {
	ID               int32
	ParticleMetaType int32
	Type             int32
	Sprite           int32
	Size             int32
	PosX             int32
	PosY             int32
	VelocityX        float32
	VelocityY        float32
	TimeToLive       int32
}
