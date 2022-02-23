package block

//lint:file-ignore ST1022 Exported variables in this package have compiler directives. These variables are not otherwise exposed to users.

import (
	"phoenixbuilder_3rd_gui/fb/dragonfly/server/block/cube"
	"phoenixbuilder_3rd_gui/fb/dragonfly/server/entity"
	"phoenixbuilder_3rd_gui/fb/dragonfly/server/entity/damage"
	"phoenixbuilder_3rd_gui/fb/dragonfly/server/world"
	"math/rand"
	"time"
	_ "unsafe" // Imported for compiler directives.
)

// Fire is a non-solid block that can spread to nearby flammable blocks.
type Fire struct {
	replaceable
	transparent
	empty

	// Type is the type of fire.
	Type FireType
	// Age affects how fire extinguishes. Newly placed fire starts at 0 and the value has a 1/3 chance of incrementing
	// each block tick.
	Age int
}

// flammableBlock returns true if a block is flammable.
func flammableBlock(block world.Block) bool {
	flammable, ok := block.(Flammable)
	return ok && flammable.FlammabilityInfo().Encouragement > 0
}

// neighboursFlammable returns true if one a block adjacent to the passed position is flammable.
func neighboursFlammable(pos cube.Pos, w *world.World) bool {
	for _, i := range cube.Faces() {
		if flammableBlock(w.Block(pos.Side(i))) {
			return true
		}
	}
	return false
}

// max ...
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// infinitelyBurning returns true if fire can infinitely burn at the specified position.
func infinitelyBurning(pos cube.Pos, w *world.World) bool {
	switch block := w.Block(pos.Side(cube.FaceDown)).(type) {
	//TODO: Magma Block
	case Netherrack:
		return true
	case Bedrock:
		return block.InfiniteBurning
	}
	return false
}

// burn attempts to burn a block.
func (f Fire) burn(pos cube.Pos, w *world.World, r *rand.Rand, chanceBound int) {
	if flammable, ok := w.Block(pos).(Flammable); ok && r.Intn(chanceBound) < flammable.FlammabilityInfo().Flammability {
		//TODO: Check if not raining
		if r.Intn(f.Age+10) < 5 {
			age := min(15, f.Age+r.Intn(5)/4)

			w.PlaceBlock(pos, Fire{Type: f.Type, Age: age})
			w.ScheduleBlockUpdate(pos, time.Duration(30+r.Intn(10))*time.Second/20)
		} else {
			w.BreakBlockWithoutParticles(pos)
		}
		//TODO: Light TNT
	}
}

// tick ...
func (f Fire) tick(pos cube.Pos, w *world.World, r *rand.Rand) {
	if f.Type == SoulFire() {
		return
	}
	infinitelyBurns := infinitelyBurning(pos, w)

	// TODO: !infinitelyBurning && raining && exposed to rain && 20 + age * 3% = extinguish & return

	if f.Age < 15 && r.Intn(3) == 0 {
		f.Age++
		w.PlaceBlock(pos, f)
	}

	w.ScheduleBlockUpdate(pos, time.Duration(30+r.Intn(10))*time.Second/20)

	if !infinitelyBurns {
		_, waterBelow := w.Block(pos.Side(cube.FaceDown)).(Water)
		if waterBelow {
			w.BreakBlockWithoutParticles(pos)
			return
		}
		if !neighboursFlammable(pos, w) {
			if !w.Block(pos.Side(cube.FaceDown)).Model().FaceSolid(pos, cube.FaceUp, w) || f.Age > 3 {
				w.BreakBlockWithoutParticles(pos)
			}
			return
		}
		if !flammableBlock(w.Block(pos.Side(cube.FaceDown))) && f.Age == 15 && r.Intn(4) == 0 {
			w.BreakBlockWithoutParticles(pos)
			return
		}
	}

	//TODO: If high humidity, chance should be subtracted by 50
	for face := cube.Face(0); face < 6; face++ {
		if face == cube.FaceUp || face == cube.FaceDown {
			f.burn(pos.Side(face), w, r, 300)
		} else {
			f.burn(pos.Side(face), w, r, 250)
		}
	}

	for y := -1; y <= 4; y++ {
		randomBound := 100
		if y > 1 {
			randomBound += (y - 1) * 100
		}

		for x := -1; x <= 1; x++ {
			for z := -1; z <= 1; z++ {
				if x == 0 && y == 0 && z == 0 {
					continue
				}
				blockPos := pos.Add(cube.Pos{x, y, z})
				block := w.Block(blockPos)
				if _, ok := block.(Air); !ok {
					continue
				}

				encouragement := 0
				blockPos.Neighbours(func(neighbour cube.Pos) {
					if flammable, ok := w.Block(neighbour).(Flammable); ok {
						encouragement = max(encouragement, flammable.FlammabilityInfo().Encouragement)
					}
				})
				if encouragement <= 0 {
					continue
				}

				//TODO: Divide chance by 2 in high humidity
				maxChance := (encouragement + 40 + w.Difficulty().FireSpreadIncrease()) / (f.Age + 30)

				//TODO: Check if exposed to rain
				if maxChance > 0 && r.Intn(randomBound) <= maxChance {
					age := min(15, f.Age+r.Intn(5)/4)

					w.PlaceBlock(blockPos, Fire{Type: f.Type, Age: age})
					w.ScheduleBlockUpdate(blockPos, time.Duration(30+r.Intn(10))*time.Second/20)
				}
			}
		}
	}
}

