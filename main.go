package main

import (
	"fmt"

	"github.com/NicoNex/gogrep/ui"
	"github.com/NicoNex/gogrep/backend"
)

func listen(datach chan ui.Data, ui ui.Ui) {
	var buf string
	var grep = backend.NewGrep()

	for d := range datach {
		buf = ""
		ui.Display("Searching...")
		ch, err := grep.Find(d)
		if err != nil {
			ui.Display(err.Error())
			continue
		}

		for s := range ch {
			buf = fmt.Sprintf("%s\n%s", buf, s)
			ui.Display(buf)
		}
	}
}

func main() {
	ui, datach := ui.NewUi()
	go listen(datach, ui)
	ui.Run()
}
