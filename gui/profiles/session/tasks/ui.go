package tasks

import (
	"fmt"
	"log"
	bot_session "phoenixbuilder_3rd_gui/fb/session"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type GUI struct {
	setContent   func(v fyne.CanvasObject)
	getContent   func() fyne.CanvasObject
	origContent  fyne.CanvasObject
	masterWindow fyne.Window

	BotSession *bot_session.Session
	sendCmdFn  func(string)
	startPos   *PosWidget
	endPos     *PosWidget
	// every ui element
	content fyne.CanvasObject
	// every ui element except return btn and two pos buttons
	majorContent        fyne.CanvasObject
	addMonkeyPathReader func(path string, fp fyne.URIReadCloser)
	addMonkeyPathWriter func(path string, fp fyne.URIWriteCloser)
}

func New(BotSession *bot_session.Session, sendCmdFn func(string), addMonkeyPathReader func(path string, fp fyne.URIReadCloser), addMonkeyPathWriter func(path string, fp fyne.URIWriteCloser)) *GUI {
	gui := &GUI{
		BotSession:          BotSession,
		sendCmdFn:           sendCmdFn,
		addMonkeyPathReader: addMonkeyPathReader,
		addMonkeyPathWriter: addMonkeyPathWriter,
	}
	gui.makePosWidgets()
	gui.majorContent = gui.makeMajorContent()
	return gui
}

type PosWidget struct {
	PosContent fyne.CanvasObject
	UpdateBtn  *widget.Button
	WX, WY, WZ *widget.Entry
	dX, dY, dZ binding.Int
}

func NewPosWidget(x, y, z int, btn *widget.Button) *PosWidget {
	w := &PosWidget{}
	w.dX = binding.BindInt(&x)
	w.dY = binding.BindInt(&y)
	w.dZ = binding.BindInt(&z)

	w.WX = widget.NewEntryWithData(binding.IntToString(w.dX))
	w.WY = widget.NewEntryWithData(binding.IntToString(w.dY))
	w.WZ = widget.NewEntryWithData(binding.IntToString(w.dZ))
	w.PosContent = container.NewGridWithColumns(3, w.WX, w.WY, w.WZ)
	w.UpdateBtn = btn
	return w
}

func (pw *PosWidget) GetPos() (x, y, z int, err error) {
	x, err = pw.dX.Get()
	if err != nil {
		return
	}
	y, err = pw.dY.Get()
	if err != nil {
		return
	}
	z, err = pw.dZ.Get()
	return
}

func (pw *PosWidget) SetPos(x, y, z int) {
	pw.dX.Set(x)
	pw.dY.Set(y)
	pw.dZ.Set(z)
}

func (g *GUI) makePosWidgets() {
	x, y, z := g.BotSession.GetPos()
	startPos := NewPosWidget(x, y, z, &widget.Button{
		Text:       "使用[" + g.BotSession.Config.RespondUser + "]的位置",
		OnTapped:   func() { g.sendCmdFn("get") },
		Importance: widget.HighImportance,
	})
	ex, ey, ez := g.BotSession.GetEndPos()
	endPos := NewPosWidget(ex, ey, ez, &widget.Button{
		Text:       "使用[" + g.BotSession.Config.RespondUser + "]的位置",
		OnTapped:   func() { g.sendCmdFn("get end") },
		Importance: widget.HighImportance,
	})
	g.BotSession.CmdSetCbFn = func(x, y, z int) {
		startPos.SetPos(x, y, z)
	}
	g.BotSession.CmdSetEndCbFn = func(x, y, z int) {
		endPos.SetPos(x, y, z)
	}
	g.startPos = startPos
	g.endPos = endPos
}

func (g *GUI) sendCmdAndClose(cmd string) {
	g.sendCmdFn(cmd)
	g.setContent(g.origContent)
}

func (g *GUI) makeIntEntry(v int, name string, hint string) (*widget.FormItem, func() (int, error)) {
	cv := v
	bv := binding.BindInt(&cv)
	getter := func() (int, error) {
		gv, err := bv.Get()
		if err != nil {
			err = fmt.Errorf("%v数据错误\n%v", name, err)
			dialog.NewError(err, g.masterWindow).Show()
		}
		return gv, err
	}
	return &widget.FormItem{Text: name, Widget: widget.NewEntryWithData(binding.IntToString(bv)), HintText: hint}, getter
}

func (g *GUI) makeSelectEntry(options []string, name string, hint string) (*widget.FormItem, func() (string, error)) {
	coptions := make([]string, len(options))
	copy(coptions, options)
	w := widget.NewSelectEntry(coptions)
	w.SetText(options[0])
	getter := func() (string, error) {
		v := w.Text
		for _, o := range coptions {
			if o == v {
				return v, nil
			}
		}
		dialog.NewError(fmt.Errorf("%v选项错误\n%v不是可选项\n可选项为%v", name, v, coptions), g.masterWindow).Show()
		return "", fmt.Errorf("%v选项错误", hint)
	}
	return &widget.FormItem{Text: name, Widget: w, HintText: hint}, getter
}

func (g *GUI) makeTranslateSelectEntry(translateOptions []string, options []string, name string, hint string) (*widget.FormItem, func() (string, error)) {
	ctransOptions := make([]string, len(translateOptions))
	copy(ctransOptions, translateOptions)
	coptions := make([]string, len(options))
	copy(coptions, options)
	if len(translateOptions) != len(options) {
		panic("内部错误：翻译选项与选项数量不一致")
	}
	w := widget.NewSelectEntry(ctransOptions)
	w.SetText(ctransOptions[0])
	getter := func() (string, error) {
		v := w.Text
		for i, o := range ctransOptions {
			if o == v {
				return coptions[i], nil
			}
		}
		dialog.NewError(fmt.Errorf("%v选项错误\n%v不是可选项\n可选项为%v", name, v, coptions), g.masterWindow).Show()
		return "", fmt.Errorf("%v选项错误", hint)
	}
	return &widget.FormItem{Text: name, Widget: w, HintText: hint}, getter
}

func (g *GUI) makeStringEntry(v string, name string, hint string) (*widget.FormItem, func() (string, error)) {
	cv := v
	bv := binding.BindString(&cv)
	getter := func() (string, error) {
		gv, err := bv.Get()
		gv = strings.TrimSpace(gv)

		if err != nil {
			err = fmt.Errorf("%v数据错误\n%v", name, err)
			dialog.NewError(err, g.masterWindow).Show()
		} else if gv == "" {
			err = fmt.Errorf("%v不能为空", name, err)
			dialog.NewError(err, g.masterWindow).Show()
		}
		return gv, err
	}
	return &widget.FormItem{Text: name, Widget: widget.NewEntryWithData(bv), HintText: hint}, getter
}

func (g *GUI) makeBoolEntry(b bool, name string, hint string) (*widget.FormItem, func() (bool, error)) {
	cb := b
	bv := binding.BindBool(&cb)
	getter := func() (bool, error) {
		gv, err := bv.Get()
		if err != nil {
			err = fmt.Errorf("%v数据错误\n%v", name, err)
			dialog.NewError(err, g.masterWindow).Show()
		}
		return gv, err
	}
	return &widget.FormItem{Text: name, Widget: widget.NewCheckWithData("是", bv), HintText: hint}, getter
}

func (g *GUI) makeBoolOption(b bool, description string) (fyne.CanvasObject, func() (bool, error)) {
	cb := b
	bv := binding.BindBool(&cb)
	getter := func() (bool, error) {
		gv, err := bv.Get()
		if err != nil {
			err = fmt.Errorf("%v数据错误\n%v", description, err)
			dialog.NewError(err, g.masterWindow).Show()
		}
		return gv, err
	}
	return widget.NewCheckWithData(description, bv), getter
}

func (g *GUI) makeReadPathOption(description string, placeHolderStr string, filter []string) (fyne.CanvasObject, func() (string, fyne.URIReadCloser, error)) {
	filePath := ""
	bv := binding.BindString(&filePath)
	var fp fyne.URIReadCloser
	fileNameEntry := widget.NewEntryWithData(bv)
	fileNameEntry.SetPlaceHolder(placeHolderStr)

	getter := func() (string, fyne.URIReadCloser, error) {
		gv, err := bv.Get()
		if fp == nil {
			err = fmt.Errorf("%v必须选择文件", description)
			dialog.NewError(err, g.masterWindow).Show()
			return "", nil, err
		}
		if err != nil {
			err = fmt.Errorf("%v数据错误\n%v", description, err)
			dialog.NewError(err, g.masterWindow).Show()
			return "", nil, err
		}
		//for _, f := range filter {
		//	if strings.HasSuffix(gv, f) {
		//		return gv, fp, nil
		//	}
		//}
		//err = fmt.Errorf("%v\n不具有后缀\n%v", gv, filter)
		//dialog.NewError(err, g.masterWindow).Show()
		return gv, fp, nil
		//return "", nil, err
	}
	return container.NewBorder(nil, nil, nil, &widget.Button{
		Text: description,
		OnTapped: func() {
			fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
				if err != nil {
					dialog.ShowError(err, g.masterWindow)
					return
				}
				if reader == nil {
					log.Println("Cancelled")
					return
				}
				//fake path string
				bv.Set(reader.URI().String() + reader.URI().Extension())
				fp = reader
			}, g.masterWindow)
			fd.SetFilter(storage.NewExtensionFileFilter(filter))
			fd.Show()
		},
	}, fileNameEntry), getter
}

