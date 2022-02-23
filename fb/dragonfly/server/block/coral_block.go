package block

import (
	"phoenixbuilder_3rd_gui/fb/dragonfly/server/block/cube"
	"phoenixbuilder_3rd_gui/fb/dragonfly/server/world"
	"math/rand"
	"time"
)

// CoralBlock is a solid block that comes in 5 variants.
type CoralBlock struct {
	solid
	bassDrum

	// Type is the type of coral of the block.
	Type CoralType
	// Dead is whether the coral block is dead.
	Dead bool
}

// NeighbourUpdateTick ...
func (c CoralBlock) NeighbourUpdateTick(pos, changedNeighbour cube.Pos, w *world.World) {
	if c.Dead {
		return
	}
	w.ScheduleBlockUpdate(pos, time.Second*5/2)
}

// ScheduledTick ...
func (c CoralBlock) ScheduledTick(pos cube.Pos, w *world.World, _ *rand.Rand) {
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
func (c CoralBlock) BreakInfo() BreakInfo {
	return newBreakInfo(7, pickaxeHarvestable, pickaxeEffective, silkTouchOneOf(CoralBlock{Type: c.Type, Dead: true}, c))
}

// EncodeBlock ...
func (c CoralBlock) EncodeBlock() (name string, properties map[string]interface{}) {
	return "minecraft:coral_block", map[string]interface{}{"coral_color": c.Type.Colour().String(), "dead_bit": c.Dead}
}

// EncodeItem ...
func (c CoralBlock) EncodeItem() (name string, meta int16) {
	if c.Dead {
		return "minecraft:coral_block", int16(c.Type.Uint8() | 8)
	}
	return "minecraft:coral_block", int16(c.Type.Uint8())
}

// allCoralBlocks returns a list of all coral block variants
func allCoralBlocks() (c []world.Block) {
	f := func(dead bool) {
		for _, t := range CoralTypes() {
			c = append(c, CoralBlock{Type: t, Dead: dead})
		}
	}
	f(true)
	f(false)
	return
}
