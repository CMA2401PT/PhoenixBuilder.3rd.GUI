package block

import (
	"phoenixbuilder/dragonfly/server/item/tool"
)

// Obsidian is a dark purple block known for its high blast resistance and strength, most commonly found when
// water flows over lava.
type Obsidian struct {
	solid
	bassDrum
	// Crying specifies if the block is a crying obsidian block. If true, the block is blue and emits light.
	Crying bool
}

// LightEmissionLevel ...
func (o Obsidian) LightEmissionLevel() uint8 {
	if o.Crying {
		return 10
	}
	return 0
}

// EncodeItem ...
func (o Obsidian) EncodeItem() (name string, meta int16) {
	if o.Crying {
		return "minecraft:crying_obsidian", 0
	}
	return "minecraft:obsidian", 0
}

// EncodeBlock ...
func (o Obsidian) EncodeBlock() (string, map[string]interface{}) {
	if o.Crying {
		return "minecraft:crying_obsidian", nil
	}
	return "minecraft:obsidian", nil
}

// BreakInfo ...
func (o Obsidian) BreakInfo() BreakInfo {
	return newBreakInfo(50, func(t tool.Tool) bool {
		return t.ToolType() == tool.TypePickaxe && t.HarvestLevel() >= tool.TierDiamond.HarvestLevel
	}, pickaxeEffective, oneOf(o))
}
