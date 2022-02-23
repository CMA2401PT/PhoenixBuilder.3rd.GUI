package effect

import (
	"phoenixbuilder_3rd_gui/fb/dragonfly/server/entity/healing"
	"phoenixbuilder_3rd_gui/fb/dragonfly/server/world"
	"image/color"
	"time"
)

// Regeneration is an effect that causes the entity that it is added to to slowly regenerate health. The level
// of the effect influences the speed with which the entity regenerates.
type Regeneration struct {
	nopLasting
}

// Apply applies health to the world.Entity passed if the duration of the effect is at the right tick.
func (Regeneration) Apply(e world.Entity, lvl int, d time.Duration) {
	interval := 50 >> lvl
	if tickDuration(d)%interval == 0 {
		if l, ok := e.(living); ok {
			l.Heal(1, healing.SourceRegenerationEffect{})
		}
	}
}

// RGBA ...
func (Regeneration) RGBA() color.RGBA {
	return color.RGBA{R: 0xcd, G: 0x5c, B: 0xab, A: 0xff}
}
