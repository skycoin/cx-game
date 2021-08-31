package constants

const (
	DRAW_SPRITE_BATCH_SIZE int32 = 100
	DRAW_COLOR_BATCH_SIZE  int32 = 100
	PIXELS_PER_TILE        int   = 32
)

const (
	BACKGROUND_Z float32 = -100
	//2_gap between layers so we can draw something between them
	BGLAYER_Z    float32 = -12
	PIPELAYER_Z  float32 = -11
	MIDLAYER_Z   float32 = -10
	FRONTLAYER_Z float32 = -8
	AGENT_Z      float32 = -5
	PLAYER_Z     float32 = -3
	HUD_Z        float32 = 0
)
