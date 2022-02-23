package command

import (
	"fmt"
	"phoenixbuilder_3rd_gui/fb/minecraft"
	"phoenixbuilder_3rd_gui/fb/minecraft/protocol"
	"phoenixbuilder_3rd_gui/fb/minecraft/protocol/packet"
	"sync"

	"github.com/google/uuid"
)

var UUIDMap sync.Map //= make(map[string]func(*minecraft.Conn,*[]protocol.CommandOutputMessage))
var BlockUpdateSubscribeMap sync.Map
var AdditionalChatCb func(string)

func init() {
	AdditionalChatCb = func(s string) {}
}

func ClearUUIDMap() {
	UUIDMap = sync.Map{}
}

func SendCommand(command string, UUID uuid.UUID, conn *minecraft.Conn) error {
	requestId, _ := uuid.Parse("96045347-a6a3-4114-94c0-1bc4cc561694")
	origin := protocol.CommandOrigin{
		Origin:         protocol.CommandOriginPlayer,
		UUID:           UUID,
		RequestID:      requestId.String(),
		PlayerUniqueID: 0,
	}
	commandRequest := &packet.CommandRequest{
		CommandLine:   command,
		CommandOrigin: origin,
		Internal:      false,
		UnLimited:     false,
	}
	return conn.WritePacket(commandRequest)
}

func SendWSCommand(command string, UUID uuid.UUID, conn *minecraft.Conn) error {
	requestId, _ := uuid.Parse("96045347-a6a3-4114-94c0-1bc4cc561694")
	origin := protocol.CommandOrigin{
		Origin:         protocol.CommandOriginAutomationPlayer,
		UUID:           UUID,
		RequestID:      requestId.String(),
		PlayerUniqueID: 0,
	}
	commandRequest := &packet.CommandRequest{
		CommandLine:   command,
		CommandOrigin: origin,
		Internal:      false,
		UnLimited:     false,
	}
	return conn.WritePacket(commandRequest)
}

func SendSizukanaCommand(command string, conn *minecraft.Conn) error {
	return conn.WritePacket(&packet.SettingsCommand{
		CommandLine:    command,
		SuppressOutput: true,
	})
}

func SendChat(content string, conn *minecraft.Conn) error {
	AdditionalChatCb(content)
	idd := conn.IdentityData()
	return conn.WritePacket(&packet.Text{
		TextType:         packet.TextTypeChat,
		NeedsTranslation: false,
		SourceName:       idd.DisplayName,
		Message:          content,
		XUID:             idd.XUID,
		PlayerRuntimeID:  fmt.Sprintf("%d", conn.GameData().EntityUniqueID),
	})
}
