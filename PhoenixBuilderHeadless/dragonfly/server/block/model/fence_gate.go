package model

import (
	"phoenixbuilder/dragonfly/server/block/cube"
	"phoenixbuilder/dragonfly/server/entity/physics"
	"phoenixbuilder/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
)

// FenceGate is a model used by fence gates.
type FenceGate struct {
	Facing cube.Direction
	Open   bool
}

// AABB ...
func (f FenceGate) AABB(pos cube.Pos, w *world.World) []physics.AABB {
	if f.Open {
		return nil
	}
	return []physics.AABB{physics.NewAABB(mgl64.Vec3{}, mgl64.Vec3{1, 1.5, 1}).Stretch(f.Facing.Face().Axis(), -0.375)}
}

// FaceSolid ...
func (f FenceGate) FaceSolid(pos cube.Pos, face cube.Face, w *world.World) bool {
	return false
}
