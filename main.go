package main

import (
	"net/http"
	"phoenixbuilder_3rd_gui/gui/assets"
	"phoenixbuilder_3rd_gui/gui/global"
	"phoenixbuilder_3rd_gui/gui/profiles"
	my_theme "phoenixbuilder_3rd_gui/gui/theme"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

var topWindow fyne.Window
var appTheme *my_theme.MyTheme

func main() {
	app := app.NewWithID("gui.3rd.PhoenixBuilder")
	appStorage := app.Storage()
	//appStorage.Create("config.yaml")

	go func() {
		// popup a network permission dialog
		http.Get("http://captive.apple.com")
	}()

	appTheme = my_theme.NewTheme()
	setThemeChineseFont(appTheme)
	appTheme.SetLight()
	app.Settings().SetTheme(appTheme)

	topWindow = app.NewWindow("Fastbuilder.3rd.Gui")
	icon := canvas.NewImageFromResource(assets.ResourceIconPng)
	icon.FillMode = canvas.ImageFillContain
	app.SetIcon(icon.Resource)
	//iconRes, err := utils.LoadFromAssets("Icon", "Icon.png")
	//if err == nil {
	//	icon := canvas.NewImageFromResource(iconRes)
	//	icon.FillMode = canvas.ImageFillContain
	//	app.SetIcon(icon.Resource)
	//} else {
	//	dialog.ShowError(fmt.Errorf("无法加载图标：\n\n%v", err), topWindow)
	//}
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

	global.MakeThemeToggleBtn(app, appTheme)
	global.MakeInformPopButton(topWindow)
	// global.MakeDebugButton(app, setContent, getContent)
	global.MakeBannner("v0.0.4")

	//vsplit := container.NewVSplit(debugContent, majorContent)
	//vsplit.Offset = 0.05
	content := container.NewBorder(global.Banner, nil, nil, nil, majorContent)
	topWindow.SetContent(content)

	//onPanicFn := func(err error) {
	//	dialog.ShowError(fmt.Errorf("发生了严重错误，程序即将退出：\n\n%v", err), topWindow)
	//	// os.Exit(-1)
	//}

	//stat, err := os.Stat(dataFolder)
	//if !(err == nil && stat.IsDir()) {
	//	err = os.Mkdir(dataFolder, 0755)
	//	if err != nil {
	//		onPanicFn(fmt.Errorf("权限错误，无法创建必要的数据文件夹 %v (%v)", dataFolder, err))
	//	}
	//}

	profilesObject := profiles.New(appStorage)
	setContent(profilesObject.GetContent(setContent, getContent, topWindow, app))

	topWindow.Resize(fyne.NewSize(480, 640))
	topWindow.ShowAndRun()
}

func setThemeChineseFont(t *my_theme.MyTheme) {
	appTheme.Regular = assets.ResourceRegularFont
	appTheme.Italic = assets.ResourceRegularFont
	appTheme.Monospace = assets.ResourceRegularFont
	appTheme.Bold = assets.ResourceBoldFont
	appTheme.BoldItalic = assets.ResourceBoldFont
}
