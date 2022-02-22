package block

import "phoenixbuilder/dragonfly/server/item"

// NetherQuartzOre is ore found in the Nether.
type NetherQuartzOre struct {
	solid
	bassDrum
}

// BreakInfo ...
func (q NetherQuartzOre) BreakInfo() BreakInfo {
	return BreakInfo{
		Hardness:    3,
		Harvestable: pickaxeHarvestable,
		Effective:   pickaxeEffective,
		Drops:       silkTouchOneOf(item.NetherQuartz{}, q),
		XPDrops:     XPDropRange{0, 3},
	}
}

// EncodeItem ...
func (NetherQuartzOre) EncodeItem() (name string, meta int16) {
	return "minecraft:quartz_ore", 0
}

// EncodeBlock ...
func (NetherQuartzOre) EncodeBlock() (string, map[string]interface{}) {
	return "minecraft:quartz_ore", nil
}
