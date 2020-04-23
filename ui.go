package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
	"fyne.io/fyne/app"
	"fyne.io/fyne/theme"
)

// TODO: move all these objects in the ui package and in different files.
// TODO: add glob and path option.

type patternEntry struct {
	widget.Entry
	pch chan string
}

type Ui struct {
	// output *widget.Label
	window fyne.Window
	app fyne.App
	output *widget.Entry
}

func newPatternEntry(patternCh chan string) *patternEntry {
	return &patternEntry{pch: patternCh}
}

func NewUi(pch chan string) Ui {
	var ui Ui
	ui.loadUi(pch)

	return ui
}

func (p *patternEntry) KeyDown(ev *fyne.KeyEvent) {
	if ev.Name == fyne.KeyReturn || ev.Name == fyne.KeyEnter {
		p.pch <- p.Entry.Text
	}
}

func (u *Ui) Display(text string) {
	u.output.SetText(text)
}

func (u *Ui) Clear() {
	u.Display("")
}

func (u *Ui) loadUi(pch chan string) {
	u.app = app.New()
	u.app.Settings().SetTheme(theme.LightTheme())

	// u.output = widget.NewLabel("")
	u.output = widget.NewMultiLineEntry()
	u.output.Disable()
	// u.output.Alignment = fyne.TextAlignLeading
	// u.output.TextStyle.Monospace = true

	pe := newPatternEntry(pch)
	pe.ExtendBaseWidget(pe)
	pe.SetPlaceHolder("pattern")

	u.window = u.app.NewWindow("gogrep")
	u.window.Resize(fyne.NewSize(960, 540))
	u.window.SetContent(widget.NewVBox(
		pe,
		u.output,
	))

	// u.window.Canvas().SetOnTypedKey(u.onKeyPress)
}

func (u *Ui) run() {
	u.window.ShowAndRun()
}