func (g *GUI) makeWritePathOption(description string, placeHolderStr string, filter []string) (fyne.CanvasObject, func() (string, fyne.URIWriteCloser, error)) {
	filePath := ""
	bv := binding.BindString(&filePath)
	var fp fyne.URIWriteCloser
	fileNameEntry := widget.NewEntryWithData(bv)
	fileNameEntry.SetPlaceHolder(placeHolderStr)

	getter := func() (string, fyne.URIWriteCloser, error) {
		gv, err := bv.Get()
		if fp == nil {
			err = fmt.Errorf("%v必须选择文件", description)
			dialog.NewError(err, g.masterWindow).Show()
			return "", nil, err
		}
		if err != nil {
			err = fmt.Errorf("%v数据错误\n%v", description, err)
			dialog.NewError(err, g.masterWindow).Show()
			return "", nil, err
		}
		//for _, f := range filter {
		//	if strings.HasSuffix(gv, f) {
		//		return gv, fp, nil
		//	}
		//}
		//err = fmt.Errorf("%v\n不具有后缀\n%v", gv, filter)
		//dialog.NewError(err, g.masterWindow).Show()
		return gv, fp, nil
		//return "", nil, err
	}
	return container.NewBorder(nil, nil, nil, &widget.Button{
		Text: description,
		OnTapped: func() {
			fd := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
				if err != nil {
					dialog.ShowError(err, g.masterWindow)
					return
				}
				if writer == nil {
					log.Println("Cancelled")
					return
				}
				//fake path string
				bv.Set(writer.URI().String() + writer.URI().Extension())
				fp = writer
			}, g.masterWindow)
			fd.SetFilter(storage.NewExtensionFileFilter(filter))
			fd.Show()
		},
	}, fileNameEntry), getter
}

