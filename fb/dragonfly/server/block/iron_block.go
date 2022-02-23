package block

import (
	"phoenixbuilder_3rd_gui/fb/dragonfly/server/block/instrument"
	"phoenixbuilder_3rd_gui/fb/dragonfly/server/item/tool"
)

// IronBlock is a precious metal block made from 9 iron ingots.
type IronBlock struct {
	solid
}

// Instrument ...
func (i IronBlock) Instrument() instrument.Instrument {
	return instrument.IronXylophone()
}

// BreakInfo ...
func (i IronBlock) BreakInfo() BreakInfo {
	return newBreakInfo(5, func(t tool.Tool) bool {
		return t.ToolType() == tool.TypePickaxe && t.HarvestLevel() >= tool.TierStone.HarvestLevel
	}, pickaxeEffective, oneOf(i))
}

// PowersBeacon ...
func (IronBlock) PowersBeacon() bool {
	return true
}

// EncodeItem ...
func (IronBlock) EncodeItem() (name string, meta int16) {
	return "minecraft:iron_block", 0
}

// EncodeBlock ...
func (IronBlock) EncodeBlock() (string, map[string]interface{}) {
	return "minecraft:iron_block", nil
}
