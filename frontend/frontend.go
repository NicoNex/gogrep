package frontend

import (
	"os"
	"runtime"

	"fyne.io/fyne"
	"fyne.io/fyne/widget"
	"fyne.io/fyne/app"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/layout"
)

type Data struct {
	Pattern string
	Glob string
	Path string
}

type Ui struct {
	Datach chan Data
	Stop chan bool
	window fyne.Window
	app fyne.App
	output *widget.Entry
	rentry *widget.Entry
	gentry *widget.Entry
	pentry *widget.Entry
}

func NewUi() Ui {
	var ui = Ui{
		Datach: make(chan Data, 1),
		Stop: make(chan bool, 1),
	}
	ui.loadUi()

	return ui
}

func (u *Ui) Display(text string) {
	u.output.SetText(text)
}

func (u *Ui) Clear() {
	u.Display("")
}

func (u *Ui) onKeyPress(ev *fyne.KeyEvent) {
	if ev.Name == fyne.KeyReturn || ev.Name == fyne.KeyEnter {
		u.Datach <- Data{
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

	// Search button.
	searchBtn := widget.NewButton("Search", u.onSearchPressed)
	searchBtn.Style = widget.PrimaryButton

	// Stop button.
	stopBtn := widget.NewButton("Stop", u.onStopPressed)
	stopBtn.Style = widget.PrimaryButton

	u.window = u.app.NewWindow("gogrep")
	u.window.Resize(fyne.NewSize(960, 540))

	content := widget.NewVBox(
		u.rentry,
		widget.NewHBox(
			u.gentry,
			u.pentry,
			layout.NewSpacer(),
			searchBtn,
			stopBtn,
		),
	)

	u.window.SetContent(
		fyne.NewContainerWithLayout(
			layout.NewBorderLayout(content, nil, nil, nil),
			content,
			widget.NewScrollContainer(u.output),
		),
	)

	u.window.Canvas().SetOnTypedKey(u.onKeyPress)
}

func (u *Ui) onSearchPressed() {
	u.Datach <- Data{
		Pattern: u.rentry.Text,
		Glob: u.gentry.Text,
		Path: u.pentry.Text,
	}
}

func (u *Ui) onStopPressed() {
	u.Stop <- true
}

func (u *Ui) Run() {
	u.window.ShowAndRun()
}