func (g *GUI) setStartPos() error {
	x, y, z, err := g.startPos.GetPos()
	if err != nil {
		dialog.NewError(fmt.Errorf("坐标错误\n%v", err), g.masterWindow).Show()
		return err
	}
	g.sendCmdFn(fmt.Sprintf("set %d %d %d", x, y, z))
	return nil
}

func (g *GUI) setEndPos() error {
	x, y, z, err := g.endPos.GetPos()
	if err != nil {
		dialog.NewError(fmt.Errorf("坐标错误\n%v", err), g.masterWindow).Show()
		return err
	}
	g.sendCmdFn(fmt.Sprintf("setend %d %d %d", x, y, z))
	return nil
}

func (g *GUI) makeConfirmButton(hint string, onTapped func()) *widget.Button {
	return &widget.Button{
		Text:          hint,
		Icon:          theme.ConfirmIcon(),
		IconPlacement: widget.ButtonIconTrailingText,
		Importance:    widget.HighImportance,
		OnTapped:      onTapped,
	}
}

func (g *GUI) makeGeoCmdContent() fyne.CanvasObject {
	rund_circleFormItem, rund_circleGetter := g.makeTranslateSelectEntry([]string{"圆面", "圈"}, []string{"round", "circle"}, "目标(圆面/圈):", "圈只有外面一圈")
	radiusFormItem, radiusGet := g.makeIntEntry(0, "半径", "圆或圈的半径")
	facingFormItem, facingGet := g.makeSelectEntry([]string{"y", "x", "z"}, "朝向", "可选值有x,y,z")
	heightFormItem, heightGet := g.makeIntEntry(0, "高度", "")
	lengthFormItem, lengthGet := g.makeIntEntry(0, "长度", "")
	widthFormItem, widthGet := g.makeIntEntry(0, "宽度", "")
	blockFormItem, blockGet := g.makeStringEntry("air", "方块", "方块名称")
	blockdataFormItem, blockdataGet := g.makeIntEntry(0, "值", "方块特殊值")
	shpere_shapeFormItem, shpere_shapeGet := g.makeTranslateSelectEntry([]string{"空心", "实心"}, []string{"hollow", "solid"}, "球填充", "空心/实心")

	c := container.NewDocTabs(
		&container.TabItem{
			Text: "圆面/圈",
			Content: container.NewVBox(
				widget.NewForm(
					rund_circleFormItem,
					radiusFormItem,
					facingFormItem,
					// heightFormItem, this doesn't work
					blockFormItem,
					blockdataFormItem,
				),
				container.NewGridWithColumns(2, widget.NewLabel("圆心位置"), g.startPos.UpdateBtn),
				g.startPos.PosContent,
				g.makeConfirmButton("绘制", func() {
					target, err := rund_circleGetter()
					if err != nil {
						return
					}
					radius, err := radiusGet()
					if err != nil {
						return
					}
					facing, err := facingGet()
					if err != nil {
						return
					}
					block, err := blockGet()
					if err != nil {
						return
					}
					blockData, err := blockdataGet()
					if err != nil {
						return
					}
					err = g.setStartPos()
					if err != nil {
						return
					}
					g.sendCmdAndClose(fmt.Sprintf("%v -r %v -f %v -h 1 -b %v -d %v", target, radius, facing, block, blockData))
				}),
			),
		},
		&container.TabItem{
			Text: "球",
			Content: container.NewVBox(widget.NewForm(
				radiusFormItem,
				shpere_shapeFormItem),
				container.NewGridWithColumns(2, widget.NewLabel("球心位置"), g.startPos.UpdateBtn),
				g.startPos.PosContent,
				g.makeConfirmButton("绘制", func() {
					radius, err := radiusGet()
					if err != nil {
						return
					}
					shape, err := shpere_shapeGet()
					if err != nil {
						return
					}
					err = g.setStartPos()
					if err != nil {
						return
					}
					g.sendCmdAndClose(fmt.Sprintf("sphere -r %v -s %v", radius, shape))
				}),
			),
		},
		&container.TabItem{
			Text: "椭圆",
			Content: container.NewVBox(widget.NewForm(
				lengthFormItem,
				widthFormItem,
				facingFormItem),
				container.NewGridWithColumns(2, widget.NewLabel("圆心位置"), g.startPos.UpdateBtn),
				g.startPos.PosContent,
				g.makeConfirmButton("绘制", func() {
					length, err := lengthGet()
					if err != nil {
						return
					}
					width, err := widthGet()
					if err != nil {
						return
					}
					facing, err := facingGet()
					if err != nil {
						return
					}
					err = g.setStartPos()
					if err != nil {
						return
					}
					g.sendCmdAndClose(fmt.Sprintf("ellipse -l %v -w %v -f %v", length, width, facing))
				}),
			),
		},
		&container.TabItem{
			Text: "椭球",
			Content: container.NewVBox(widget.NewForm(
				lengthFormItem,
				widthFormItem,
				heightFormItem),
				container.NewGridWithColumns(2, widget.NewLabel("球心位置"), g.startPos.UpdateBtn),
				g.startPos.PosContent,
				g.makeConfirmButton("绘制", func() {
					length, err := lengthGet()
					if err != nil {
						return
					}
					width, err := widthGet()
					if err != nil {
						return
					}
					height, err := heightGet()
					if err != nil {
						return
					}
					err = g.setStartPos()
					if err != nil {
						return
					}
					g.sendCmdAndClose(fmt.Sprintf("ellipsoid -l %v -w %v -h %v", length, width, height))
				}),
			),
		},
	)
	return c
}

