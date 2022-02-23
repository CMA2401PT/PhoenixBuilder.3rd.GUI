package session

import (
	"bytes"
	"phoenixbuilder_3rd_gui/fb/dragonfly/server/player/scoreboard"
	"github.com/google/uuid"
	"phoenixbuilder_3rd_gui/fb/minecraft/protocol"
	"phoenixbuilder_3rd_gui/fb/minecraft/protocol/packet"
	"strings"
	"time"
)

// SendMessage ...
func (s *Session) SendMessage(message string) {
	s.writePacket(&packet.Text{
		TextType: packet.TextTypeRaw,
		Message:  message,
	})
}

// SendTip ...
func (s *Session) SendTip(message string) {
	s.writePacket(&packet.Text{
		TextType: packet.TextTypeTip,
		Message:  message,
	})
}

// SendAnnouncement ...
func (s *Session) SendAnnouncement(message string) {
	s.writePacket(&packet.Text{
		TextType: packet.TextTypeAnnouncement,
		Message:  message,
	})
}

// SendPopup ...
func (s *Session) SendPopup(message string) {
	s.writePacket(&packet.Text{
		TextType: packet.TextTypePopup,
		Message:  message,
	})
}

// SendJukeboxPopup ...
func (s *Session) SendJukeboxPopup(message string) {
	s.writePacket(&packet.Text{
		TextType: packet.TextTypeJukeboxPopup,
		Message:  message,
	})
}

// SendScoreboard ...
func (s *Session) SendScoreboard(sb *scoreboard.Scoreboard) {
	if s.scoreboardObj.Load() != "" {
		s.RemoveScoreboard()
	}
	obj := uuid.New().String()
	s.scoreboardObj.Store(obj)

	s.writePacket(&packet.SetDisplayObjective{
		DisplaySlot:   "sidebar",
		ObjectiveName: obj,
		DisplayName:   sb.Name(),
		CriteriaName:  "dummy",
	})
	pk := &packet.SetScore{
		ActionType: packet.ScoreboardActionModify,
	}
	for k, line := range bytes.Split(sb.Bytes(), []byte{'\n'}) {
		if len(line) == 0 {
			line = []byte("§" + colours[k])
		}
		pk.Entries = append(pk.Entries, protocol.ScoreboardEntry{
			EntryID:       int64(k),
			ObjectiveName: s.scoreboardObj.Load(),
			Score:         int32(k),
			IdentityType:  protocol.ScoreboardIdentityFakePlayer,
			DisplayName:   padScoreboardString(sb, string(line)),
		})
	}
	s.writePacket(pk)
}

// pad pads the string passed for as much as needed to achieve the same length as the name of the scoreboard.
// If the string passed is already of the same length as the name of the scoreboard or longer, the string will
// receive one space of padding.
func padScoreboardString(sb *scoreboard.Scoreboard, s string) string {
	if len(sb.Name())-len(s)-2 <= 0 {
		return " " + s + " "
	}
	return " " + s + strings.Repeat(" ", len(sb.Name())-len(s)-2)
}

// colours holds a list of colour codes to be filled out for empty lines in a scoreboard.
var colours = [15]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "a", "b", "c", "d", "e", "f"}

// RemoveScoreboard ...
func (s *Session) RemoveScoreboard() {
	s.writePacket(&packet.RemoveObjective{
		ObjectiveName: s.scoreboardObj.Load(),
	})
}

// SendBossBar sends a boss bar to the player with the text passed and the health percentage of the bar.
// SendBossBar removes any boss bar that might be active before sending the new one.
func (s *Session) SendBossBar(text string, healthPercentage float64) {
	s.RemoveBossBar()
	s.writePacket(&packet.BossEvent{
		BossEntityUniqueID: selfEntityRuntimeID,
		EventType:          packet.BossEventShow,
		BossBarTitle:       text,
		HealthPercentage:   float32(healthPercentage),
	})
}

// RemoveBossBar removes any boss bar currently active on the player's screen.
func (s *Session) RemoveBossBar() {
	s.writePacket(&packet.BossEvent{
		BossEntityUniqueID: selfEntityRuntimeID,
		EventType:          packet.BossEventHide,
	})
}

const tickLength = time.Second / 20

// SetTitleDurations ...
func (s *Session) SetTitleDurations(fadeInDuration, remainDuration, fadeOutDuration time.Duration) {
	s.writePacket(&packet.SetTitle{
		ActionType:      packet.TitleActionSetDurations,
		FadeInDuration:  int32(fadeInDuration / tickLength),
		RemainDuration:  int32(remainDuration / tickLength),
		FadeOutDuration: int32(fadeOutDuration / tickLength),
	})
}

// SendTitle ...
func (s *Session) SendTitle(text string) {
	s.writePacket(&packet.SetTitle{ActionType: packet.TitleActionSetTitle, Text: text})
}

// SendSubtitle ...
func (s *Session) SendSubtitle(text string) {
	s.writePacket(&packet.SetTitle{ActionType: packet.TitleActionSetSubtitle, Text: text})
}

// SendActionBarMessage ...
func (s *Session) SendActionBarMessage(text string) {
	s.writePacket(&packet.SetTitle{ActionType: packet.TitleActionSetActionBar, Text: text})
}