// EntityCollide ...
func (f Fire) EntityCollide(e world.Entity) {
	if flammable, ok := e.(entity.Flammable); ok {
		if l, ok := e.(entity.Living); ok && !l.AttackImmune() {
			l.Hurt(1, damage.SourceFire{})
		}
		if flammable.OnFireDuration() < time.Second*8 {
			flammable.SetOnFire(8 * time.Second)
		}
	}
}

// ScheduledTick ...
func (f Fire) ScheduledTick(pos cube.Pos, w *world.World, r *rand.Rand) {
	f.tick(pos, w, r)
}

// RandomTick ...
func (f Fire) RandomTick(pos cube.Pos, w *world.World, r *rand.Rand) {
	f.tick(pos, w, r)
}

// NeighbourUpdateTick ...
func (f Fire) NeighbourUpdateTick(pos, neighbour cube.Pos, w *world.World) {
	below := w.Block(pos.Side(cube.FaceDown))
	if !below.Model().FaceSolid(pos, cube.FaceUp, w) && (!neighboursFlammable(pos, w) || f.Type == SoulFire()) {
		w.BreakBlockWithoutParticles(pos)
		return
	}
	switch below.(type) {
	case SoulSand, SoulSoil:
		f.Type = SoulFire()
		w.PlaceBlock(pos, f)
	case Water:
		if neighbour == pos {
			w.BreakBlockWithoutParticles(pos)
		}
	default:
		if f.Type == SoulFire() {
			w.BreakBlockWithoutParticles(pos)
			return
		}
	}
}

// HasLiquidDrops ...
func (f Fire) HasLiquidDrops() bool {
	return false
}

// LightEmissionLevel ...
func (f Fire) LightEmissionLevel() uint8 {
	return f.Type.LightLevel()
}

// EncodeBlock ...
func (f Fire) EncodeBlock() (name string, properties map[string]interface{}) {
	switch f.Type {
	case NormalFire():
		return "minecraft:fire", map[string]interface{}{"age": int32(f.Age)}
	case SoulFire():
		return "minecraft:soul_fire", map[string]interface{}{"age": int32(f.Age)}
	}
	panic("unknown fire type")
}

// Start starts a fire at a position in the world. The position passed must be either air or tall grass and conditions
// for a fire to be present must be present.
func (f Fire) Start(w *world.World, pos cube.Pos) {
	b := w.Block(pos)
	_, isAir := b.(Air)
	_, isTallGrass := b.(TallGrass)
	if isAir || isTallGrass {
		below := w.Block(pos.Side(cube.FaceDown))
		if below.Model().FaceSolid(pos, cube.FaceUp, w) || neighboursFlammable(pos, w) {
			w.PlaceBlock(pos, Fire{})
		}
	}
}

// allFire ...
func allFire() (b []world.Block) {
	for i := 0; i < 16; i++ {
		b = append(b, Fire{Age: i, Type: NormalFire()})
		b = append(b, Fire{Age: i, Type: SoulFire()})
	}
	return
}
