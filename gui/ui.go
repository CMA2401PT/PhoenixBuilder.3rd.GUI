package gui

import (
	"fmt"
	"os"
	"phoenixbuilder_3rd_gui/gui/profiles"

	"fyne.io/fyne/v2/dialog"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type GUI struct {
	app          fyne.App
	masterWindow fyne.Window
	dataFolder   string
}

func NewGUI(dataFolder string) *GUI {
	return &GUI{
		dataFolder: dataFolder,
	}
}

func (g *GUI) Run() {
	g.app = app.NewWithID("gui.3rd.PhoenixBuilder")
	chineseTheme := &myTheme{}
	chineseTheme.SetLight()
	chineseTheme.SetFonts("./assets/font/Consolas-with-Yahei Regular Nerd Font.ttf", "")
	g.app.Settings().SetTheme(chineseTheme)
	g.masterWindow = g.app.NewWindow("Fastbuilder.3rd.Gui")
	icon := canvas.NewImageFromFile("Icon.png")
	icon.FillMode = canvas.ImageFillContain
	g.app.SetIcon(icon.Resource)
	g.masterWindow.SetMaster()

	majorContent := container.NewMax()

	getContent := func() fyne.CanvasObject {
		if len(majorContent.Objects) != 0 {
			return majorContent.Objects[0]
		} else {
			return nil
		}
	}
	setContent := func(v fyne.CanvasObject) {
		majorContent.Objects = []fyne.CanvasObject{v}
		majorContent.Refresh()
	}

	debugContent := makeDebugContent(g, setContent, getContent)
	vsplit := container.NewVSplit(debugContent, majorContent)
	vsplit.Offset = 0.05
	g.masterWindow.SetContent(vsplit)

	onPanicFn := func(err error) {
		dialog.ShowError(fmt.Errorf("发生了严重错误，程序即将退出：\n\n%v", err), g.masterWindow)
		os.Exit(-1)
	}

	stat, err := os.Stat(g.dataFolder)
	if !(err == nil && stat.IsDir()) {
		err = os.Mkdir(g.dataFolder, 0755)
		if err != nil {
			onPanicFn(fmt.Errorf("权限错误，无法创建必要的数据文件夹 %v (%v)", g.dataFolder, err))
		}
	}

	profilesObject := profiles.New(g.dataFolder)
	setContent(profilesObject.GetContent(setContent, getContent, g.masterWindow))

	g.masterWindow.Resize(fyne.NewSize(480, 640))
	g.masterWindow.ShowAndRun()
}

func makeDebugContent(g *GUI, setContent func(v fyne.CanvasObject), getContent func() fyne.CanvasObject) fyne.CanvasObject {
	content := container.NewVBox(
		container.New(layout.NewGridLayout(5),
			widget.NewLabel("WIP"),
			widget.NewButton("Dark", func() {
				t := g.app.Settings().Theme()
				t.(*myTheme).SetDark()
				g.app.Settings().SetTheme(t)
			}),
			widget.NewButton("Light", func() {
				t := g.app.Settings().Theme()
				t.(*myTheme).SetLight()
				g.app.Settings().SetTheme(t)
			}),
		),
	)
	return content
}
