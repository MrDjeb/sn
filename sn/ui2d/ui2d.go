package ui2d

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"

	"github.com/MrDjeb/sn/game"
	"github.com/veandco/go-sdl2/sdl"
)

func cherr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

type UI2d struct{}

const SizeSegm, WinW, WinH = 40, 1280, 720

var renderer *sdl.Renderer
var window *sdl.Window
var textures map[string]*sdl.Texture
var keyboardState []uint8
var prevKeyboardState []uint8

func (ui *UI2d) GetInput() *game.Input {
	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				return &game.Input{Typ: game.Quit}
				//case *sdl.WindowEvent:
				/////	window.SetSize(e.Data1, e.Data2)
			}

			var input game.Input
			if keyboardState[sdl.SCANCODE_UP] == 0 && prevKeyboardState[sdl.SCANCODE_UP] != 0 {
				input.Typ = game.Up.Inp
			}
			if keyboardState[sdl.SCANCODE_DOWN] == 0 && prevKeyboardState[sdl.SCANCODE_DOWN] != 0 {
				input.Typ = game.Down.Inp
			}
			if keyboardState[sdl.SCANCODE_LEFT] == 0 && prevKeyboardState[sdl.SCANCODE_LEFT] != 0 {
				input.Typ = game.Left.Inp
			}
			if keyboardState[sdl.SCANCODE_RIGHT] == 0 && prevKeyboardState[sdl.SCANCODE_RIGHT] != 0 {
				input.Typ = game.Right.Inp
			}
			for i, v := range keyboardState {
				prevKeyboardState[i] = v
			}
			if input.Typ != game.None {
				return &input
			}

		}
	}
}

func getBodyType(segmPrev, segmNow, segmNext game.Segm) string {
	par1 := segmNow.GetLocation(segmPrev).Inp
	par2 := segmNow.GetLocation(segmNext).Inp
	switch {
	case (par1 == game.Up.Inp && par2 == game.Down.Inp) || (par2 == game.Up.Inp && par1 == game.Down.Inp):
		return game.Up.Str + game.Down.Str
	case (par1 == game.Right.Inp && par2 == game.Left.Inp) || (par2 == game.Right.Inp && par1 == game.Left.Inp):
		return game.Right.Str + game.Left.Str
	case (par1 == game.Up.Inp && par2 == game.Right.Inp) || (par2 == game.Up.Inp && par1 == game.Right.Inp):
		return game.Up.Str + game.Right.Str
	case (par1 == game.Down.Inp && par2 == game.Right.Inp) || (par2 == game.Down.Inp && par1 == game.Right.Inp):
		return game.Down.Str + game.Right.Str
	case (par1 == game.Down.Inp && par2 == game.Left.Inp) || (par2 == game.Down.Inp && par1 == game.Left.Inp):
		return game.Down.Str + game.Left.Str
	case (par1 == game.Up.Inp && par2 == game.Left.Inp) || (par2 == game.Up.Inp && par1 == game.Left.Inp):
		return game.Up.Str + game.Left.Str
	default:
		return ""
	}
}

func getHeadTailType(segmNow, segmT game.Segm) string {
	switch segmNow.GetLocation(segmT).Inp {
	case game.Up.Inp:
		return game.Up.Str
	case game.Right.Inp:
		return game.Right.Str
	case game.Down.Inp:
		return game.Down.Str
	case game.Left.Inp:
		return game.Left.Str
	default:
		return ""
	}
}

func (ui *UI2d) Draw(field *game.Field) {
	drawBackground()

	for _, snake := range field.Snakes {

		renderer.Copy(textures["head_"+getHeadTailType(snake.Body.Front().Value.(game.Segm), snake.Body.Front().Next().Value.(game.Segm))],
			nil, &sdl.Rect{snake.Body.Front().Value.(game.Segm).X * SizeSegm,
				snake.Body.Front().Value.(game.Segm).Y * SizeSegm, SizeSegm, SizeSegm})

		renderer.Copy(textures["tail_"+getHeadTailType(snake.Body.Back().Value.(game.Segm), snake.Body.Back().Prev().Value.(game.Segm))],
			nil, &sdl.Rect{snake.Body.Back().Value.(game.Segm).X * SizeSegm,
				snake.Body.Back().Value.(game.Segm).Y * SizeSegm, SizeSegm, SizeSegm})

		for segm := snake.Body.Front().Next(); segm != snake.Body.Back(); segm = segm.Next() {
			renderer.Copy(textures["body_"+getBodyType(segm.Prev().Value.(game.Segm), segm.Value.(game.Segm), segm.Next().Value.(game.Segm))],
				nil, &sdl.Rect{segm.Value.(game.Segm).X * SizeSegm, segm.Value.(game.Segm).Y * SizeSegm, SizeSegm, SizeSegm})
		}
	}

	renderer.Copy(textures["apple"], nil, &sdl.Rect{field.Fruit.X * SizeSegm, field.Fruit.Y * SizeSegm, SizeSegm, SizeSegm})
	renderer.Present()
}

func drawBackground() {
	rand.Seed(2)
	for y := int32(0); y < WinH; y += SizeSegm {
		for x := int32(0); x < WinW; x += SizeSegm {
			//fmt.Println(string('0'+rune(rand.Intn('8'-'0'+1))))
			renderer.Copy(textures["grass"+string('0'+rune(rand.Intn('7'-'0'+1)))], nil, &sdl.Rect{x, y, SizeSegm, SizeSegm})
		}
	}
}

func imgToTexture(img image.Image) *sdl.Texture {

	w := img.Bounds().Max.X
	h := img.Bounds().Max.Y

	pixels := make([]byte, w*h*4)
	ind := 0
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			pixels[ind], pixels[ind+1], pixels[ind+2], pixels[ind+3] = byte(r/256), byte(g/256), byte(b/256), byte(a/256)
			ind += 4
		}
	}
	tex, err := renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STATIC, int32(w), int32(h))
	cherr(err)
	tex.Update(nil, pixels, w*4)

	cherr(tex.SetBlendMode(sdl.BLENDMODE_BLEND))
	return tex
}

func assetsToTextures(dirPath string) { //Supported formats: png, jpg, gif
	textures = map[string]*sdl.Texture{}
	files, err := ioutil.ReadDir(dirPath)
	cherr(err)
	for _, file := range files {
		infile, err := os.Open(dirPath + file.Name())
		cherr(err)
		defer infile.Close()

		img, _, err := image.Decode(infile)
		if err == image.ErrFormat {
			fmt.Println(err.Error(), "at:", file.Name())
			continue
		}
		cherr(err)

		textures[strings.Join(strings.Split(file.Name(), ".")[:len(strings.Split(file.Name(), "."))-1], "")] = imgToTexture(img)
	}
}

func init() {
	sdl.LogSetAllPriority(sdl.LOG_PRIORITY_VERBOSE)
	cherr(sdl.Init(sdl.INIT_EVERYTHING))
	window, err := sdl.CreateWindow("Snake", sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED, WinW, WinH, sdl.WINDOW_RESIZABLE)
	cherr(err)
	renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	cherr(err)
	sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "1")

	assetsToTextures("./ui2d/assetsV/")

	keyboardState = sdl.GetKeyboardState()
	prevKeyboardState = make([]uint8, len(keyboardState))
	for i, v := range keyboardState {
		prevKeyboardState[i] = v
	}
}
