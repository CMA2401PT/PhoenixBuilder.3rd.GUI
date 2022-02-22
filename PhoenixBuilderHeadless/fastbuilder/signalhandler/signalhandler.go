package signalhandler

import (
	"os"
	"os/signal"
	I18n "phoenixbuilder/fastbuilder/i18n"
	"phoenixbuilder/minecraft"
	bridge_fmt "phoenixbuilder/session/bridge/fmt"
	"syscall"
)

func Init(conn *minecraft.Conn) {
	go func() {
		signalchannel := make(chan os.Signal)
		signal.Notify(signalchannel, os.Interrupt) // ^C
		signal.Notify(signalchannel, syscall.SIGTERM)
		signal.Notify(signalchannel, syscall.SIGQUIT) // ^\
		<-signalchannel
		conn.Close()
		bridge_fmt.Printf("%s.\n", I18n.T(I18n.QuitCorrectly))
		os.Exit(0)
	}()
}
