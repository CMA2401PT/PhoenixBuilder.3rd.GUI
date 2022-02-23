package session

import (
	"fmt"
	"phoenixbuilder_3rd_gui/gui/profiles/config"
	"phoenixbuilder_3rd_gui/gui/profiles/session/list_terminal"
	"phoenixbuilder_3rd_gui/gui/profiles/session/tasks"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	bot_bridge_command "phoenixbuilder_3rd_gui/fb/fastbuilder/command"
	bot_session "phoenixbuilder_3rd_gui/fb/session"
	bot_bridge_fmt "phoenixbuilder_3rd_gui/fb/session/bridge/fmt"
)

type GUI struct {
	setContent   func(v fyne.CanvasObject)
	getContent   func() fyne.CanvasObject
	origContent  fyne.CanvasObject
	masterWindow fyne.Window

	writeBackConfigFn func()
	sessionConfig     *config.SessionConfigWithName
	term              *list_terminal.Terminal
	content           fyne.CanvasObject

	loadingBar                      *widget.ProgressBarInfinite
	loadinglabel                    *widget.Label
	loadingIndicator                *fyne.Container
	cmdInputBar                     *widget.Entry
	quitButton                      *widget.Button
	createFromTemplateBtn           *widget.Button
	titleRedirectBarHiderActivated  bool
	titleRedirectBar                *widget.Label
	titleRedirectBarLastUpdatedTime time.Time
	functionGroup                   *fyne.Container
	taskMenu                        *tasks.GUI

	alreadyClosed bool
	terminateChan chan string
	BotSession    *bot_session.Session
}

func New(config *config.SessionConfigWithName, writeBackConfigFn func()) *GUI {
	gui := &GUI{
		sessionConfig:     config,
		writeBackConfigFn: writeBackConfigFn,
	}
	return gui
}

func (g *GUI) setLoading(hint string) {
	g.functionGroup.Hide()
	g.loadingIndicator.Show()
	g.loadingBar.Start()
	g.loadinglabel.SetText(hint)
}

func (g *GUI) doneLoading() {
	g.functionGroup.Show()
	g.loadingIndicator.Hide()
	g.loadingBar.Stop()
}

func (g *GUI) closeGUI() {
	g.alreadyClosed = true
	g.setContent(g.origContent)
	g.BotSession.Stop()
}

func (g *GUI) sendCmd(s string) {
	s = strings.TrimSpace(s)
	fmt.Println("Cmd:", s)
	g.cmdInputBar.SetText("")
	g.term.AppendNewLine(s, true)
	g.BotSession.Execute(s)
}

func (g *GUI) redirectCliOutput(s string) {
	s = strings.TrimSpace(s)
	g.term.AppendNewLine(s, false)
}

func (g *GUI) redirectTitleDisplay(s string) {
	s = strings.TrimSpace(s)
	g.titleRedirectBar.Text = s
	g.titleRedirectBarLastUpdatedTime = time.Now()
	if g.titleRedirectBar.Hidden {
		g.titleRedirectBar.Show()
		if !g.titleRedirectBarHiderActivated {
			g.titleRedirectBarHiderActivated = true
			go func() {
				for {
					time.Sleep(3 * time.Second)
					if time.Since(g.titleRedirectBarLastUpdatedTime) > time.Second*3 {
						g.titleRedirectBar.Hide()
						g.titleRedirectBarHiderActivated = false
						break
					}
				}
			}()
		}
	}
	g.titleRedirectBar.Refresh()
}

func (g *GUI) onLoginError(err error) {
	dialog.NewError(err, g.masterWindow).Show()
	g.closeGUI()
}

func (g *GUI) onRuntimeError(err error) {
	dialog.NewError(err, g.masterWindow).Show()
	g.closeGUI()
}

func (g *GUI) makeToolContent() fyne.CanvasObject {
	g.loadingBar = widget.NewProgressBarInfinite()
	g.loadinglabel = widget.NewLabel("正在加载...")
	g.loadinglabel.Alignment = fyne.TextAlignCenter
	g.loadingIndicator = container.NewVBox(
		g.loadinglabel, g.loadingBar)
	g.cmdInputBar = widget.NewEntry()
	g.cmdInputBar.PlaceHolder = "输入/黏贴命令 (中文可能有问题)"
	g.cmdInputBar.OnSubmitted = func(s string) {
		g.sendCmd(s)
	}
	g.quitButton = widget.NewButton("结束会话", func() {
		g.closeGUI()
	})
	g.quitButton.Icon = theme.NavigateBackIcon()
	g.quitButton.IconPlacement = widget.ButtonIconLeadingText
	g.createFromTemplateBtn = widget.NewButton("可用命令", func() {})
	g.createFromTemplateBtn.Icon = theme.ContentAddIcon()
	g.createFromTemplateBtn.IconPlacement = widget.ButtonIconTrailingText
	g.createFromTemplateBtn.Importance = widget.HighImportance
	g.titleRedirectBar = widget.NewLabel("")
	g.titleRedirectBar.Wrapping = fyne.TextWrapWord
	g.titleRedirectBar.Hide()
	g.functionGroup = container.NewVBox(
		g.titleRedirectBar,
		container.NewBorder(nil, nil,
			g.quitButton, g.createFromTemplateBtn, g.cmdInputBar,
		),
	)

	g.functionGroup.Hide()
	return container.NewVBox(g.loadingIndicator, g.functionGroup)
}

func (g *GUI) GetContent(setContent func(v fyne.CanvasObject), getContent func() fyne.CanvasObject, masterWindow fyne.Window) fyne.CanvasObject {
	g.origContent = getContent()
	g.setContent = setContent
	g.getContent = getContent
	g.masterWindow = masterWindow
	g.term = list_terminal.New()
	g.term.OnPasteFn = func(s string) {
		g.cmdInputBar.SetText(s)
	}
	toolbar := g.makeToolContent()
	g.content = container.NewBorder(
		nil, toolbar, nil, nil,
		g.term.GetContent(g.masterWindow),
	)

	return g.content
}

func (g *GUI) AfterMount() {
	bot_bridge_fmt.HookFunc = g.redirectCliOutput
	bot_bridge_command.AdditionalChatCb = g.redirectCliOutput
	bot_bridge_command.AdditionalTitleCb = g.redirectTitleDisplay

	g.setLoading("正在登录，最长可能需要30s...")
	go func() {
		g.BotSession = bot_session.NewSession(g.sessionConfig.Config)
		if g.BotSession == nil {
			g.onLoginError(fmt.Errorf("一个现有会话未正常退出，或许你需要重启程序"))
			return
		}
		terminateChan, err := g.BotSession.Start()
		if err != nil {
			g.onLoginError(fmt.Errorf("无法顺利登陆到租赁服中\n%v", err))
			return
		}
		g.writeBackConfigFn()
		g.taskMenu = tasks.New(g.BotSession, g.sendCmd)
		g.createFromTemplateBtn.OnTapped = func() {
			g.setContent(g.taskMenu.GetContent(g.setContent, g.getContent, g.masterWindow))
		}
		g.terminateChan = terminateChan
		g.doneLoading()
		closeReason := <-g.terminateChan
		if !g.alreadyClosed {
			g.onRuntimeError(fmt.Errorf("和租赁服的连接被迫断开了\n%v", closeReason))
			return
		}
	}()
}
