package world_provider

import (
	"phoenixbuilder_3rd_gui/fb/dragonfly/server/world"
	"phoenixbuilder_3rd_gui/fb/minecraft"
	"phoenixbuilder_3rd_gui/fb/minecraft/protocol/packet"
)

var CurrentWorld *world.World = nil

func Create(conn *minecraft.Conn) *world.World {
	intw := world.New(&StubLogger{}, 32)
	intw.Provider(NewOnlineWorldProvider(conn))
	return intw
}

func NewWorld(conn *minecraft.Conn) {
	ChunkCache = make(map[world.ChunkPos]*packet.LevelChunk)
	CurrentWorld = Create(conn)
	firstLoaded = false
}

func DestroyWorld() {
	firstLoaded = false
	CurrentWorld = nil
	ChunkCache = nil
}

func init() {
	InitRuntimeIdsWithoutMinecraftPrefix()
}
