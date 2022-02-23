package block

import (
	"phoenixbuilder_3rd_gui/fb/dragonfly/server/block/cube"
	"phoenixbuilder_3rd_gui/fb/dragonfly/server/item"
	"phoenixbuilder_3rd_gui/fb/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
)

// SporeBlossom is a decorative block.
type SporeBlossom struct {
	empty
	transparent
}

// HasLiquidDrops ...
func (s SporeBlossom) HasLiquidDrops() bool {
	return true
}

// NeighbourUpdateTick ...
func (s SporeBlossom) NeighbourUpdateTick(pos, _ cube.Pos, w *world.World) {
	if !w.Block(pos.Side(cube.FaceUp)).Model().FaceSolid(pos.Side(cube.FaceUp), cube.FaceDown, w) {
		w.BreakBlock(pos)
	}
}

// UseOnBlock ...
func (s SporeBlossom) UseOnBlock(pos cube.Pos, face cube.Face, clickPos mgl64.Vec3, w *world.World, user item.User, ctx *item.UseContext) (used bool) {
	pos, _, used = firstReplaceable(w, pos, face, s)
	if !used {
		return
	}
	if !w.Block(pos.Side(cube.FaceUp)).Model().FaceSolid(pos.Side(cube.FaceUp), cube.FaceDown, w) {
		return
	}

	place(w, pos, s, user, ctx)
	return placed(ctx)
}

// BreakInfo ...
func (s SporeBlossom) BreakInfo() BreakInfo {
	return newBreakInfo(0, alwaysHarvestable, nothingEffective, oneOf(s))
}

// EncodeItem ...
func (s SporeBlossom) EncodeItem() (name string, meta int16) {
	return "minecraft:spore_blossom", 0
}

// EncodeBlock ...
func (s SporeBlossom) EncodeBlock() (string, map[string]interface{}) {
	return "minecraft:spore_blossom", nil
}
