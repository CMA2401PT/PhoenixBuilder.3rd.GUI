package item

import (
	"phoenixbuilder/dragonfly/server/entity/effect"
	"phoenixbuilder/dragonfly/server/world"
	"time"
)

// Pufferfish is a poisonous type of fish that is used to brew water breathing potions.
type Pufferfish struct {
	defaultFood
}

// Consume ...
func (p Pufferfish) Consume(_ *world.World, c Consumer) Stack {
	c.Saturate(1, 0.2)
	c.AddEffect(effect.New(effect.Hunger{}, 3, 15*time.Second))
	c.AddEffect(effect.New(effect.Poison{}, 4, time.Minute))
	c.AddEffect(effect.New(effect.Nausea{}, 1, 15*time.Second))
	return Stack{}
}

// EncodeItem ...
func (p Pufferfish) EncodeItem() (name string, meta int16) {
	return "minecraft:pufferfish", 0
}
