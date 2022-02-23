package block

import (
	"phoenixbuilder_3rd_gui/fb/dragonfly/server/block/cube"
	"phoenixbuilder_3rd_gui/fb/dragonfly/server/block/instrument"
	"phoenixbuilder_3rd_gui/fb/dragonfly/server/item"
	"phoenixbuilder_3rd_gui/fb/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
)

// Pumpkin is a crop block. Interacting with shears results in the carved variant.
type Pumpkin struct {
	solid

	// Carved is whether the pumpkin is carved.
	Carved bool
	// Facing is the direction the pumpkin is facing.
	Facing cube.Direction
}

// Instrument ...
func (p Pumpkin) Instrument() instrument.Instrument {
	if !p.Carved {
		return instrument.Didgeridoo()
	}
	return instrument.Piano()
}

// UseOnBlock ...
func (p Pumpkin) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, w *world.World, user item.User, ctx *item.UseContext) (used bool) {
	pos, _, used = firstReplaceable(w, pos, face, p)
	if !used {
		return
	}
	p.Facing = user.Facing().Opposite()

	place(w, pos, p, user, ctx)
	return placed(ctx)
}

// BreakInfo ...
func (p Pumpkin) BreakInfo() BreakInfo {
	return newBreakInfo(1, alwaysHarvestable, axeEffective, oneOf(p))
}

// Carve ...
func (p Pumpkin) Carve(f cube.Face) (world.Block, bool) {
	return Pumpkin{Facing: f.Direction(), Carved: true}, !p.Carved
}

// Helmet ...
func (p Pumpkin) Helmet() bool {
	return p.Carved
}

// DefencePoints ...
func (p Pumpkin) DefencePoints() float64 {
	return 0
}

// KnockBackResistance ...
func (p Pumpkin) KnockBackResistance() float64 {
	return 0
}

// EncodeItem ...
func (p Pumpkin) EncodeItem() (name string, meta int16) {
	if p.Carved {
		return "minecraft:carved_pumpkin", 0
	}
	return "minecraft:pumpkin", 0
}

// EncodeBlock ...
func (p Pumpkin) EncodeBlock() (name string, properties map[string]interface{}) {
	direction := 2
	switch p.Facing {
	case cube.South:
		direction = 0
	case cube.West:
		direction = 1
	case cube.East:
		direction = 3
	}

	if p.Carved {
		return "minecraft:carved_pumpkin", map[string]interface{}{"direction": int32(direction)}
	}
	return "minecraft:pumpkin", map[string]interface{}{"direction": int32(direction)}
}

func allPumpkins() (pumpkins []world.Block) {
	for i := cube.Direction(0); i <= 3; i++ {
		pumpkins = append(pumpkins, Pumpkin{Facing: i})
		pumpkins = append(pumpkins, Pumpkin{Facing: i, Carved: true})
	}
	return
}
