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

// Carrot is a crop that can be consumed raw.
type Carrot struct {
	crop
}

// SameCrop ...
func (Carrot) SameCrop(c Crop) bool {
	_, ok := c.(Carrot)
	return ok
}

// AlwaysConsumable ...
func (c Carrot) AlwaysConsumable() bool {
	return false
}

// ConsumeDuration ...
func (c Carrot) ConsumeDuration() time.Duration {
	return item.DefaultConsumeDuration
}

// Consume ...
func (c Carrot) Consume(_ *world.World, consumer item.Consumer) item.Stack {
	consumer.Saturate(3, 3.6)
	return item.Stack{}
}

// BoneMeal ...
func (c Carrot) BoneMeal(pos cube.Pos, w *world.World) bool {
	if c.Growth == 7 {
		return false
	}
	c.Growth = min(c.Growth+rand.Intn(4)+2, 7)
	w.PlaceBlock(pos, c)
	return true
}

// UseOnBlock ...
func (c Carrot) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, w *world.World, user item.User, ctx *item.UseContext) bool {
	pos, _, used := firstReplaceable(w, pos, face, c)
	if !used {
		return false
	}

	if _, ok := w.Block(pos.Side(cube.FaceDown)).(Farmland); !ok {
		return false
	}

	place(w, pos, c, user, ctx)
	return placed(ctx)
}

// BreakInfo ...
func (c Carrot) BreakInfo() BreakInfo {
	return newBreakInfo(0, alwaysHarvestable, nothingEffective, func(tool.Tool, []item.Enchantment) []item.Stack {
		if c.Growth < 7 {
			return []item.Stack{item.NewStack(c, 1)}
		}
		return []item.Stack{item.NewStack(c, rand.Intn(4)+2)}
	})
}

// EncodeItem ...
func (c Carrot) EncodeItem() (name string, meta int16) {
	return "minecraft:carrot", 0
}

// RandomTick ...
func (c Carrot) RandomTick(pos cube.Pos, w *world.World, r *rand.Rand) {
	if w.Light(pos) < 8 {
		w.BreakBlock(pos)
	} else if c.Growth < 7 && r.Float64() <= c.CalculateGrowthChance(pos, w) {
		c.Growth++
		w.PlaceBlock(pos, c)
	}
}

// EncodeBlock ...
func (c Carrot) EncodeBlock() (name string, properties map[string]interface{}) {
	return "minecraft:carrots", map[string]interface{}{"growth": int32(c.Growth)}
}

// allCarrots ...
func allCarrots() (carrots []world.Block) {
	for growth := 0; growth < 8; growth++ {
		carrots = append(carrots, Carrot{crop{Growth: growth}})
	}
	return
}
