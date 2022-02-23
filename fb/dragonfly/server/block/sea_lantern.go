package block

import (
	"phoenixbuilder_3rd_gui/fb/dragonfly/server/item"
	"math/rand"
)

// SeaLantern is an underwater light sources that appear in ocean monuments and underwater ruins.
type SeaLantern struct {
	transparent
	solid
	clicksAndSticks
}

// LightEmissionLevel ...
func (SeaLantern) LightEmissionLevel() uint8 {
	return 15
}

// BreakInfo ...
func (s SeaLantern) BreakInfo() BreakInfo {
	return newBreakInfo(0.3, alwaysHarvestable, nothingEffective, silkTouchDrop(item.NewStack(item.PrismarineCrystals{}, rand.Intn(2)+2), item.NewStack(s, 1)))
}

// EncodeItem ...
func (SeaLantern) EncodeItem() (name string, meta int16) {
	return "minecraft:sealantern", 0
}

// EncodeBlock ...
func (SeaLantern) EncodeBlock() (string, map[string]interface{}) {
	return "minecraft:seaLantern", nil
}
