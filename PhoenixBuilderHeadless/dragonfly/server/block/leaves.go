package block

import (
	"phoenixbuilder/dragonfly/server/block/cube"
	"phoenixbuilder/dragonfly/server/item"
	"phoenixbuilder/dragonfly/server/item/tool"
	"phoenixbuilder/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"math/rand"
)

// Leaves are blocks that grow as part of trees which mainly drop saplings and sticks.
type Leaves struct {
	leaves

	// Wood is the type of wood of the leaves. This field must have one of the values found in the material
	// package.
	Wood WoodType
	// Persistent specifies if the leaves are persistent, meaning they will not decay as a result of no wood
	// being nearby.
	Persistent bool

	ShouldUpdate bool
}

// UseOnBlock makes leaves persistent when they are placed so that they don't decay.
func (l Leaves) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, w *world.World, user item.User, ctx *item.UseContext) (used bool) {
	pos, _, used = firstReplaceable(w, pos, face, l)
	if !used {
		return
	}
	l.Persistent = true

	place(w, pos, l, user, ctx)
	return placed(ctx)
}

// findLog ...
func findLog(pos cube.Pos, w *world.World, visited *[]cube.Pos, distance int) bool {
	for _, v := range *visited {
		if v == pos {
			return false
		}
	}
	*visited = append(*visited, pos)

	if log, ok := w.Block(pos).(Log); ok && !log.Stripped {
		return true
	}
	if _, ok := w.Block(pos).(Leaves); !ok || distance > 6 {
		return false
	}
	logFound := false
	pos.Neighbours(func(neighbour cube.Pos) {
		if !logFound && findLog(neighbour, w, visited, distance+1) {
			logFound = true
		}
	})
	return logFound
}

// RandomTick ...
func (l Leaves) RandomTick(pos cube.Pos, w *world.World, _ *rand.Rand) {
	if !l.Persistent && l.ShouldUpdate {
		if findLog(pos, w, &[]cube.Pos{}, 0) {
			l.ShouldUpdate = false
			w.PlaceBlock(pos, l)
		} else {
			w.BreakBlockWithoutParticles(pos)
		}
	}
}

// NeighbourUpdateTick ...
func (l Leaves) NeighbourUpdateTick(pos, _ cube.Pos, w *world.World) {
	if !l.Persistent && !l.ShouldUpdate {
		l.ShouldUpdate = true
		w.PlaceBlock(pos, l)
	}
}

// FlammabilityInfo ...
func (l Leaves) FlammabilityInfo() FlammabilityInfo {
	return newFlammabilityInfo(30, 60, true)
}

// BreakInfo ...
func (l Leaves) BreakInfo() BreakInfo {
	return newBreakInfo(0.2, alwaysHarvestable, func(t tool.Tool) bool {
		return t.ToolType() == tool.TypeShears || t.ToolType() == tool.TypeHoe
	}, func(t tool.Tool, enchantments []item.Enchantment) []item.Stack {
		if t.ToolType() == tool.TypeShears || hasSilkTouch(enchantments) {
			return []item.Stack{item.NewStack(l, 1)}
		}
		var drops []item.Stack
		if (l.Wood == OakWood() || l.Wood == DarkOakWood()) && rand.Float64() < 0.005 {
			drops = append(drops, item.NewStack(item.Apple{}, 1))
		}
		// TODO: Saplings and sticks can drop along with apples
		return drops
	})
}

// EncodeItem ...
func (l Leaves) EncodeItem() (name string, meta int16) {
	switch l.Wood {
	case OakWood(), SpruceWood(), BirchWood(), JungleWood():
		return "minecraft:leaves", int16(l.Wood.Uint8())
	case AcaciaWood(), DarkOakWood():
		return "minecraft:leaves2", int16(l.Wood.Uint8() - 4)
	}
	panic("invalid wood type")
}

// LightDiffusionLevel ...
func (Leaves) LightDiffusionLevel() uint8 {
	return 1
}

// CanDisplace ...
func (Leaves) CanDisplace(b world.Liquid) bool {
	_, ok := b.(Water)
	return ok
}

// SideClosed ...
func (Leaves) SideClosed(cube.Pos, cube.Pos, *world.World) bool {
	return false
}

// EncodeBlock ...
func (l Leaves) EncodeBlock() (name string, properties map[string]interface{}) {
	switch l.Wood {
	case OakWood(), SpruceWood(), BirchWood(), JungleWood():
		return "minecraft:leaves", map[string]interface{}{"old_leaf_type": l.Wood.String(), "persistent_bit": l.Persistent, "update_bit": l.ShouldUpdate}
	case AcaciaWood(), DarkOakWood():
		return "minecraft:leaves2", map[string]interface{}{"new_leaf_type": l.Wood.String(), "persistent_bit": l.Persistent, "update_bit": l.ShouldUpdate}
	}
	panic("invalid wood type")
}

// allLogs returns a list of all possible leaves states.
func allLeaves() (leaves []world.Block) {
	f := func(persistent, update bool) {
		leaves = append(leaves, Leaves{Wood: OakWood(), Persistent: persistent, ShouldUpdate: update})
		leaves = append(leaves, Leaves{Wood: SpruceWood(), Persistent: persistent, ShouldUpdate: update})
		leaves = append(leaves, Leaves{Wood: BirchWood(), Persistent: persistent, ShouldUpdate: update})
		leaves = append(leaves, Leaves{Wood: JungleWood(), Persistent: persistent, ShouldUpdate: update})
		leaves = append(leaves, Leaves{Wood: AcaciaWood(), Persistent: persistent, ShouldUpdate: update})
		leaves = append(leaves, Leaves{Wood: DarkOakWood(), Persistent: persistent, ShouldUpdate: update})
	}
	f(true, true)
	f(true, false)
	f(false, true)
	f(false, false)
	return
}
