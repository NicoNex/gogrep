package ui

import (
	"os"
	"runtime"

	"fyne.io/fyne"
	"fyne.io/fyne/widget"
	"fyne.io/fyne/app"
	"fyne.io/fyne/theme"
)

// TODO: move all these objects in the ui package and in different files.
// TODO: add glob and path option.

// type patternEntry struct {
// 	widget.Entry
// 	pch chan string
// }

type Data struct {
	Pattern string
	Glob string
	Path string
}

type Ui struct {
	dch chan Data
	window fyne.Window
	app fyne.App
	output *widget.Entry
	rentry *widget.Entry
	gentry *widget.Entry
	pentry *widget.Entry
}

// func newPatternEntry(patternCh chan string) *patternEntry {
// 	return &patternEntry{pch: patternCh}
// }

func NewUi() (Ui, chan Data) {
	var ui Ui
	ui.dch = make(chan Data)
	ui.loadUi()

	return ui, ui.dch
}

// func (p *patternEntry) KeyDown(ev *fyne.KeyEvent) {
// 	if ev.Name == fyne.KeyReturn || ev.Name == fyne.KeyEnter {
// 		p.pch <- p.Entry.Text
// 	}
// }

func (u *Ui) Display(text string) {
	u.output.SetText(text)
}

func (u *Ui) Clear() {
	u.Display("")
}

func (u *Ui) onKeyPress(ev *fyne.KeyEvent) {
	if ev.Name == fyne.KeyReturn || ev.Name == fyne.KeyEnter {
		u.dch <- Data{
			Pattern: u.rentry.Text,
			Glob: u.gentry.Text,
			Path: u.pentry.Text,
		}
	}
}

func (u *Ui) loadUi() {
	u.app = app.New()
	u.app.Settings().SetTheme(theme.LightTheme())

	u.output = widget.NewMultiLineEntry()
	u.output.Disable()

	// Regex entry.
	u.rentry = widget.NewEntry()
	u.rentry.SetPlaceHolder("pattern")

	// Glob entry.
	u.gentry = widget.NewEntry()
	u.gentry.SetPlaceHolder("glob")

	// Path entry.
	u.pentry = widget.NewEntry()
	if runtime.GOOS == "windows" {
		u.pentry.SetText(os.Getenv("UserProfile"))
	} else {
		u.pentry.SetText(os.Getenv("HOME"))
	}

	u.window = u.app.NewWindow("gogrep")
	u.window.Resize(fyne.NewSize(960, 540))
	u.window.SetContent(
		widget.NewVBox(
			u.rentry,
			widget.NewHBox(
				u.gentry,
				u.pentry,
			),
			u.output,
		),
	)

	u.window.Canvas().SetOnTypedKey(u.onKeyPress)
}

func (u *Ui) Run() {
	u.window.ShowAndRun()
}
