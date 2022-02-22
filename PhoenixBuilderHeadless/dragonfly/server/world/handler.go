package world

import (
	"phoenixbuilder/dragonfly/server/block/cube"
	"phoenixbuilder/dragonfly/server/event"
)

// Handler handles events that are called by a world. Implementations of Handler may be used to listen to
// specific events such as when an entity is added to the world.
type Handler interface {
	// HandleLiquidFlow handles the flowing of a liquid from one block position from into another block
	// position into. The liquid that will replace the block replaced is also passed.
	HandleLiquidFlow(ctx *event.Context, from, into cube.Pos, liquid, replaced Block)
	// HandleLiquidHarden handles the hardening of a liquid at hardenedPos. The liquid that was hardened,
	// liquidHardened, and the liquid that caused it to harden, otherLiquid, are passed. The block created
	// as a result is also passed.
	HandleLiquidHarden(ctx *event.Context, hardenedPos cube.Pos, liquidHardened, otherLiquid, newBlock Block)
}

// NopHandler implements the Handler interface but does not execute any code when an event is called. The
// default Handler of worlds is set to NopHandler.
// Users may embed NopHandler to avoid having to implement each method.
type NopHandler struct{}

// Compile time check to make sure NopHandler implements Handler.
var _ Handler = (*NopHandler)(nil)

// HandleLiquidFlow ...
func (NopHandler) HandleLiquidFlow(*event.Context, cube.Pos, cube.Pos, Block, Block) {}

// HandleLiquidHarden ...
func (NopHandler) HandleLiquidHarden(*event.Context, cube.Pos, Block, Block, Block) {}
