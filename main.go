package main

import (
	"fmt"

	"github.com/NicoNex/gogrep/frontend"
	"github.com/NicoNex/gogrep/backend"
)

var grep backend.Grep
var ui frontend.Ui
var sem chan int

func search(data frontend.Data) {
	var buf string

	ui.Display("Searching...")
	err := grep.Find(data)
	if err != nil {
		ui.Display(err.Error())
		return
	}

	for s := range grep.Outch {
		buf = fmt.Sprintf("%s\n%s", buf, s)
		sem <- 1
		ui.Display(buf)
		<-sem
	}
}

func listen() {
	for {
		select {
		case data := <- ui.Datach:
			go search(data)

		case <- ui.Stop:
			fmt.Println("stop")
			if _, ok := <-grep.Stop; ok {
				grep.Stop <- 1
			}
		}
	}
}

func main() {
	sem = make(chan int, 64)
	ui = frontend.NewUi()
	grep = backend.NewGrep()
	go listen()
	ui.Run()
}
