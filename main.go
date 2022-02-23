package main

import (
	"fmt"
	"os"
	"phoenixbuilder_3rd_gui/gui/profiles"
	my_theme "phoenixbuilder_3rd_gui/gui/theme"

	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

var topWindow fyne.Window

func main() {
	dataFolder := "./data"
	app := app.NewWithID("gui.3rd.PhoenixBuilder")
	chineseTheme := &my_theme.MyTheme{}
	chineseTheme.SetLight()
	chineseTheme.SetFonts("./assets/font/Consolas-with-Yahei Regular Nerd Font.ttf", "")
	app.Settings().SetTheme(chineseTheme)
	topWindow = app.NewWindow("Fastbuilder.3rd.Gui")
	icon := canvas.NewImageFromFile("Icon.png")
	icon.FillMode = canvas.ImageFillContain
	app.SetIcon(icon.Resource)
	topWindow.SetMaster()

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

	debugContent := makeDebugContent(app, setContent, getContent)
	vsplit := container.NewVSplit(debugContent, majorContent)
	vsplit.Offset = 0.05
	topWindow.SetContent(vsplit)

	onPanicFn := func(err error) {
		dialog.ShowError(fmt.Errorf("发生了严重错误，程序即将退出：\n\n%v", err), topWindow)
		os.Exit(-1)
	}

	stat, err := os.Stat(dataFolder)
	if !(err == nil && stat.IsDir()) {
		err = os.Mkdir(dataFolder, 0755)
		if err != nil {
			onPanicFn(fmt.Errorf("权限错误，无法创建必要的数据文件夹 %v (%v)", dataFolder, err))
		}
	}

	profilesObject := profiles.New(dataFolder)
	setContent(profilesObject.GetContent(setContent, getContent, topWindow))

	topWindow.Resize(fyne.NewSize(480, 640))
	topWindow.ShowAndRun()
}

func makeDebugContent(app fyne.App, setContent func(v fyne.CanvasObject), getContent func() fyne.CanvasObject) fyne.CanvasObject {
	content := container.NewVBox(
		container.New(layout.NewGridLayout(5),
			widget.NewLabel("WIP"),
			widget.NewButton("Dark", func() {
				t := app.Settings().Theme()
				t.(*my_theme.MyTheme).SetDark()
				app.Settings().SetTheme(t)
			}),
			widget.NewButton("Light", func() {
				t := app.Settings().Theme()
				t.(*my_theme.MyTheme).SetLight()
				app.Settings().SetTheme(t)
			}),
		),
	)
	return content
}
