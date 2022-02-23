package block

import (
	"phoenixbuilder_3rd_gui/fb/dragonfly/server/item"
	"phoenixbuilder_3rd_gui/fb/dragonfly/server/item/tool"
)

// CopperOre is a rare mineral block found underground.
type CopperOre struct {
	solid
	bassDrum

	// Type is the type of copper ore.
	Type OreType
}

// BreakInfo ...
func (c CopperOre) BreakInfo() BreakInfo {
	return newBreakInfo(c.Type.Hardness(), func(t tool.Tool) bool {
		return t.ToolType() == tool.TypePickaxe && t.HarvestLevel() >= tool.TierStone.HarvestLevel
	}, pickaxeEffective, silkTouchOneOf(item.RawCopper{}, c))
}

// EncodeItem ...
func (c CopperOre) EncodeItem() (name string, meta int16) {
	return "minecraft:" + c.Type.Prefix() + "copper_ore", 0
}

// EncodeBlock ...
func (c CopperOre) EncodeBlock() (string, map[string]interface{}) {
	return "minecraft:" + c.Type.Prefix() + "copper_ore", nil
}
