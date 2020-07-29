package frontend

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

type Data struct {
	Pattern string
	Glob    string
	Path    string
}

type Ui struct {
	Datach chan Data
	Stop   chan bool
	window fyne.Window
	app    fyne.App
	output *widget.Entry
	rentry *widget.Entry
	gentry *widget.Entry
	pentry *widget.Entry
}

func NewUi() Ui {
	var ui = Ui{
		Datach: make(chan Data, 1),
		Stop:   make(chan bool, 1),
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

func (u *Ui) showErr(err error) {
	fmt.Println(err)
	dialog.ShowError(err, u.window)
}

func (u *Ui) export() {
	var expdir string

	if runtime.GOOS == "windows" {
		expdir = filepath.Join(os.Getenv("UserProfile"), "gogrep")
	} else {
		expdir = filepath.Join(os.Getenv("HOME"), "gogrep")
	}

	if _, err := os.Stat(expdir); err != nil {
		if os.IsNotExist(err) {
			if err := os.Mkdir(expdir, 0755); err != nil {
				u.showErr(err)
				return
			}
		} else {
			u.showErr(err)
			return
		}
	}

	now := time.Now()
	fname := fmt.Sprintf(
		"export_%d-%02d-%02d_%02d:%02d.txt",
		now.Year(),
		now.Month(),
		now.Day(),
		now.Hour(),
		now.Minute(),
	)
	fpath := filepath.Join(expdir, fname)
	err := ioutil.WriteFile(fpath, []byte(u.output.Text), 0644)
	if err != nil {
		u.showErr(err)
		return
	}
	msg := fmt.Sprintf("Successfully exported in %s", fpath)
	dialog.ShowInformation("Success!", msg, u.window)
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

	// Export button.
	expBtn := widget.NewButton("Export", u.export)

	// Clear button.
	clrBtn := widget.NewButton("Clear", u.Clear)

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
			expBtn,
			clrBtn,
		),
	)

	u.window.SetContent(
		fyne.NewContainerWithLayout(
			layout.NewBorderLayout(content, nil, nil, nil),
			content,
			widget.NewScrollContainer(u.output),
		),
	)
}

func (u *Ui) onSearchPressed() {
	u.Datach <- Data{
		Pattern: u.rentry.Text,
		Glob:    u.gentry.Text,
		Path:    u.pentry.Text,
	}
}

func (u *Ui) onStopPressed() {
	u.Stop <- true
}

func (u *Ui) Run() {
	u.window.ShowAndRun()
}
