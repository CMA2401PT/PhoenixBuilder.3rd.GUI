package model

import (
	"phoenixbuilder/dragonfly/server/block/cube"
	"phoenixbuilder/dragonfly/server/entity/physics"
	"phoenixbuilder/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
)

// Chest is the model of a chest. It is just barely not a full block, having a slightly reduced with on all
// axes.
type Chest struct{}

// AABB ...
func (Chest) AABB(pos cube.Pos, w *world.World) []physics.AABB {
	return []physics.AABB{physics.NewAABB(mgl64.Vec3{0.025, 0, 0.025}, mgl64.Vec3{0.975, 0.95, 0.975})}
}

// FaceSolid ...
func (Chest) FaceSolid(pos cube.Pos, face cube.Face, w *world.World) bool {
	return false
}