func (g *GUI) makeBuildingContent() fyne.CanvasObject {
	excludecommandsOption, excludecommandsGet := g.makeBoolOption(false, "不导入命令方块中的命令")
	invalidatecommandsOption, invalidateCommandsGet := g.makeBoolOption(false, "导入，但无效化命令方块中的命令")
	strictOption, strictGet := g.makeBoolOption(true, "验证文件签名")
	pathOption, pathGet := g.makeReadPathOption("选择建筑文件", ".schematic/.bdx/.mcacblock", []string{".schematic", ".bdx", ".mcacblock"})
	return container.NewVBox(
		pathOption,
		excludecommandsOption,
		invalidatecommandsOption,
		strictOption,
		container.NewGridWithColumns(2, widget.NewLabel("建筑起点位置"), g.startPos.UpdateBtn),
		g.startPos.PosContent,
		g.makeConfirmButton("导入", func() {
			path, fp, err := pathGet()
			if err != nil {
				return
			}
			flags := make([]string, 0)
			excludecommands, err := excludecommandsGet()
			if err != nil {
				return
			}
			if excludecommands {
				flags = append(flags, "--excludecommands")
			}
			invalidatecommands, err := invalidateCommandsGet()
			if err != nil {
				return
			}
			if invalidatecommands {
				flags = append(flags, "--invalidatecommands")
			}
			strict, err := strictGet()
			if err != nil {
				return
			}
			if strict {
				flags = append(flags, "--strict")
			}
			flagStr := strings.Join(flags, " ")
			err = g.setStartPos()
			if err != nil {
				return
			}
			//path = strings.TrimPrefix(path, "file://")
			cmd := path + " " + flagStr
			if strings.HasSuffix(path, "schematic") {
				cmd = "schem -p " + cmd
			} else if strings.HasSuffix(path, "mcacblock") {
				cmd = "acme -p " + cmd
			} else if strings.HasSuffix(path, "bdx") {
				cmd = "bdump -p " + cmd
			}
			g.addMonkeyPathReader(path, fp)
			g.sendCmdAndClose(cmd)
		}),
	)
}

