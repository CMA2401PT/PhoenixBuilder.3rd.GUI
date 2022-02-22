package block

import (
	"phoenixbuilder/dragonfly/server/block/cube"
	"phoenixbuilder/dragonfly/server/item"
	"phoenixbuilder/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"math/rand"
	"time"
)

// Coral is a non solid block that comes in 5 variants.
type Coral struct {
	empty
	transparent
	bassDrum

	// Type is the type of coral of the block.
	Type CoralType
	// Dead is whether the coral is dead.
	Dead bool
}

// UseOnBlock ...
func (c Coral) UseOnBlock(pos cube.Pos, face cube.Face, clickPos mgl64.Vec3, w *world.World, user item.User, ctx *item.UseContext) bool {
	pos, _, used := firstReplaceable(w, pos, face, c)
	if !used {
		return false
	}
	if !w.Block(pos.Side(cube.FaceDown)).Model().FaceSolid(pos.Side(cube.FaceDown), cube.FaceUp, w) {
		return false
	}
	if liquid, ok := w.Liquid(pos); ok {
		if water, ok := liquid.(Water); ok {
			if water.Depth != 8 {
				return false
			}
		}
	}

	place(w, pos, c, user, ctx)
	return placed(ctx)
}

// HasLiquidDrops ...
func (c Coral) HasLiquidDrops() bool {
	return false
}

// CanDisplace ...
func (c Coral) CanDisplace(b world.Liquid) bool {
	_, water := b.(Water)
	return water
}

// SideClosed ...
func (c Coral) SideClosed(pos, side cube.Pos, w *world.World) bool {
	return false
}

// NeighbourUpdateTick ...
func (c Coral) NeighbourUpdateTick(pos, changedNeighbour cube.Pos, w *world.World) {
	if !w.Block(pos.Side(cube.FaceDown)).Model().FaceSolid(pos.Side(cube.FaceDown), cube.FaceUp, w) {
		w.BreakBlock(pos)
		return
	}
	if c.Dead {
		return
	}
	w.ScheduleBlockUpdate(pos, time.Second*5/2)
}

// ScheduledTick ...
func (c Coral) ScheduledTick(pos cube.Pos, w *world.World, _ *rand.Rand) {
	if c.Dead {
		return
	}

	adjacentWater := false
	pos.Neighbours(func(neighbour cube.Pos) {
		if liquid, ok := w.Liquid(neighbour); ok {
			if _, ok := liquid.(Water); ok {
				adjacentWater = true
			}
		}
	})
	if !adjacentWater {
		c.Dead = true
		w.PlaceBlock(pos, c)
	}
}

// BreakInfo ...
func (c Coral) BreakInfo() BreakInfo {
	return newBreakInfo(0, alwaysHarvestable, nothingEffective, silkTouchOnlyDrop(c))
}

// EncodeBlock ...
func (c Coral) EncodeBlock() (name string, properties map[string]interface{}) {
	return "minecraft:coral", map[string]interface{}{"coral_color": c.Type.Colour().String(), "dead_bit": c.Dead}
}

// EncodeItem ...
func (c Coral) EncodeItem() (name string, meta int16) {
	if c.Dead {
		return "minecraft:coral", int16(c.Type.Uint8() | 8)
	}
	return "minecraft:coral", int16(c.Type.Uint8())
}

// allCoral returns a list of all coral block variants
func allCoral() (c []world.Block) {
	f := func(dead bool) {
		for _, t := range CoralTypes() {
			c = append(c, Coral{Type: t, Dead: dead})
		}
	}
	f(true)
	f(false)
	return
}
