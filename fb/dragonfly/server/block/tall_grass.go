package block

import (
	"phoenixbuilder_3rd_gui/fb/dragonfly/server/block/cube"
	"phoenixbuilder_3rd_gui/fb/dragonfly/server/item"
	"phoenixbuilder_3rd_gui/fb/dragonfly/server/item/tool"
	"phoenixbuilder_3rd_gui/fb/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"math/rand"
)

// TallGrass is a transparent plant block which can be used to obtain seeds and as decoration.
type TallGrass struct {
	replaceable
	transparent
	empty

	// Type is the type of grass that the plant represents.
	Type GrassType
}

// FlammabilityInfo ...
func (g TallGrass) FlammabilityInfo() FlammabilityInfo {
	return newFlammabilityInfo(30, 100, false)
}

// BreakInfo ...
func (g TallGrass) BreakInfo() BreakInfo {
	return newBreakInfo(0, alwaysHarvestable, nothingEffective, func(t tool.Tool, enchantments []item.Enchantment) []item.Stack {
		if t.ToolType() == tool.TypeShears || hasSilkTouch(enchantments) {
			return []item.Stack{item.NewStack(g, 1)}
		}
		if rand.Float32() > 0.57 {
			return []item.Stack{item.NewStack(WheatSeeds{}, 1)}
		}
		return []item.Stack{}
	})
}

// BoneMeal attempts to affect the block using a bone meal item.
func (g TallGrass) BoneMeal(pos cube.Pos, w *world.World) bool {
	upper := DoubleTallGrass{Type: g.Type, UpperPart: true}
	if replaceableWith(w, pos.Side(cube.FaceUp), upper) {
		w.SetBlock(pos, DoubleTallGrass{Type: g.Type})
		w.SetBlock(pos.Side(cube.FaceUp), upper)
		return true
	}
	return false
}

// NeighbourUpdateTick ...
func (g TallGrass) NeighbourUpdateTick(pos, _ cube.Pos, w *world.World) {
	if !supportsVegetation(g, w.Block(pos.Side(cube.FaceDown))) {
		w.BreakBlock(pos)
	}
}

// HasLiquidDrops ...
func (g TallGrass) HasLiquidDrops() bool {
	return true
}

// UseOnBlock ...
func (g TallGrass) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, w *world.World, user item.User, ctx *item.UseContext) bool {
	pos, _, used := firstReplaceable(w, pos, face, g)
	if !used {
		return false
	}
	if !supportsVegetation(g, w.Block(pos.Side(cube.FaceDown))) {
		return false
	}

	place(w, pos, g, user, ctx)
	return placed(ctx)
}

// EncodeItem ...
func (g TallGrass) EncodeItem() (name string, meta int16) {
	return "minecraft:tallgrass", int16(g.Type.Uint8() + 1)
}

// EncodeBlock ...
func (g TallGrass) EncodeBlock() (name string, properties map[string]interface{}) {
	switch g.Type {
	case NormalGrass():
		return "minecraft:tallgrass", map[string]interface{}{"tall_grass_type": "tall"}
	case Fern():
		return "minecraft:tallgrass", map[string]interface{}{"tall_grass_type": "fern"}
	}
	panic("should never happen")
}

// allTallGrass ...
func allTallGrass() (grasses []world.Block) {
	for _, g := range GrassTypes() {
		grasses = append(grasses, TallGrass{Type: g})
	}
	return
}
