package block

import (
	"phoenixbuilder_3rd_gui/fb/dragonfly/server/block/instrument"
	"phoenixbuilder_3rd_gui/fb/dragonfly/server/item"
	"math/rand"
)

// Glowstone is commonly found on the ceiling of the nether dimension.
type Glowstone struct {
	solid
}

// Instrument ...
func (g Glowstone) Instrument() instrument.Instrument {
	return instrument.Pling()
}

// BreakInfo ...
func (g Glowstone) BreakInfo() BreakInfo {
	return newBreakInfo(0.3, alwaysHarvestable, nothingEffective, silkTouchDrop(item.NewStack(item.GlowstoneDust{}, rand.Intn(3)+2), item.NewStack(g, 1)))
}

// EncodeItem ...
func (Glowstone) EncodeItem() (name string, meta int16) {
	return "minecraft:glowstone", 0
}

// EncodeBlock ...
func (Glowstone) EncodeBlock() (string, map[string]interface{}) {
	return "minecraft:glowstone", nil
}

// LightEmissionLevel returns 15.
func (Glowstone) LightEmissionLevel() uint8 {
	return 15
}
