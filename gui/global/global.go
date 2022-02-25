package global

import (
	my_theme "phoenixbuilder_3rd_gui/gui/theme"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var ThemeToggleBtn *ThemeToggler
var InformBtn *widget.Button
var Banner *fyne.Container

type ThemeToggler struct {
	app      fyne.App
	appTheme *my_theme.MyTheme
	Btn      *widget.Button
}

func (tt *ThemeToggler) DataChanged() {
	newIcon := theme.RadioButtonCheckedIcon()
	iv, _ := tt.appTheme.IsLight.Get()
	if iv {
		newIcon = theme.RadioButtonIcon()
	}
	tt.Btn.Icon = newIcon
}

func MakeThemeToggleBtn(app fyne.App, appTheme *my_theme.MyTheme) *ThemeToggler {
	if ThemeToggleBtn != nil {
		return ThemeToggleBtn
	}
	t := &ThemeToggler{appTheme: appTheme, app: app}
	initIcon := theme.RadioButtonCheckedIcon()
	iv, _ := appTheme.IsLight.Get()
	if iv {
		initIcon = theme.RadioButtonIcon()
	}
	toggleBtn := &widget.Button{
		DisableableWidget: widget.DisableableWidget{},
		Text:              "",
		Icon:              initIcon,
		Importance:        widget.LowImportance,
		Alignment:         0,
		IconPlacement:     0,
		OnTapped: func() {
			iv, _ := t.appTheme.IsLight.Get()
			if iv {
				t.appTheme.SetDark()
				app.Settings().SetTheme(t.appTheme)
			} else {
				t.appTheme.SetLight()
				app.Settings().SetTheme(t.appTheme)
			}
		},
	}
	appTheme.IsLight.AddListener(t)
	t.Btn = toggleBtn
	ThemeToggleBtn = t
	return ThemeToggleBtn
}

func MakeInformPopButton(win fyne.Window) *widget.Button {
	if InformBtn != nil {
		return InformBtn
	}
	InformBtn = &widget.Button{
		DisableableWidget: widget.DisableableWidget{},
		Text:              "",
		Icon:              theme.InfoIcon(),
		Importance:        widget.LowImportance,
		Alignment:         0,
		IconPlacement:     0,
		OnTapped: func() {
			dialog.NewInformation("说明", "本项目是PhoenixBuilder的第三方GUI版本\n项目的核心(FB)为:\nhttps://github.com/LNSSPsd/PhoenixBuilder\n核心功能开发者为: Ruphane, CAIMEO\n界面开发者: CMA2401PT", win).Show()
		},
	}
	return InformBtn
}

func MakeBannner(build string) *fyne.Container {
	if Banner != nil {
		return Banner
	}
	Banner = container.NewBorder(nil, &widget.Separator{},
		widget.NewLabel("FB.3rd.GUI (Alpha) "+build),
		container.NewGridWithColumns(2, InformBtn, ThemeToggleBtn.Btn),
		widget.NewLabel(""),
	)
	return Banner
}
