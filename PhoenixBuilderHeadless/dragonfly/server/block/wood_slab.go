package block

import (
	"phoenixbuilder/dragonfly/server/block/cube"
	"phoenixbuilder/dragonfly/server/block/model"
	"phoenixbuilder/dragonfly/server/entity/physics"
	"phoenixbuilder/dragonfly/server/item"
	"phoenixbuilder/dragonfly/server/item/tool"
	"phoenixbuilder/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
)

// WoodSlab is a half block that allows entities to walk up blocks without jumping.
type WoodSlab struct {
	bass

	// Wood is the type of wood of the slabs. This field must have one of the values found in the material
	// package.
	Wood WoodType
	// Top specifies if the slab is in the top part of the block.
	Top bool
	// Double specifies if the slab is a double slab. These double slabs can be made by placing another slab
	// on an existing slab.
	Double bool
}

// FlammabilityInfo ...
func (s WoodSlab) FlammabilityInfo() FlammabilityInfo {
	if !s.Wood.Flammable() {
		return newFlammabilityInfo(0, 0, false)
	}
	return newFlammabilityInfo(5, 20, true)
}

// Model ...
func (s WoodSlab) Model() world.BlockModel {
	return model.Slab{Double: s.Double, Top: s.Top}
}

// UseOnBlock handles the placement of slabs with relation to them being upside down or not and handles slabs
// being turned into double slabs.
func (s WoodSlab) UseOnBlock(pos cube.Pos, face cube.Face, clickPos mgl64.Vec3, w *world.World, user item.User, ctx *item.UseContext) (used bool) {
	clickedBlock := w.Block(pos)
	if clickedSlab, ok := clickedBlock.(WoodSlab); ok && !s.Double {
		if (face == cube.FaceUp && !clickedSlab.Double && clickedSlab.Wood == s.Wood && !clickedSlab.Top) ||
			(face == cube.FaceDown && !clickedSlab.Double && clickedSlab.Wood == s.Wood && clickedSlab.Top) {
			// A half slab of the same type was clicked at the top, so we can make it full.
			clickedSlab.Double = true

			place(w, pos, clickedSlab, user, ctx)
			return placed(ctx)
		}
	}
	if sideSlab, ok := w.Block(pos.Side(face)).(WoodSlab); ok && !replaceableWith(w, pos, s) && !s.Double {
		// The block on the side of the one clicked was a slab and the block clicked was not replaceableWith, so
		// the slab on the side must've been half and may now be filled if the wood types are the same.
		if !sideSlab.Double && sideSlab.Wood == s.Wood {
			sideSlab.Double = true

			place(w, pos.Side(face), sideSlab, user, ctx)
			return placed(ctx)
		}
	}
	pos, face, used = firstReplaceable(w, pos, face, s)
	if !used {
		return
	}
	if face == cube.FaceDown || (clickPos[1] > 0.5 && face != cube.FaceUp) {
		s.Top = true
	}

	place(w, pos, s, user, ctx)
	return placed(ctx)
}

// BreakInfo ...
func (s WoodSlab) BreakInfo() BreakInfo {
	return newBreakInfo(2, alwaysHarvestable, axeEffective, func(tool.Tool, []item.Enchantment) []item.Stack {
		if s.Double {
			s.Double = false
			// If the slab is double, it should drop two single slabs.
			return []item.Stack{item.NewStack(s, 2)}
		}
		return []item.Stack{item.NewStack(s, 1)}
	})
}

// LightDiffusionLevel returns 0 if the slab is a half slab, or 15 if it is double.
func (s WoodSlab) LightDiffusionLevel() uint8 {
	if s.Double {
		return 15
	}
	return 0
}

// AABB ...
func (s WoodSlab) AABB(cube.Pos, *world.World) []physics.AABB {
	if s.Double {
		return []physics.AABB{physics.NewAABB(mgl64.Vec3{}, mgl64.Vec3{1, 1, 1})}
	}
	if s.Top {
		return []physics.AABB{physics.NewAABB(mgl64.Vec3{0, 0.5, 0}, mgl64.Vec3{1, 1, 1})}
	}
	return []physics.AABB{physics.NewAABB(mgl64.Vec3{}, mgl64.Vec3{1, 0.5, 1})}
}

// EncodeItem ...
func (s WoodSlab) EncodeItem() (name string, meta int16) {
	switch s.Wood {
	case OakWood(), SpruceWood(), BirchWood(), JungleWood(), AcaciaWood(), DarkOakWood():
		if s.Double {
			return "minecraft:double_wooden_slab", int16(s.Wood.Uint8())
		}
		return "minecraft:wooden_slab", int16(s.Wood.Uint8())
	case CrimsonWood():
		if s.Double {
			return "minecraft:crimson_double_slab", 0
		}
		return "minecraft:crimson_slab", 0
	case WarpedWood():
		if s.Double {
			return "minecraft:warped_double_slab", 0
		}
		return "minecraft:warped_slab", 0
	}
	panic("invalid wood type")
}

// EncodeBlock ...
func (s WoodSlab) EncodeBlock() (name string, properties map[string]interface{}) {
	if s.Double {
		if s.Wood == CrimsonWood() || s.Wood == WarpedWood() {
			return "minecraft:" + s.Wood.String() + "_double_slab", map[string]interface{}{"top_slot_bit": s.Top}
		}
		return "minecraft:double_wooden_slab", map[string]interface{}{"top_slot_bit": s.Top, "wood_type": s.Wood.String()}
	}
	if s.Wood == CrimsonWood() || s.Wood == WarpedWood() {
		return "minecraft:" + s.Wood.String() + "_slab", map[string]interface{}{"top_slot_bit": s.Top}
	}
	return "minecraft:wooden_slab", map[string]interface{}{"top_slot_bit": s.Top, "wood_type": s.Wood.String()}
}

// CanDisplace ...
func (s WoodSlab) CanDisplace(b world.Liquid) bool {
	_, ok := b.(Water)
	return !s.Double && ok
}

// SideClosed ...
func (s WoodSlab) SideClosed(pos, side cube.Pos, _ *world.World) bool {
	// Only returns true if the side is below the slab and if the slab is not upside down.
	return !s.Top && side[1] == pos[1]-1
}

// allWoodSlabs returns all states of wood slabs.
func allWoodSlabs() (slabs []world.Block) {
	f := func(double bool, upsideDown bool) {
		for _, w := range WoodTypes() {
			slabs = append(slabs, WoodSlab{Double: double, Top: upsideDown, Wood: w})
		}
	}
	f(false, false)
	f(false, true)
	f(true, false)
	f(true, true)
	return
}
