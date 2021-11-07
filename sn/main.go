package main

import (
	"log"

	"github.com/MrDjeb/sn/game"
	"github.com/MrDjeb/sn/ui2d"
)

func cherr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	ui := &ui2d.UI2d{}
	game.Run(ui)
}
