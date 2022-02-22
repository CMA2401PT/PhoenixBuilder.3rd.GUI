package world

import (
	"phoenixbuilder/dragonfly/server/block/cube"
)

// Settings holds the settings of a World. These are typically saved to a level.dat file.
type Settings struct {
	// Name is the display name of the World.
	Name string
	// Spawn is the spawn position of the World. New players that join the world will be spawned here.
	Spawn cube.Pos
	// Time is the current time of the World. It advances every tick if TimeCycle is set to true.
	Time int64
	// TimeCycle specifies if the time should advance every tick. If set to false, time won't change.
	TimeCycle bool
	// CurrentTick is the current tick of the world. This is similar to the Time, except that it has no visible effect
	// to the client. It can also not be changed through commands and will only ever go up.
	CurrentTick int64
	// DefaultGameMode is the GameMode assigned to players that join the World for the first time.
	DefaultGameMode GameMode
	// Difficulty is the difficulty of the World. Behaviour of hunger, regeneration and monsters differs based on the
	// difficulty of the world.
	Difficulty Difficulty
}

// defaultSettings returns the default Settings for a new World.
func defaultSettings() Settings {
	return Settings{Name: "World", DefaultGameMode: GameModeSurvival{}, Difficulty: DifficultyNormal{}, TimeCycle: true}
}
