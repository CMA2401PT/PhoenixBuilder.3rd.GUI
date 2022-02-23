package model

import (
	"phoenixbuilder_3rd_gui/fb/dragonfly/server/block/cube"
	"phoenixbuilder_3rd_gui/fb/dragonfly/server/entity/physics"
	"phoenixbuilder_3rd_gui/fb/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
)

// Carpet is a model for carpet-like extremely thin blocks.
type Carpet struct{}

// AABB ...
func (Carpet) AABB(cube.Pos, *world.World) []physics.AABB {
	return []physics.AABB{physics.NewAABB(mgl64.Vec3{}, mgl64.Vec3{1, 0.0625, 1})}
}

// FaceSolid ...
func (Carpet) FaceSolid(cube.Pos, cube.Face, *world.World) bool {
	return false
}
