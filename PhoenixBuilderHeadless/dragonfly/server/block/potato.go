package block

import (
	"phoenixbuilder/dragonfly/server/block/cube"
	"phoenixbuilder/dragonfly/server/item"
	"phoenixbuilder/dragonfly/server/item/tool"
	"phoenixbuilder/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"math/rand"
	"time"
)

// Potato is a crop that can be consumed raw or cooked to make baked potatoes.
type Potato struct {
	crop
}

// SameCrop ...
func (Potato) SameCrop(c Crop) bool {
	_, ok := c.(Potato)
	return ok
}

// AlwaysConsumable ...
func (p Potato) AlwaysConsumable() bool {
	return false
}

// ConsumeDuration ...
func (p Potato) ConsumeDuration() time.Duration {
	return item.DefaultConsumeDuration
}

// Consume ...
func (p Potato) Consume(_ *world.World, c item.Consumer) item.Stack {
	c.Saturate(1, 0.6)
	return item.Stack{}
}

// BoneMeal ...
func (p Potato) BoneMeal(pos cube.Pos, w *world.World) bool {
	if p.Growth == 7 {
		return false
	}
	p.Growth = min(p.Growth+rand.Intn(4)+2, 7)
	w.PlaceBlock(pos, p)
	return true
}

// UseOnBlock ...
func (p Potato) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, w *world.World, user item.User, ctx *item.UseContext) bool {
	pos, _, used := firstReplaceable(w, pos, face, p)
	if !used {
		return false
	}

	if _, ok := w.Block(pos.Side(cube.FaceDown)).(Farmland); !ok {
		return false
	}

	place(w, pos, p, user, ctx)
	return placed(ctx)
}

// BreakInfo ...
func (p Potato) BreakInfo() BreakInfo {
	return newBreakInfo(0, alwaysHarvestable, nothingEffective, func(tool.Tool, []item.Enchantment) []item.Stack {
		if rand.Float64() < 0.02 {
			return []item.Stack{item.NewStack(p, rand.Intn(5)+1), item.NewStack(item.PoisonousPotato{}, 1)}
		}
		return []item.Stack{item.NewStack(p, rand.Intn(5)+1)}
	})
}

// EncodeItem ...
func (p Potato) EncodeItem() (name string, meta int16) {
	return "minecraft:potato", 0
}

// RandomTick ...
func (p Potato) RandomTick(pos cube.Pos, w *world.World, r *rand.Rand) {
	if w.Light(pos) < 8 {
		w.BreakBlock(pos)
	} else if p.Growth < 7 && r.Float64() <= p.CalculateGrowthChance(pos, w) {
		p.Growth++
		w.PlaceBlock(pos, p)
	}
}

// EncodeBlock ...
func (p Potato) EncodeBlock() (name string, properties map[string]interface{}) {
	return "minecraft:potatoes", map[string]interface{}{"growth": int32(p.Growth)}
}

// allPotato ...
func allPotato() (potato []world.Block) {
	for i := 0; i <= 7; i++ {
		potato = append(potato, Potato{crop{Growth: i}})
	}
	return
}
