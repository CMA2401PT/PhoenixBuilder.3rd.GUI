package effect

import (
	"phoenixbuilder_3rd_gui/fb/dragonfly/server/entity/damage"
	"phoenixbuilder_3rd_gui/fb/dragonfly/server/world"
	"image/color"
	"time"
)

// Wither is a lasting effect that causes an entity to take continuous damage that is capable of killing an
// entity.
type Wither struct {
	nopLasting
}

// Apply ...
func (Wither) Apply(e world.Entity, lvl int, d time.Duration) {
	interval := 80 >> lvl
	if tickDuration(d)%interval == 0 {
		if l, ok := e.(living); ok {
			l.Hurt(1, damage.SourceWitherEffect{})
		}
	}
}

// RGBA ...
func (Wither) RGBA() color.RGBA {
	return color.RGBA{R: 0x35, G: 0x2a, B: 0x27, A: 0xff}
}
