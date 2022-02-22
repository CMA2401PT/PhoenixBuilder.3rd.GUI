package block

import (
	"phoenixbuilder/dragonfly/server/block/instrument"
	"phoenixbuilder/dragonfly/server/item/tool"
)

// EmeraldBlock is a precious mineral block crafted using 9 emeralds.
type EmeraldBlock struct {
	solid
}

// Instrument ...
func (e EmeraldBlock) Instrument() instrument.Instrument {
	return instrument.Bit()
}

// BreakInfo ...
func (e EmeraldBlock) BreakInfo() BreakInfo {
	return newBreakInfo(5, func(t tool.Tool) bool {
		return t.ToolType() == tool.TypePickaxe && t.HarvestLevel() >= tool.TierIron.HarvestLevel
	}, pickaxeEffective, oneOf(e))
}

// PowersBeacon ...
func (EmeraldBlock) PowersBeacon() bool {
	return true
}

// EncodeItem ...
func (EmeraldBlock) EncodeItem() (name string, meta int16) {
	return "minecraft:emerald_block", 0
}

// EncodeBlock ...
func (EmeraldBlock) EncodeBlock() (string, map[string]interface{}) {
	return "minecraft:emerald_block", nil
}
