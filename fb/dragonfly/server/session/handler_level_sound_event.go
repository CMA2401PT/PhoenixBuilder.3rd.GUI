package session

import (
	"phoenixbuilder_3rd_gui/fb/dragonfly/server/world"
	"phoenixbuilder_3rd_gui/fb/minecraft/protocol/packet"
)

type LevelSoundEventHandler struct{}

func (l LevelSoundEventHandler) Handle(p packet.Packet, s *Session) error {
	pk := p.(*packet.LevelSoundEvent)
	if pk.SoundType == packet.SoundEventAttackNoDamage && (s.c.GameMode() != world.GameModeSpectator{}) {
		s.swingingArm.Store(true)
		defer s.swingingArm.Store(false)
		s.c.PunchAir()
	}
	return nil
}
