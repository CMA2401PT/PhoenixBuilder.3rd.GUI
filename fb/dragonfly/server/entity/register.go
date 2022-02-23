package entity

import "phoenixbuilder_3rd_gui/fb/dragonfly/server/world"

// init registers all entities that can be saved in a world.World, so that they can be loaded when found in the world
// save.
func init() {
	world.RegisterEntity(&Text{})
	world.RegisterEntity(&FallingBlock{})
	world.RegisterEntity(&Item{})
}
