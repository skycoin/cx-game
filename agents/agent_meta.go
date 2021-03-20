package agents

type AgentMeta struct {
	MovementSpeed    float32
	PlayerControlled bool
	AgentType        int
	AiControllerId   int
	PosX             float32
	PosY             float32
	MomentumX        float32
	MomentumY        float32
	JumpHeight       float32
	// StaticFriction float32
	DynamicFriction float32
}
