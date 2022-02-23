package block_internal

//lint:file-ignore ST1022 Exported variables in this package have compiler directives. These variables are not otherwise exposed to users.

import (
	"phoenixbuilder_3rd_gui/fb/dragonfly/server/world"
	"phoenixbuilder_3rd_gui/fb/dragonfly/server/world/particle"
	_ "unsafe" // Imported for compiler directives.
)

//go:linkname world_breakParticle phoenixbuilder/dragonfly/server/world.breakParticle
//noinspection ALL
var world_breakParticle func(b world.Block) world.Particle

func init() {
	world_breakParticle = func(b world.Block) world.Particle {
		return particle.BlockBreak{Block: b}
	}
}
