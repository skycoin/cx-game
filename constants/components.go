package constants

type AiHandlerID int
type DrawHandlerID int

const (
	AI_HANDLER_NULL AiHandlerID = iota

	NUM_AI_HANDLERS
)

const (
	DRAW_HANDLER_NULL DrawHandlerID = iota
	DRAW_HANDLER_QUAD

	NUM_DRAW_HANDLERS
)
