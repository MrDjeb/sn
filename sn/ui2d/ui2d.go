package ui2d

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"log"
	"os"

	"github.com/MrDjeb/sn/game"
	"github.com/veandco/go-sdl2/sdl"
)

func cherr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

type UI2d struct {
}

const SizeSegm, WinW, WinH = 40, 1280, 720

var renderer *sdl.Renderer
var textures map[string]*sdl.Texture
var keyboardState []uint8
var prevKeyboardState []uint8

func (ui *UI2d) GetInput() *game.Input {
	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				return &game.Input{Typ: game.Quit}
			}

			var input game.Input
			if keyboardState[sdl.SCANCODE_UP] == 0 && prevKeyboardState[sdl.SCANCODE_UP] != 0 {
				input.Typ = game.Up
			}
			if keyboardState[sdl.SCANCODE_DOWN] == 0 && prevKeyboardState[sdl.SCANCODE_DOWN] != 0 {
				input.Typ = game.Down
			}
			if keyboardState[sdl.SCANCODE_LEFT] == 0 && prevKeyboardState[sdl.SCANCODE_LEFT] != 0 {
				input.Typ = game.Left
			}
			if keyboardState[sdl.SCANCODE_RIGHT] == 0 && prevKeyboardState[sdl.SCANCODE_RIGHT] != 0 {
				input.Typ = game.Right
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

func (ui *UI2d) Draw(field *game.Field) {
	drawBackground()
	_, _, w, h, _ := textures["apple.png"].Query()
	renderer.Copy(textures["apple.png"], nil, &sdl.Rect{0, 0, w, h})
	renderer.Copy(textures["apple.png"], nil, &sdl.Rect{40, 0, SizeSegm, SizeSegm})
	renderer.Present()
}

func drawBackground() {
	_, _, w, h, err := textures["grass5.jpg"].Query()
	w /= 4
	h /= 4

	cherr(err)
	for y := int32(0); y < WinH; y += h {
		for x := int32(0); x < WinW; x += w {
			renderer.Copy(textures["grass5.jpg"], nil, &sdl.Rect{x, y, w, h})
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

func assetsToTextures(dirPath string) {
	textures = map[string]*sdl.Texture{}

	files, err := ioutil.ReadDir(dirPath)
	cherr(err)

	for _, file := range files {
		infile, err := os.Open(dirPath + file.Name())
		cherr(err)
		defer infile.Close()

		img, _, err := image.Decode(infile)
		cherr(err)

		//fmt.Println(file.Name()[:len(file.Name())-(len(format)+1)], format)
		textures[file.Name()] = imgToTexture(img)
	}
}

func init() {
	sdl.LogSetAllPriority(sdl.LOG_PRIORITY_VERBOSE)
	cherr(sdl.Init(sdl.INIT_EVERYTHING))
	window, err := sdl.CreateWindow("Snake", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, WinW, WinH, sdl.WINDOW_SHOWN)
	cherr(err)
	renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	cherr(err)
	sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "1")

	assetsToTextures("./ui2d/assets/")

	keyboardState = sdl.GetKeyboardState()
	prevKeyboardState = make([]uint8, len(keyboardState))
	for i, v := range keyboardState {
		prevKeyboardState[i] = v
	}
}

/*func loadSprites(renderer *sdl.Renderer) []sprite {
	spriteStr := []string{""}
	sprites := make([]sprite, len(spriteStr))

	for i, str := range spriteStr {
		infile, err := os.Open(str)
		cherr(err)
		defer infile.Close()

		img, err := png.Decode(infile)
		cherr(err)

		w := img.Bounds().Max.X
		h := img.Bounds().Max.Y

		spritePixels := make([]byte, w*h*4)
		ind := 0
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				r, g, b, a := img.At(x, y).RGBA()
				spritePixels[ind], spritePixels[ind+1], spritePixels[ind+2], spritePixels[ind+3] = byte(r/256), byte(g/256), byte(b/256), byte(a/256)
				ind += 4
			}
		}
		tex, err := renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STATIC, int32(w), int32(h))
		cherr(err)
		tex.Update(nil, spritePixels, w*4)

		cherr(tex.SetBlendMode(sdl.BLENDMODE_BLEND))
		sprites[i] = sprite{tex, pos{float32(winW) / 2, float32(winH) / 2}, float32(1), w, h}
	}
	return sprites
}*/
