package gui

import (
	"image/color"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type myTheme struct {
	defaultTheme                                 fyne.Theme
	regular, bold, italic, boldItalic, monospace fyne.Resource
}

func (t *myTheme) SetDark() {
	t.defaultTheme = theme.DarkTheme()
}
func (t *myTheme) SetLight() {
	t.defaultTheme = theme.LightTheme()
}

func (t *myTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	return t.defaultTheme.Color(name, variant)
}

func (t *myTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return t.defaultTheme.Icon(name)
}

func (t *myTheme) Font(style fyne.TextStyle) fyne.Resource {
	if style.Monospace {
		return t.monospace
	}
	if style.Bold {
		if style.Italic {
			return t.boldItalic
		}
		return t.bold
	}
	if style.Italic {
		return t.italic
	}
	return t.regular
}

func (t *myTheme) Size(name fyne.ThemeSizeName) float32 {
	return t.defaultTheme.Size(name)
}

func (t *myTheme) SetFonts(regularFontPath string, monoFontPath string) {
	t.regular = theme.TextFont()
	t.bold = theme.TextBoldFont()
	t.italic = theme.TextItalicFont()
	t.boldItalic = theme.TextBoldItalicFont()
	t.monospace = theme.TextMonospaceFont()

	if regularFontPath != "" {
		t.regular = loadCustomFont(regularFontPath, "Regular", t.regular)
		t.bold = loadCustomFont(regularFontPath, "Bold", t.bold)
		t.italic = loadCustomFont(regularFontPath, "Italic", t.italic)
		t.boldItalic = loadCustomFont(regularFontPath, "BoldItalic", t.boldItalic)
	}
	if monoFontPath != "" {
		t.monospace = loadCustomFont(monoFontPath, "Regular", t.monospace)
	} else {
		t.monospace = t.regular
	}
}

func loadCustomFont(env, variant string, fallback fyne.Resource) fyne.Resource {
	variantPath := strings.Replace(env, "Regular", variant, -1)

	res, err := fyne.LoadResourceFromPath(variantPath)
	if err != nil {
		fyne.LogError("Error loading specified font", err)
		return fallback
	}

	return res
}
