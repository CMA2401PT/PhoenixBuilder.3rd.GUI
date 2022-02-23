package block

import (
	"phoenixbuilder_3rd_gui/fb/dragonfly/server/block/cube"
	"phoenixbuilder_3rd_gui/fb/dragonfly/server/item"
	"phoenixbuilder_3rd_gui/fb/dragonfly/server/item/tool"
	"phoenixbuilder_3rd_gui/fb/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"math/rand"
)

// BeetrootSeeds are a crop that can be harvested to craft soup or red dye.
type BeetrootSeeds struct {
	crop
}

// SameCrop ...
func (BeetrootSeeds) SameCrop(c Crop) bool {
	_, ok := c.(BeetrootSeeds)
	return ok
}

// BoneMeal ...
func (b BeetrootSeeds) BoneMeal(pos cube.Pos, w *world.World) bool {
	if b.Growth == 7 {
		return false
	}
	if rand.Float64() < 0.75 {
		b.Growth++
		w.PlaceBlock(pos, b)
		return true
	}
	return false
}

// UseOnBlock ...
func (b BeetrootSeeds) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, w *world.World, user item.User, ctx *item.UseContext) bool {
	pos, _, used := firstReplaceable(w, pos, face, b)
	if !used {
		return false
	}

	if _, ok := w.Block(pos.Side(cube.FaceDown)).(Farmland); !ok {
		return false
	}

	place(w, pos, b, user, ctx)
	return placed(ctx)
}

// BreakInfo ...
func (b BeetrootSeeds) BreakInfo() BreakInfo {
	return newBreakInfo(0, alwaysHarvestable, nothingEffective, func(tool.Tool, []item.Enchantment) []item.Stack {
		if b.Growth < 7 {
			return []item.Stack{item.NewStack(b, 1)}
		}
		return []item.Stack{item.NewStack(item.Beetroot{}, 1), item.NewStack(b, rand.Intn(4)+1)}
	})
}

// EncodeItem ...
func (b BeetrootSeeds) EncodeItem() (name string, meta int16) {
	return "minecraft:beetroot_seeds", 0
}

// RandomTick ...
func (b BeetrootSeeds) RandomTick(pos cube.Pos, w *world.World, r *rand.Rand) {
	if w.Light(pos) < 8 {
		w.BreakBlock(pos)
	} else if b.Growth < 7 && r.Intn(3) > 0 && r.Float64() <= b.CalculateGrowthChance(pos, w) {
		b.Growth++
		w.PlaceBlock(pos, b)
	}
}

// EncodeBlock ...
func (b BeetrootSeeds) EncodeBlock() (name string, properties map[string]interface{}) {
	return "minecraft:beetroot", map[string]interface{}{"growth": int32(b.Growth)}
}

// allBeetroot ...
func allBeetroot() (beetroot []world.Block) {
	for i := 0; i <= 7; i++ {
		beetroot = append(beetroot, BeetrootSeeds{crop: crop{Growth: i}})
	}
	return
}
