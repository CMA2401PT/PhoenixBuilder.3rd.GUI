package task

import (
	"fmt"
	"phoenixbuilder_3rd_gui/fb/dragonfly/server/block/cube"
	"phoenixbuilder_3rd_gui/fb/dragonfly/server/world"
	"phoenixbuilder_3rd_gui/fb/fastbuilder/bdump"
	"phoenixbuilder_3rd_gui/fb/fastbuilder/command"
	"phoenixbuilder_3rd_gui/fb/fastbuilder/configuration"
	"phoenixbuilder_3rd_gui/fb/fastbuilder/parsing"
	"phoenixbuilder_3rd_gui/fb/fastbuilder/types"
	"phoenixbuilder_3rd_gui/fb/fastbuilder/world_provider"
	"phoenixbuilder_3rd_gui/fb/minecraft"
	"phoenixbuilder_3rd_gui/fb/minecraft/protocol/packet"
	"runtime"
	"strings"
)

type SolidSimplePos struct {
	X int64 `json:"x"`
	Y int64 `json:"y"`
	Z int64 `json:"z"`
}

type SolidRet struct {
	BlockName  string         `json:"blockName"`
	Position   SolidSimplePos `json:"position"`
	StatusCode int64          `json:"statusCode"`
}

var ExportWaiter chan map[string]interface{}

func CreateExportTask(commandLine string, conn *minecraft.Conn) *Task {
	cfg, err := parsing.Parse(commandLine, configuration.GlobalFullConfig().Main())
	if err != nil {
		command.Tellraw(conn, fmt.Sprintf("Failed to parse command: %v", err))
		return nil
	}
	beginPos := cfg.Position
	endPos := cfg.End
	if endPos.X-beginPos.X < 0 {
		temp := endPos.X
		endPos.X = beginPos.X
		beginPos.X = temp
	}
	if endPos.Y-beginPos.Y < 0 {
		temp := endPos.Y
		endPos.Y = beginPos.Y
		beginPos.Y = temp
	}
	if endPos.Z-beginPos.Z < 0 {
		temp := endPos.Z
		endPos.Z = beginPos.Z
		beginPos.Z = temp
	}
	if world_provider.CurrentWorld != nil {
		command.Tellraw(conn, "EXPORT >> World interaction interface is occupied, failing")
		return nil
	}
	world_provider.NewWorld(conn)
	go func() {
		command.Tellraw(conn, "EXPORT >> Exporting...")
		V := (endPos.X - beginPos.X + 1) * (endPos.Y - beginPos.Y + 1) * (endPos.Z - beginPos.Z + 1)
		blocks := make([]*types.RuntimeModule, V)
		counter := 0
		for x := beginPos.X; x <= endPos.X; x++ {
			for z := beginPos.Z; z <= endPos.Z; z++ {
				for y := beginPos.Y; y <= endPos.Y; y++ {
					blk := world_provider.CurrentWorld.Block(cube.Pos{x, y, z})
					runtimeId := world.LoadRuntimeID(blk)
					if runtimeId == world_provider.AirRuntimeId {
						continue
					}
					block, item := blk.EncodeBlock()
					var cbdata *types.CommandBlockData = nil
					var chestData *types.ChestData = nil
					if block == "chest" || strings.Contains(block, "shulker_box") {
						content := item["Items"].([]interface{})
						chest := make(types.ChestData, len(content))
						for index, iface := range content {
							i := iface.(map[string]interface{})
							name := i["Name"].(string)
							count := i["Count"].(uint8)
							damage := i["Damage"].(int16)
							slot := i["Slot"].(uint8)
							name_mcnk := name[10:]
							chest[index] = types.ChestSlot{
								Name:   name_mcnk,
								Count:  count,
								Damage: uint16(int(damage)),
								Slot:   slot,
							}
						}
						chestData = &chest
					}
					if strings.Contains(block, "command_block") {
						var mode uint32
						if block == "command_block" {
							mode = packet.CommandBlockImpulse
						} else if block == "repeating_command_block" {
							mode = packet.CommandBlockRepeat
						} else if block == "chain_command_block" {
							mode = packet.CommandBlockChain
						}
						cmd := item["Command"].(string)
						cusname := item["CustomName"].(string)
						exeft := item["ExecuteOnFirstTick"].(uint8)
						tickdelay := item["TickDelay"].(int32)
						aut := item["auto"].(uint8)
						trackoutput := item["TrackOutput"].(uint8)
						lo := item["LastOutput"].(string)
						//conditionalmode:=item["conditionalMode"].(uint8)
						data := item["data"].(int32)
						var conb bool
						if (data>>3)&1 == 1 {
							conb = true
						} else {
							conb = false
						}
						var exeftb bool
						if exeft == 0 {
							exeftb = true
						} else {
							exeftb = true
						}
						var tob bool
						if trackoutput == 1 {
							tob = true
						} else {
							tob = false
						}
						var nrb bool
						if aut == 1 {
							nrb = false
							//REVERSED!!
						} else {
							nrb = true
						}
						cbdata = &types.CommandBlockData{
							Mode:               mode,
							Command:            cmd,
							CustomName:         cusname,
							ExecuteOnFirstTick: exeftb,
							LastOutput:         lo,
							TickDelay:          tickdelay,
							TrackOutput:        tob,
							Conditional:        conb,
							NeedRedstone:       nrb,
						}
					}
					blocks[counter] = &types.RuntimeModule{
						BlockRuntimeId:   runtimeId,
						CommandBlockData: cbdata,
						ChestData:        chestData,
						Point: types.Position{
							X: x,
							Y: y,
							Z: z,
						},
					}
					counter++
				}
			}
		}
		world_provider.DestroyWorld()
		blocks = blocks[:counter]
		runtime.GC()
		out := bdump.BDump{
			Blocks: blocks,
			//Blocks: nil,
		}
		if strings.LastIndex(cfg.Path, ".bdx") != len(cfg.Path)-4 || len(cfg.Path) < 4 {
			cfg.Path += ".bdx"
		}
		command.Tellraw(conn, "EXPORT >> Writing output file")
		err, signerr := out.WriteToFile(cfg.Path)
		if err != nil {
			command.Tellraw(conn, fmt.Sprintf("EXPORT >> ERROR: Failed to export: %v", err))
			return
		} else if signerr != nil {
			command.Tellraw(conn, fmt.Sprintf("EXPORT >> Note: The file is unsigned since the following error was trapped: %v", signerr))
		} else {
			command.Tellraw(conn, fmt.Sprintf("EXPORT >> File signed successfully"))
		}
		command.Tellraw(conn, fmt.Sprintf("EXPORT >> Successfully exported your structure to %v", cfg.Path))
		runtime.GC()
	}()
	return nil
}
