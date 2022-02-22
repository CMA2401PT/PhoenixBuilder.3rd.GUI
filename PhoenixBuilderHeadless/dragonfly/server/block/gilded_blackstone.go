package block

import (
	"phoenixbuilder/dragonfly/server/item"
	"phoenixbuilder/dragonfly/server/item/tool"
	"math/rand"
)

// GildedBlackstone is a variant of blackstone that can drop itself or gold nuggets when mined.
type GildedBlackstone struct {
	solid
}

// BreakInfo ...
func (b GildedBlackstone) BreakInfo() BreakInfo {
	return newBreakInfo(1.5, pickaxeHarvestable, pickaxeEffective, func(tool.Tool, []item.Enchantment) []item.Stack {
		if rand.Float64() < 0.1 {
			return []item.Stack{item.NewStack(item.GoldNugget{}, rand.Intn(4)+2)}
		}
		return []item.Stack{item.NewStack(b, 1)}
	})
}

// EncodeItem ...
func (GildedBlackstone) EncodeItem() (name string, meta int16) {
	return "minecraft:gilded_blackstone", 0
}

// EncodeBlock ...
func (GildedBlackstone) EncodeBlock() (string, map[string]interface{}) {
	return "minecraft:gilded_blackstone", nil
}
