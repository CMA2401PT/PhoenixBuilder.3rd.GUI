package model

import (
	"phoenixbuilder_3rd_gui/fb/dragonfly/server/block/cube"
	"phoenixbuilder_3rd_gui/fb/dragonfly/server/entity/physics"
	"phoenixbuilder_3rd_gui/fb/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
)

// TilledGrass is a model used for grass that has been tilled in some way, such as dirt paths and farmland.
type TilledGrass struct {
}

// AABB ...
func (TilledGrass) AABB(pos cube.Pos, w *world.World) []physics.AABB {
	// TODO: Make the max Y value 0.9375 once https://bugs.mojang.com/browse/MCPE-12109 gets fixed.
	return []physics.AABB{physics.NewAABB(mgl64.Vec3{}, mgl64.Vec3{1, 1, 1})}
}

// FaceSolid ...
func (TilledGrass) FaceSolid(pos cube.Pos, face cube.Face, w *world.World) bool {
	return true
}
