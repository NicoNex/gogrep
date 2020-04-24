package main

import (
	"fmt"

	"github.com/NicoNex/gogrep/frontend"
	"github.com/NicoNex/gogrep/backend"
)

var grep backend.Grep
var ui frontend.Ui

func search(data frontend.Data) {
	var buf string

	ui.Display("Searching...")
	ch, err := grep.Find(data)
	if err != nil {
		ui.Display(err.Error())
		return
	}

	for s := range ch {
		buf = fmt.Sprintf("%s\n%s", buf, s)
		ui.Display(buf)
	}
}

func listen() {
	for {
		select {
		case data := <- ui.Datach:
			go search(data)

		case <- ui.Stop:
			grep.Stop <- true
		}
	}
}

func main() {
	ui = frontend.NewUi()
	grep = backend.NewGrep()
	go listen()
	ui.Run()
}
