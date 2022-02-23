package block

import (
	"phoenixbuilder_3rd_gui/fb/dragonfly/server/item/tool"
)

// CoalBlock is a precious mineral block made from 9 coal.
type CoalBlock struct {
	solid
	bassDrum
}

// FlammabilityInfo ...
func (c CoalBlock) FlammabilityInfo() FlammabilityInfo {
	return newFlammabilityInfo(5, 5, false)
}

// BreakInfo ...
func (c CoalBlock) BreakInfo() BreakInfo {
	return newBreakInfo(5, func(t tool.Tool) bool {
		return t.ToolType() == tool.TypePickaxe && t.HarvestLevel() >= tool.TierWood.HarvestLevel
	}, pickaxeEffective, oneOf(c))
}

// EncodeItem ...
func (CoalBlock) EncodeItem() (name string, meta int16) {
	return "minecraft:coal_block", 0
}

// EncodeBlock ...
func (CoalBlock) EncodeBlock() (name string, properties map[string]interface{}) {
	return "minecraft:coal_block", nil
}
