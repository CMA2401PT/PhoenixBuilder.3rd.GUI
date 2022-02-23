package block

import (
	"phoenixbuilder_3rd_gui/fb/dragonfly/server/item"
	"phoenixbuilder_3rd_gui/fb/dragonfly/server/world"
)

// Concrete is a solid block which comes in the 16 regular dye colors, created by placing concrete powder
// adjacent to water.
type Concrete struct {
	solid
	bassDrum

	// Colour is the colour of the concrete block.
	Colour item.Colour
}

// BreakInfo ...
func (c Concrete) BreakInfo() BreakInfo {
	return newBreakInfo(1.8, pickaxeHarvestable, pickaxeEffective, oneOf(c))
}

// EncodeItem ...
func (c Concrete) EncodeItem() (name string, meta int16) {
	return "minecraft:concrete", int16(c.Colour.Uint8())
}

// EncodeBlock ...
func (c Concrete) EncodeBlock() (name string, properties map[string]interface{}) {
	return "minecraft:concrete", map[string]interface{}{"color": c.Colour.String()}
}

// allConcrete returns concrete blocks with all possible colours.
func allConcrete() []world.Block {
	b := make([]world.Block, 0, 16)
	for _, c := range item.Colours() {
		b = append(b, Concrete{Colour: c})
	}
	return b
}
