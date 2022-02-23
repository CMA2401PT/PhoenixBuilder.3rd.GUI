package model

import (
	"phoenixbuilder_3rd_gui/fb/dragonfly/server/block/cube"
	"phoenixbuilder_3rd_gui/fb/dragonfly/server/entity/physics"
	"phoenixbuilder_3rd_gui/fb/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
)

// Leaves is a model for leaves-like blocks. These blocks have a full collision box, but none of their faces
// are solid.
type Leaves struct{}

// AABB ...
func (Leaves) AABB(cube.Pos, *world.World) []physics.AABB {
	return []physics.AABB{physics.NewAABB(mgl64.Vec3{}, mgl64.Vec3{1, 1, 1})}
}

// FaceSolid ...
func (Leaves) FaceSolid(cube.Pos, cube.Face, *world.World) bool {
	return false
}
