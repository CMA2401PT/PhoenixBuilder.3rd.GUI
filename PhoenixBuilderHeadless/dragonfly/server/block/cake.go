package block

import (
	"phoenixbuilder/dragonfly/server/block/cube"
	"phoenixbuilder/dragonfly/server/block/model"
	"phoenixbuilder/dragonfly/server/item"
	"phoenixbuilder/dragonfly/server/world"
	"phoenixbuilder/dragonfly/server/world/sound"
	"github.com/go-gl/mathgl/mgl64"
)

// Cake is an edible block.
type Cake struct {
	transparent

	// Bites is the amount of bites taken out of the cake.
	Bites int
}

// CanDisplace ...
func (c Cake) CanDisplace(b world.Liquid) bool {
	_, water := b.(Water)
	return water
}

// SideClosed ...
func (c Cake) SideClosed(cube.Pos, cube.Pos, *world.World) bool {
	return false
}

// UseOnBlock ...
func (c Cake) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, w *world.World, user item.User, ctx *item.UseContext) bool {
	pos, _, used := firstReplaceable(w, pos, face, c)
	if !used {
		return false
	}

	if _, air := w.Block(pos.Side(cube.FaceDown)).(Air); air {
		return false
	}

	place(w, pos, c, user, ctx)
	return placed(ctx)
}

// NeighbourUpdateTick ...
func (c Cake) NeighbourUpdateTick(pos, _ cube.Pos, w *world.World) {
	if _, air := w.Block(pos.Side(cube.FaceDown)).(Air); air {
		w.BreakBlock(pos)
	}
}

// Activate ...
func (c Cake) Activate(pos cube.Pos, _ cube.Face, w *world.World, u item.User) {
	if i, ok := u.(interface {
		Saturate(food int, saturation float64)
	}); ok {
		i.Saturate(2, 0.4)
		w.PlaySound(u.Position().Add(mgl64.Vec3{0, 1.5}), sound.Burp{})
		c.Bites++
		if c.Bites > 6 {
			w.BreakBlockWithoutParticles(pos)
			return
		}
		w.PlaceBlock(pos, c)
	}
}

// BreakInfo ...
func (c Cake) BreakInfo() BreakInfo {
	return newBreakInfo(0.5, neverHarvestable, nothingEffective, simpleDrops())
}

// EncodeItem ...
func (c Cake) EncodeItem() (name string, meta int16) {
	return "minecraft:cake", 0
}

// EncodeBlock ...
func (c Cake) EncodeBlock() (name string, properties map[string]interface{}) {
	return "minecraft:cake", map[string]interface{}{"bite_counter": int32(c.Bites)}
}

// Model ...
func (c Cake) Model() world.BlockModel {
	return model.Cake{Bites: c.Bites}
}

// allCake ...
func allCake() (cake []world.Block) {
	for bites := 0; bites < 7; bites++ {
		cake = append(cake, Cake{Bites: bites})
	}
	return
}
