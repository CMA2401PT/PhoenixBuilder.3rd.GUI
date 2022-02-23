package block

import (
	"phoenixbuilder_3rd_gui/fb/dragonfly/server/item/tool"
)

// NetheriteBlock is a precious mineral block made from 9 netherite ingots.
type NetheriteBlock struct {
	solid
	bassDrum
}

// BreakInfo ...
func (n NetheriteBlock) BreakInfo() BreakInfo {
	return newBreakInfo(50, func(t tool.Tool) bool {
		return t.ToolType() == tool.TypePickaxe && t.HarvestLevel() >= tool.TierDiamond.HarvestLevel
	}, pickaxeEffective, oneOf(n))
}

// PowersBeacon ...
func (NetheriteBlock) PowersBeacon() bool {
	return true
}

// EncodeItem ...
func (NetheriteBlock) EncodeItem() (name string, meta int16) {
	return "minecraft:netherite_block", 0
}

// EncodeBlock ...
func (NetheriteBlock) EncodeBlock() (string, map[string]interface{}) {
	return "minecraft:netherite_block", nil
}