func (g *GUI) makePlotContent() fyne.CanvasObject {
	pathOption, pathGet := g.makeReadPathOption("选择图片", "png/jpg", []string{".png", ".PNG", ".jpg", ".jpeg", ".JPG"})
	facingFormItem, facingGet := g.makeSelectEntry([]string{"y", "x", "z"}, "朝向", "可选值有x,y,z")
	mapXFormItem, mapXGet := g.makeIntEntry(1, "横向", "横向由几张地图构成")
	mapZFormItem, mapZGet := g.makeIntEntry(1, "纵向", "纵向由几张地图构成")
	mapYFormItem, mapYGet := g.makeIntEntry(0, "允许使用高度", ">40时通过阴影产生更多颜色")
	c := container.NewDocTabs(
		&container.TabItem{
			Text: "图片",
			Content: container.NewVBox(
				pathOption,
				widget.NewForm(facingFormItem),
				container.NewGridWithColumns(2, widget.NewLabel("图片绘制起点(建议为64的奇数倍)"), g.startPos.UpdateBtn),
				g.startPos.PosContent,
				g.makeConfirmButton("制图", func() {
					path, fp, err := pathGet()
					if err != nil {
						return
					}
					facing, err := facingGet()
					if err != nil {
						return
					}
					err = g.setStartPos()
					if err != nil {
						return
					}
					g.addMonkeyPathReader(path, fp)
					g.sendCmdAndClose(fmt.Sprintf("plot -p %v -f %v", path, facing))
				}),
			),
		},
		&container.TabItem{
			Text: "地图画",
			Content: container.NewVBox(
				pathOption,
				widget.NewForm(mapXFormItem,
					mapZFormItem,
					mapYFormItem,
				),
				container.NewGridWithColumns(2, widget.NewLabel("图片绘制起点(建议为64的奇数倍)"), g.startPos.UpdateBtn),
				g.startPos.PosContent,
				g.makeConfirmButton("制图", func() {
					path, fp, err := pathGet()
					if err != nil {
						return
					}
					mapX, err := mapXGet()
					if err != nil {
						return
					}
					mapZ, err := mapZGet()
					if err != nil {
						return
					}
					mapY, err := mapYGet()
					if err != nil {
						return
					}
					if mapY < 20 {
						mapY = 0
					}
					err = g.setStartPos()
					if err != nil {
						return
					}
					g.addMonkeyPathReader(path, fp)
					g.sendCmdAndClose(fmt.Sprintf("mapart -p %v -mapX %v -mapZ %v -mapY %v", path, mapX, mapZ, mapY))
				}),
			),
		},
	)
	return c
}

