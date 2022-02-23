package session

import (
	"phoenixbuilder_3rd_gui/fb/dragonfly/server/block/cube"
	"phoenixbuilder_3rd_gui/fb/minecraft/protocol/packet"
)

// BlockPickRequestHandler handles the BlockPickRequest packet.
type BlockPickRequestHandler struct{}

// Handle ...
func (b BlockPickRequestHandler) Handle(p packet.Packet, s *Session) error {
	pk := p.(*packet.BlockPickRequest)
	s.c.PickBlock(cube.Pos{int(pk.Position.X()), int(pk.Position.Y()), int(pk.Position.Z())})
	return nil
}
