package task_config

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type GUI struct {
	setContent   func(v fyne.CanvasObject)
	getContent   func() fyne.CanvasObject
	origContent  fyne.CanvasObject
	masterWindow fyne.Window

	content      fyne.CanvasObject
	majorContent fyne.CanvasObject
}

func New() *GUI {
	gui := &GUI{}
	gui.majorContent = gui.makeMajorContent()
	return gui
}

func (g *GUI) makeMajorContent() fyne.CanvasObject {
	return widget.NewLabel("这里还没有做")
}

func (g *GUI) GetContent(setContent func(v fyne.CanvasObject), getContent func() fyne.CanvasObject, masterWindow fyne.Window) fyne.CanvasObject {
	g.origContent = getContent()
	g.setContent = setContent
	g.getContent = getContent
	g.masterWindow = masterWindow
	g.content = container.NewBorder(nil, &widget.Button{
		Text: "取消",
		OnTapped: func() {
			g.setContent(g.origContent)
		},
		Icon:          theme.CancelIcon(),
		IconPlacement: widget.ButtonIconLeadingText,
	}, nil, nil, g.majorContent)

	return g.content
}
