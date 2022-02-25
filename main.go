package main

import (
	"fmt"
	"net/http"
	"phoenixbuilder_3rd_gui/gui/assets"
	"phoenixbuilder_3rd_gui/gui/global"
	"phoenixbuilder_3rd_gui/gui/profiles"
	my_theme "phoenixbuilder_3rd_gui/gui/theme"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
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
	global.MakeThemeToggleBtn(app, appTheme)
	global.MakeInformPopButton(topWindow)
	global.MakeBannner("v0.0.3")

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
	debugContent.Hide()
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

func makeDebugContent(app fyne.App, setContent func(v fyne.CanvasObject), getContent func() fyne.CanvasObject) fyne.CanvasObject {
	content := container.NewVBox(
		container.New(layout.NewGridLayout(5),
			widget.NewLabel("WIP"),
			widget.NewButton("Dark", func() {
				appTheme.SetDark()
				app.Settings().SetTheme(appTheme)
			}),
			widget.NewButton("Light", func() {
				appTheme.SetLight()
				app.Settings().SetTheme(appTheme)
			}),
			widget.NewButton("Chinese", func() {
				//onError := func(info error) {
				//	dialog.ShowError(info, topWindow)
				//	time.Sleep(5 * time.Second)
				//}
				//
				//res, err := utils.LoadFromAssets("Consolas_with_Yahei_Regular.ttf", "Consolas_with_Yahei_Regular.ttf")
				//if err != nil {
				//	onError(err)
				//	return
				//}
				appTheme.Regular = assets.ResourceRegularFont
				appTheme.Italic = assets.ResourceRegularFont
				appTheme.Monospace = assets.ResourceRegularFont
				//res, err = utils.LoadFromAssets("Consolas_with_Yahei_Bold.ttf", "Consolas_with_Yahei_Bold.ttf")
				//if err != nil {
				//	onError(err)
				//	return
				//}
				appTheme.Bold = assets.ResourceBoldFont
				appTheme.BoldItalic = assets.ResourceBoldFont
				//chineseTheme.SetFontsFromAssets("Consolas_with_Yahei_Regular.ttf", "", onError)
				app.Settings().SetTheme(appTheme)
			}),
			widget.NewButton("File", func() {
				dialog.NewFileOpen(func(closer fyne.URIReadCloser, err error) {
					if err != nil {
						dialog.ShowError(err, topWindow)
					} else {
						dialog.ShowInformation("Selected", closer.URI().String(), topWindow)
					}
				}, topWindow).Show()
			}),
			widget.NewButton("RootStorage", func() {
				dialog.ShowInformation("RootStorage", app.Storage().RootURI().String(), topWindow)
			}),
			widget.NewButton("ListRootStorage", func() {
				dialog.ShowInformation("ListRootStorage", fmt.Sprintf("%v", app.Storage().List()), topWindow)
				//appStorage.List()
			}),
			widget.NewButton("Remove Config", func() {
				err := app.Storage().Remove("config.yaml")
				if err != nil {
					dialog.ShowInformation("Cannot Remove", fmt.Sprintf("%v\n%v", app.Storage().List(), err), topWindow)
				}
			}),
			widget.NewButton("Save Config", func() {
				_, err := app.Storage().Save("config.yaml")
				if err != nil {
					dialog.ShowInformation("Cannot Save", fmt.Sprintf("%v\n%v", app.Storage().List(), err), topWindow)
				}
			}),
			widget.NewButton("File&os.Open", func() {
				dialog.NewFileOpen(func(closer fyne.URIReadCloser, err error) {
					if err != nil {
						dialog.ShowError(err, topWindow)
					} else {
						dialog.ShowInformation("Selected", closer.URI().Extension(), topWindow)
						p := closer.URI().Extension()
						//p = closer.URI().Path()
						cp := p
						//cp = strings.TrimPrefix(cp, "content://")
						//_, err := os.Open(cp)
						closer.Close()
						if err != nil {
							//fyne.Storage.Open()
							dialog.ShowError(fmt.Errorf("os.Open error\n%v\n%v", cp, err), topWindow)
						}
					}
				}, topWindow).Show()
			}),
		),
	)
	return content
}
