package block

import (
	"phoenixbuilder_3rd_gui/fb/dragonfly/server/block/cube"
	"phoenixbuilder_3rd_gui/fb/dragonfly/server/block/model"
	"phoenixbuilder_3rd_gui/fb/dragonfly/server/world"
)

// NetherBrickFence is the nether brick variant of the fence block.
type NetherBrickFence struct {
	transparent
}

// BreakInfo ...
func (n NetherBrickFence) BreakInfo() BreakInfo {
	return newBreakInfo(2, pickaxeHarvestable, pickaxeEffective, oneOf(n))
}

// CanDisplace ...
func (NetherBrickFence) CanDisplace(b world.Liquid) bool {
	_, ok := b.(Water)
	return ok
}

// SideClosed ...
func (NetherBrickFence) SideClosed(cube.Pos, cube.Pos, *world.World) bool {
	return false
}

// Model ...
func (n NetherBrickFence) Model() world.BlockModel {
	return model.Fence{}
}

// EncodeItem ...
func (NetherBrickFence) EncodeItem() (name string, meta int16) {
	return "minecraft:nether_brick_fence", 0
}

// EncodeBlock ...
func (NetherBrickFence) EncodeBlock() (string, map[string]interface{}) {
	return "minecraft:nether_brick_fence", nil
}
