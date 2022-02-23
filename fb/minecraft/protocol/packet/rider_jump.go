package packet

import (
	"phoenixbuilder_3rd_gui/fb/minecraft/protocol"
)

// RiderJump is sent by the client to the server when it jumps while riding an entity that has the
// WASDControlled entity flag set, for example when riding a horse.
type RiderJump struct {
	// JumpStrength is the strength of the jump, depending on how long the rider has held the jump button.
	JumpStrength int32
}

// ID ...
func (*RiderJump) ID() uint32 {
	return IDRiderJump
}

// Marshal ...
func (pk *RiderJump) Marshal(w *protocol.Writer) {
	w.Varint32(&pk.JumpStrength)
}

// Unmarshal ...
func (pk *RiderJump) Unmarshal(r *protocol.Reader) {
	r.Varint32(&pk.JumpStrength)
}
