package block

import (
	"phoenixbuilder_3rd_gui/fb/dragonfly/server/block/cube"
	"phoenixbuilder_3rd_gui/fb/dragonfly/server/block/instrument"
	"phoenixbuilder_3rd_gui/fb/dragonfly/server/item"
	"phoenixbuilder_3rd_gui/fb/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
)

// BoneBlock is a decorative block that can face different directions.
type BoneBlock struct {
	solid

	// Axis is the axis which the bone block faces.
	Axis cube.Axis
}

// Instrument ...
func (b BoneBlock) Instrument() instrument.Instrument {
	return instrument.Xylophone()
}

// UseOnBlock handles the rotational placing of bone blocks.
func (b BoneBlock) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, w *world.World, user item.User, ctx *item.UseContext) (used bool) {
	pos, face, used = firstReplaceable(w, pos, face, b)
	if !used {
		return
	}
	b.Axis = face.Axis()

	place(w, pos, b, user, ctx)
	return placed(ctx)
}

// BreakInfo ...
func (b BoneBlock) BreakInfo() BreakInfo {
	return newBreakInfo(2, pickaxeHarvestable, pickaxeEffective, oneOf(b))
}

// EncodeItem ...
func (b BoneBlock) EncodeItem() (name string, meta int16) {
	return "minecraft:bone_block", 0
}

// EncodeBlock ...
func (b BoneBlock) EncodeBlock() (name string, properties map[string]interface{}) {
	return "minecraft:bone_block", map[string]interface{}{"pillar_axis": b.Axis.String(), "deprecated": int32(0)}
}

// allBoneBlock ...
func allBoneBlock() (boneBlocks []world.Block) {
	for _, axis := range cube.Axes() {
		boneBlocks = append(boneBlocks, BoneBlock{Axis: axis})
	}
	return
}
