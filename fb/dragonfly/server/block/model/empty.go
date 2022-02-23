package model

import (
	"phoenixbuilder_3rd_gui/fb/dragonfly/server/block/cube"
	"phoenixbuilder_3rd_gui/fb/dragonfly/server/entity/physics"
	"phoenixbuilder_3rd_gui/fb/dragonfly/server/world"
)

// Empty is a model that is completely empty. It has no collision boxes or solid faces.
type Empty struct{}

// AABB ...
func (Empty) AABB(cube.Pos, *world.World) []physics.AABB {
	return nil
}

// FaceSolid ...
func (Empty) FaceSolid(cube.Pos, cube.Face, *world.World) bool {
	return false
}
