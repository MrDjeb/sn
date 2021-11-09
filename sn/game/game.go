package game

import (
	"container/list"
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
	Quit
)

type Input struct {
	Typ InputType
}

type Direction struct {
	Str string
	Sgm Segm
	Inp InputType
}

var (
	Up    = Direction{"up", Segm{0, -1}, 256}
	Right = Direction{"right", Segm{1, 0}, 257}
	Down  = Direction{"down", Segm{0, 1}, 258}
	Left  = Direction{"left", Segm{-1, 0}, 259}
)

type Segm struct {
	X int32
	Y int32
}

func (segm *Segm) Res(segmT Segm) Segm {
	return Segm{segm.X - segmT.X, segm.Y - segmT.Y}
}

func (segm *Segm) Sum(segmT Segm) Segm {
	return Segm{segm.X + segmT.X, segm.Y + segmT.Y}
}

func (segm *Segm) GetLocation(segmT Segm) Direction {
	switch segmT.Res(*segm) {
	case Up.Sgm:
		return Up
	case Right.Sgm:
		return Right
	case Down.Sgm:
		return Down
	case Left.Sgm:
		return Left
	default:
		return Direction{"", Segm{0, 0}, None}
	}
}

type Snake struct {
	Body      list.List
	IsDead    bool
	IndexOfAI int
}

func (snake *Snake) Move(input *Input) {
	segmHead := snake.Body.Front().Value.(Segm)
	switch input.Typ {
	case Up.Inp:
		snake.Body.PushFront(Up.Sgm.Sum(segmHead))
	case Right.Inp:
		snake.Body.PushFront(Right.Sgm.Sum(segmHead))
	case Down.Inp:
		snake.Body.PushFront(Down.Sgm.Sum(segmHead))
	case Left.Inp:
		snake.Body.PushFront(Left.Sgm.Sum(segmHead))
	}
	snake.Body.Remove(snake.Body.Back())
}

type Field struct {
	w, h   int
	Snakes []Snake
	Fruit  Segm
}

func handleInput(field *Field, input *Input) {
	field.Snakes[0].Move(input)
}

func generate(hms int) *Field {
	rand.Seed(time.Now().UTC().UnixNano())
	field := &Field{}
	field.w, field.h = 30, 18
	field.Snakes = make([]Snake, 1)
	field.Snakes[0].Body.PushBack(Segm{2, 2})
	field.Snakes[0].Body.PushBack(Segm{2, 3})
	field.Snakes[0].Body.PushBack(Segm{3, 3})
	field.Snakes[0].Body.PushBack(Segm{3, 4})
	field.Snakes[0].IsDead = false
	field.Snakes[0].IndexOfAI = hms
	//field.Snakes = []Snake{{Body: []Segm{{0, 0}}, IsDead: false, IndexOfAI: hms}}
	field.Fruit = Segm{int32(rand.Intn(field.w)), int32(rand.Intn(field.h))}
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
