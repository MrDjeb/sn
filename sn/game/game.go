package game

import (
	"fmt"
	"math/rand"
	"time"
)

type GameUI interface {
	Draw(*Field)
	GetInput() *Input
}

type InputType int

const (
	None InputType = iota
	Up
	Down
	Left
	Right
	Quit
)

type Input struct {
	Typ InputType
}

type Segm struct {
	X int
	Y int
}

type Snake struct {
	Body      []Segm
	IsDead    bool
	IndexOfAI int
}

type Field struct {
	Snakes []Snake
	Fruit  Segm
}

func handleInput(field *Field, input *Input) {
	switch input.Typ {
	case Up:
		fmt.Println(input.Typ)
	case Down:
		fmt.Println(input.Typ)
	case Left:
		fmt.Println(input.Typ)
	case Right:
		fmt.Println(input.Typ)

	}
}

func generate(hms int) *Field {
	rand.Seed(time.Now().UTC().UnixNano())
	field := &Field{}
	field.Snakes = []Snake{{Body: []Segm{{0, 0}}, IsDead: false, IndexOfAI: hms}}
	field.Fruit = Segm{rand.Intn(1200), rand.Intn(720)}
	return field
}

func Run(ui GameUI) {
	field := generate(3)
	for {
		ui.Draw(field)
		input := ui.GetInput()

		if input != nil && input.Typ == Quit {
			return
		}

		handleInput(field, input)
	}

}
