package main

// import (
// 	"fmt"
// )

func listen(pch chan string, ui Ui) {
	for text := range pch {
		ui.Display(text)
	}
}

func main() {
	var pch = make(chan string)

	ui := NewUi(pch)
	go listen(pch, ui)
	ui.run()
}
