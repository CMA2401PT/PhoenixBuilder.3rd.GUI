package block

import (
	"phoenixbuilder/dragonfly/server/block/cube"
	"phoenixbuilder/dragonfly/server/block/model"
	"phoenixbuilder/dragonfly/server/item"
	"phoenixbuilder/dragonfly/server/item/tool"
	"phoenixbuilder/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"math/rand"
)

// CocoaBean is a crop block found in jungle biomes.
type CocoaBean struct {
	transparent

	// Facing is the direction from the cocoa bean to the log.
	Facing cube.Direction
	// Age is the stage of the cocoa bean's growth. 2 is fully grown.
	Age int
}

// BoneMeal ...
func (c CocoaBean) BoneMeal(pos cube.Pos, w *world.World) bool {
	if c.Age == 2 {
		return false
	}
	c.Age++
	w.PlaceBlock(pos, c)
	return true
}

// HasLiquidDrops ...
func (c CocoaBean) HasLiquidDrops() bool {
	return true
}

// NeighbourUpdateTick ...
func (c CocoaBean) NeighbourUpdateTick(pos, _ cube.Pos, w *world.World) {
	if log, ok := w.Block(pos.Side(c.Facing.Face())).(Log); !ok || log.Wood != JungleWood() {
		w.BreakBlockWithoutParticles(pos)
	}
}

// UseOnBlock ...
func (c CocoaBean) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, w *world.World, user item.User, ctx *item.UseContext) bool {
	pos, _, used := firstReplaceable(w, pos, face, c)
	if !used {
		return false
	}

	if face == cube.FaceUp || face == cube.FaceDown {
		return false
	}
	if log, ok := w.Block(pos.Side(face.Opposite())).(Log); ok {
		if log.Wood == JungleWood() {
			c.Facing = face.Opposite().Direction()
			ctx.IgnoreAABB = true

			place(w, pos, c, user, ctx)
			return placed(ctx)
		}
	}

	return false
}

// RandomTick ...
func (c CocoaBean) RandomTick(pos cube.Pos, w *world.World, r *rand.Rand) {
	if c.Age < 2 && r.Intn(5) == 0 {
		c.Age++
		w.PlaceBlock(pos, c)
	}
}

// BreakInfo ...
func (c CocoaBean) BreakInfo() BreakInfo {
	return newBreakInfo(0.2, alwaysHarvestable, axeEffective, func(tool.Tool, []item.Enchantment) []item.Stack {
		if c.Age == 2 {
			return []item.Stack{item.NewStack(c, rand.Intn(2)+2)}
		}
		return []item.Stack{item.NewStack(c, 1)}
	})
}

// EncodeItem ...
func (c CocoaBean) EncodeItem() (name string, meta int16) {
	return "minecraft:cocoa_beans", 0
}

// EncodeBlock ...
func (c CocoaBean) EncodeBlock() (name string, properties map[string]interface{}) {
	direction := 2
	switch c.Facing {
	case cube.South:
		direction = 0
	case cube.West:
		direction = 1
	case cube.East:
		direction = 3
	}

	return "minecraft:cocoa", map[string]interface{}{"age": int32(c.Age), "direction": int32(direction)}
}

// Model ...
func (c CocoaBean) Model() world.BlockModel {
	return model.CocoaBean{Facing: c.Facing, Age: c.Age}
}

// allCocoaBeans ...
func allCocoaBeans() (cocoa []world.Block) {
	for i := cube.Direction(0); i <= 3; i++ {
		cocoa = append(cocoa, CocoaBean{Facing: i, Age: 0})
		cocoa = append(cocoa, CocoaBean{Facing: i, Age: 1})
		cocoa = append(cocoa, CocoaBean{Facing: i, Age: 2})
	}
	return
}
