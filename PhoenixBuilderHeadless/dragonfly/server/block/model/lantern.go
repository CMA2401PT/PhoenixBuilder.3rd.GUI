package model

import (
	"phoenixbuilder/dragonfly/server/block/cube"
	"phoenixbuilder/dragonfly/server/entity/physics"
	"phoenixbuilder/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
)

// Lantern is a model for the lantern block.
type Lantern struct {
	Hanging bool
}

// AABB ...
func (l Lantern) AABB(pos cube.Pos, w *world.World) []physics.AABB {
	if l.Hanging {
		return []physics.AABB{physics.NewAABB(mgl64.Vec3{0.3125, 0.125, 0.3125}, mgl64.Vec3{0.6875, 0.625, 0.6875})}
	}
	return []physics.AABB{physics.NewAABB(mgl64.Vec3{0.3125, 0, 0.3125}, mgl64.Vec3{0.6875, 0.5, 0.6875})}
}

// FaceSolid ...
func (l Lantern) FaceSolid(pos cube.Pos, face cube.Face, w *world.World) bool {
	return false
}
