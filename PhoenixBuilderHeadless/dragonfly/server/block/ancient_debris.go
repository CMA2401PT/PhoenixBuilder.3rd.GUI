package block

import (
	"phoenixbuilder/dragonfly/server/item/tool"
)

// AncientDebris is a rare ore found within The Nether.
type AncientDebris struct {
	solid
}

// BreakInfo ...
func (a AncientDebris) BreakInfo() BreakInfo {
	return newBreakInfo(30, func(t tool.Tool) bool {
		return t.ToolType() == tool.TypePickaxe && t.HarvestLevel() >= tool.TierDiamond.HarvestLevel
	}, pickaxeEffective, oneOf(a))
}

// EncodeItem ...
func (AncientDebris) EncodeItem() (name string, meta int16) {
	return "minecraft:ancient_debris", 0
}

// EncodeBlock ...
func (AncientDebris) EncodeBlock() (string, map[string]interface{}) {
	return "minecraft:ancient_debris", nil
}
