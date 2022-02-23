package block

import (
	"phoenixbuilder_3rd_gui/fb/dragonfly/server/block/cube"
	"phoenixbuilder_3rd_gui/fb/dragonfly/server/item"
	"phoenixbuilder_3rd_gui/fb/dragonfly/server/item/tool"
	"phoenixbuilder_3rd_gui/fb/dragonfly/server/world"
	"math/rand"
)

// Gravel is a block affected by gravity. It has a 10% chance of dropping flint instead of itself on break.
type Gravel struct {
	gravityAffected
	solid
	snare
}

// NeighbourUpdateTick ...
func (g Gravel) NeighbourUpdateTick(pos, _ cube.Pos, w *world.World) {
	g.fall(g, pos, w)
}

// BreakInfo ...
func (g Gravel) BreakInfo() BreakInfo {
	return newBreakInfo(0.6, alwaysHarvestable, shovelEffective, func(t tool.Tool, enchantments []item.Enchantment) []item.Stack {
		if !hasSilkTouch(enchantments) && rand.Float64() < 0.1 {
			return []item.Stack{item.NewStack(item.Flint{}, 1)}
		}
		return []item.Stack{item.NewStack(g, 1)}
	})
}

// EncodeItem ...
func (Gravel) EncodeItem() (name string, meta int16) {
	return "minecraft:gravel", 0
}

// EncodeBlock ...
func (Gravel) EncodeBlock() (string, map[string]interface{}) {
	return "minecraft:gravel", nil
}
