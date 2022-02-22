package model

import (
	"phoenixbuilder/dragonfly/server/block/cube"
	"phoenixbuilder/dragonfly/server/entity/physics"
	"phoenixbuilder/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
)

// CocoaBean is a model used by cocoa bean blocks.
type CocoaBean struct {
	Facing cube.Direction
	Age    int
}

// AABB ...
func (c CocoaBean) AABB(pos cube.Pos, w *world.World) []physics.AABB {
	return []physics.AABB{physics.NewAABB(mgl64.Vec3{}, mgl64.Vec3{1, 1, 1}).
		Stretch(c.Facing.RotateRight().Face().Axis(), -(6-float64(c.Age))/16).
		ExtendTowards(cube.FaceDown, -0.25).
		ExtendTowards(cube.FaceUp, -((7-float64(c.Age)*2)/16)).
		ExtendTowards(c.Facing.Face(), -0.0625).
		ExtendTowards(c.Facing.Opposite().Face(), -((11 - float64(c.Age)*2) / 16))}
}

// FaceSolid ...
func (c CocoaBean) FaceSolid(pos cube.Pos, face cube.Face, w *world.World) bool {
	return false
}
