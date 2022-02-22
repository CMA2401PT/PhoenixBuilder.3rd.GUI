package world_provider

import (
	"fmt"
	"phoenixbuilder/dragonfly/server/world"
	"phoenixbuilder/dragonfly/server/world/chunk"
	"phoenixbuilder/fastbuilder/command"
	"phoenixbuilder/minecraft"
	"phoenixbuilder/minecraft/protocol/packet"
	bridge_fmt "phoenixbuilder/session/bridge/fmt"
	"runtime"
	"time"

	"github.com/google/uuid"
)

var ChunkInput chan *packet.LevelChunk = nil
var ChunkCache map[world.ChunkPos]*packet.LevelChunk = nil
var firstLoaded bool = false

type OnlineWorldProvider struct {
	connection *minecraft.Conn
	//nbtmap map[world.ChunkPos][]map[string]interface{}
}

func NewOnlineWorldProvider(conn *minecraft.Conn) *OnlineWorldProvider {
	return &OnlineWorldProvider{
		connection: conn,
		//nbtmap: make(map[world.ChunkPos][]map[string]interface{}),
	}
}

func (p *OnlineWorldProvider) Settings() world.Settings {
	return world.Settings{
		Name: "World",
	}
}

func (p *OnlineWorldProvider) SaveSettings(_ world.Settings) {

}

func DoCache(pkt *packet.LevelChunk) {
	if ChunkCache != nil {
		quickCache(pkt)
	}
}

func quickCache(pkt *packet.LevelChunk) {
	ChunkCache[world.ChunkPos{pkt.ChunkX, pkt.ChunkZ}] = pkt
}

func wander(conn *minecraft.Conn, position world.ChunkPos) {
	u_d, _ := uuid.NewUUID()
	err := command.SendWSCommand(fmt.Sprintf("tp %d 127 %d", position[0]*16+100000, 1000000-position[1]*16+100000), u_d, conn)
	if err != nil {
		panic(fmt.Errorf("Connection closed: %+v", err))
	}
	select {
	case <-ChunkInput:
		//quickCache(inp)
	case <-time.After(2 * time.Second):

	}
	u_d, _ = uuid.NewUUID()
	err = command.SendWSCommand(fmt.Sprintf("tp %d 127 %d", position[0]*16, position[1]*16), u_d, conn)
	if err != nil {
		panic(fmt.Errorf("[2]Connection closed: %+v", err))
	}
}

func (p *OnlineWorldProvider) LoadChunk(position world.ChunkPos) (c *chunk.Chunk, exists bool, err error) {
	if ChunkCache == nil {
		panic("LoadChunk() before creating a world")
	}
	cacheitem, hascacheitem := ChunkCache[position]
	if hascacheitem {
		delete(ChunkCache, position)
		chunk, err := chunk.NetworkDecode(AirRuntimeId, cacheitem.RawPayload, int(cacheitem.SubChunkCount))
		if err != nil {
			bridge_fmt.Printf("Failed to decode chunk: %v\n", err)
			return nil, true, err
		}
		return chunk, true, nil
	}
	if ChunkInput != nil {
		panic("Multithreading on OnlineWorldProvider's LoadChunk function isn't allowed")
	}
	u_d, _ := uuid.NewUUID()
	ChunkInput = make(chan *packet.LevelChunk, 32)
	err = command.SendWSCommand(fmt.Sprintf("tp %d 127 %d", position[0]*16, position[1]*16), u_d, p.connection)
	if err != nil {
		panic(fmt.Errorf("[2]Connection closed: %+v", err))
	}
	for {
		inp, hasqi := ChunkCache[position]
		if !hasqi {
			select {
			case inp = <-ChunkInput:
				quickCache(inp)
				bridge_fmt.Printf("Waiting for chunk: current: %d, %d | expected: %v\n", inp.ChunkX, inp.ChunkZ, position)
				if inp.ChunkX != position[0] || inp.ChunkZ != position[1] {
					continue
				}
			case <-time.After(2 * time.Second):
				runtime.GC()
				bridge_fmt.Printf("Expected chunk %v didn't arrive, wandering around\n", position)
				wander(p.connection, position)
				continue
			}
		} else {
			delete(ChunkCache, position)
		}
		// Hit
		close(ChunkInput)
		ChunkInput = nil
		chunk, err := chunk.NetworkDecode(AirRuntimeId, inp.RawPayload, int(inp.SubChunkCount))
		if err != nil {
			bridge_fmt.Printf("Failed to decode chunk: %v\n", err)
			return nil, true, err
		}
		/*blockentities:=bytes.NewReader(inp.RawPayload[len(inp.RawPayload)-ato:])
		blockentities.ReadByte()
		dec:=nbt.NewDecoderWithEncoding(blockentities, nbt.NetworkLittleEndian)
		nbtout:=make([]map[string]interface{},0)
		for {
			out:=make(map[string]interface{})
			err=dec.Decode(&out)
			if(err!=nil) {
				break
			}
			nbtout=append(nbtout,out)
		}
		p.nbtmap[position]=nbtout*/
		return chunk, true, nil
	}
}

func (p *OnlineWorldProvider) SaveChunk(position world.ChunkPos, c *chunk.Chunk) error {
	return nil
}

func (p *OnlineWorldProvider) LoadEntities(position world.ChunkPos) ([]world.SaveableEntity, error) {
	// Not implemented
	return []world.SaveableEntity{}, nil
}

func (p *OnlineWorldProvider) SaveEntities(position world.ChunkPos, entities []world.SaveableEntity) error {
	return nil
}

func (p *OnlineWorldProvider) LoadBlockNBT(position world.ChunkPos) ([]map[string]interface{}, error) {
	return nil, nil
	/*r, h:=p.nbtmap[position]
	if(!h) {
		fmt.Printf("No NBT for position %v.\n",position)
		return nil, fmt.Errorf("NO NBT")
	}
	return r, nil*/
}

func (p *OnlineWorldProvider) SaveBlockNBT(position world.ChunkPos, data []map[string]interface{}) error {
	return nil
}

func (p *OnlineWorldProvider) Close() error {
	return nil
}