func (g *GUI) makeNbtContent() fyne.CanvasObject {
	if !g.BotSession.Config.NBTConstructorEnabled {
		return widget.NewLabel("应FB要求，你不能在该服务器使用此功能")
	}
	pathOption, pathGet := g.makeReadPathOption("选择nbt文件", "json/txt", []string{".json", ".txt"})
	nbtEntry := widget.NewMultiLineEntry()
	nbtEntry.SetPlaceHolder(`{
		"name": "chest",
		"nbt":{
			"Findable:char":1,
			"LootTable": "loot_tables/chests/end_city_treasure.json"
			"display":{
				"Name": "Lucky",
				"Lore": ["+(DATA)"]
			}
		}
	}
	`)
	return container.NewVBox(
		widget.NewLabel("从文件构造nbt物品"),
		pathOption,
		g.makeConfirmButton("构造", func() {
			path, fp, err := pathGet()
			if err != nil {
				return
			}
			cmd := fmt.Sprintf("construct %v", path)
			g.addMonkeyPathReader(path, fp)
			g.sendCmdAndClose(cmd)
		}),
		widget.NewSeparator(),
		widget.NewLabel("从文字构造nbt物品"),
		nbtEntry,
		g.makeConfirmButton("构造", func() {
			cmd := fmt.Sprintf("simpleconstruct %v", nbtEntry.Text)
			g.sendCmdAndClose(cmd)
		}),
	)

}

func (g *GUI) makeExportContent() fyne.CanvasObject {
	pathOption, pathGet := g.makeWritePathOption("导出到建筑文件", ".bdx", []string{".bdx"})
	return container.NewVBox(
		pathOption,
		container.NewGridWithColumns(2, widget.NewLabel("导出建筑起点位置"), g.startPos.UpdateBtn),
		g.startPos.PosContent,
		container.NewGridWithColumns(2, widget.NewLabel("导出建筑终点位置"), g.endPos.UpdateBtn),
		g.endPos.PosContent,
		g.makeConfirmButton("导出", func() {
			path, fp, err := pathGet()
			if err != nil {
				return
			}
			err = g.setStartPos()
			if err != nil {
				return
			}
			err = g.setEndPos()
			if err != nil {
				return
			}
			cmd := fmt.Sprintf("export -p %v", path)
			g.addMonkeyPathWriter(path, fp)
			g.sendCmdAndClose(cmd)
		}),
	)
}

func (g *GUI) makeMajorContent() fyne.CanvasObject {
	return &widget.Accordion{
		Items: []*widget.AccordionItem{
			// &widget.AccordionItem{
			// 	Title: "Debug",
			// 	Detail: container.NewVBox(
			// 		container.NewHBox(widget.NewLabel("Get    "), g.startPos.UpdateBtn),
			// 		g.startPos.PosContent,
			// 		container.NewHBox(widget.NewLabel("Get End"), g.endPos.UpdateBtn),
			// 		g.endPos.PosContent,
			// 	),
			// },
			&widget.AccordionItem{
				Title:  "几何指令 (在空间中构造简单几何体)",
				Detail: g.makeGeoCmdContent(),
			},
			&widget.AccordionItem{
				Title:  "建筑导入 (支持 schematic/bdx/mcacblock)",
				Detail: g.makeBuildingContent(),
			},
			&widget.AccordionItem{
				Title:  "图片及地图画",
				Detail: g.makePlotContent(),
			},
			&widget.AccordionItem{
				Title:  "导出",
				Detail: g.makeExportContent(),
			},
			&widget.AccordionItem{
				Title:  "构造nbt物品",
				Detail: g.makeNbtContent(),
			},
		},
	}
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
