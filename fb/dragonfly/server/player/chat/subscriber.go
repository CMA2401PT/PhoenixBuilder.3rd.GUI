package chat

import (
	"fmt"
	"phoenixbuilder_3rd_gui/fb/minecraft/text"
	bridge_fmt "phoenixbuilder_3rd_gui/fb/session/bridge/fmt"
	"strings"
)

// Subscriber represents an entity that may subscribe to a Chat. In order to do so, the Subscriber must
// implement methods to send messages to it.
type Subscriber interface {
	// Message sends a formatted message to the subscriber. The message is formatted as it would when using
	// fmt.Println.
	Message(a ...interface{})
}

// StdoutSubscriber is an implementation of Subscriber that forwards messages sent to the chat to the stdout.
type StdoutSubscriber struct{}

// Message ...
func (c StdoutSubscriber) Message(a ...interface{}) {
	s := make([]string, len(a))
	for i, b := range a {
		s[i] = fmt.Sprint(b)
	}
	t := text.ANSI(strings.Join(s, " "))
	if !strings.HasSuffix(t, "\n") {
		bridge_fmt.Println(t)
		return
	}
	bridge_fmt.Print(t)
}
