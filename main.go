package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	mapWidth                = 100
	mapHeight               = 50
	winWidth, winHeight int = 1000, 500
)

var p = fmt.Println
var pf = fmt.Printf

// Создать карту
type LocalParam struct {
	karta           MAP
	generation      int
	speed           uint32
	mutationRate    float64
	numMut          int
	population      int
	snakes          []Snake `json:"snakes"`
	lenSnakeStart   int
	numLeaders      int
	restart         bool
	bestResultApple int
}
type MAP struct {
	width, height int
	apples        int
	maxApples     int
	eaten         int
	kletki        [mapWidth][mapHeight]int //0-пусто, 1-яблоко, 2 , Преграды: 3 - змея, 4 - голова, 5 - преграда ,6...
}
type color struct {
	r, g, b byte
}

func setPixel(x, y int, c color, pixels []byte) {
	index := (y*winWidth + x) * 4
	if index < len(pixels)-4 && index >= 0 {
		pixels[index] = c.r
		pixels[index+1] = c.g
		pixels[index+2] = c.b
		// pixels[index+3] = 0
	}
}
func paintSquare(x, y, tip int, pixels []byte) {
	startX := x * (winWidth / mapWidth)
	startY := y * (winHeight / mapHeight)
	switch tip {
	case 1:
		for y := 1; y < (winHeight / mapHeight); y++ {
			for x := 1; x < (winWidth / mapWidth); x++ {
				setPixel(x+startX, y+startY, color{230, 70, 90}, pixels)
			}
		}
	case 2:
		for y := 1; y < (winHeight / mapHeight); y++ {
			for x := 1; x < (winWidth / mapWidth); x++ {
				setPixel(x+startX, y+startY, color{50, 50, 50}, pixels)
			}
		}
	case 3:
		for y := 1; y < (winHeight / mapHeight); y++ {
			for x := 1; x < (winWidth / mapWidth); x++ {
				setPixel(x+startX, y+startY, color{77, 77, 255}, pixels)
			}
		}
	case 4:
		for y := 1; y < (winHeight / mapHeight); y++ {
			for x := 1; x < (winWidth / mapWidth); x++ {
				setPixel(x+startX, y+startY, color{60, 50, 245}, pixels)
			}
		}
	case 5:
		for y := 1; y < (winHeight / mapHeight); y++ {
			for x := 1; x < (winWidth / mapWidth); x++ {
				setPixel(x+startX, y+startY, color{0, 0, 0}, pixels)
			}
		}
	case 0:
	default:
		for y := 1; y < (winHeight / mapHeight); y++ {
			for x := 1; x < (winWidth / mapWidth); x++ {
				setPixel(x+startX, y+startY, color{0, 170, 0}, pixels)
			}
		}
	}
}

func (k *LocalParam) update(keyState []uint8) {
	if keyState[sdl.SCANCODE_UP] != 0 && k.speed >= 4 {
		k.speed /= 2
		p("Speed", 1000/k.speed, "per second")
		sdl.Delay(200)
	} else if keyState[sdl.SCANCODE_DOWN] != 0 {
		k.speed *= 2
		p("Speed", 1000/k.speed, "per second")
		sdl.Delay(200)
	} else if keyState[sdl.SCANCODE_N] != 0 {
		k.restart = true
		sdl.Delay(1000)
	} else if keyState[sdl.SCANCODE_LEFT] != 0 {
		if k.mutationRate < -10 {
			k.mutationRate = -10
			sdl.Delay(200)
		} else {
			k.mutationRate -= 0.02
			sdl.Delay(200)
		}
		p("MutationRate", k.mutationRate)
	} else if keyState[sdl.SCANCODE_RIGHT] != 0 {
		if k.mutationRate > 10 {
			k.mutationRate = 10
			sdl.Delay(200)
		} else {
			k.mutationRate += 0.02
			sdl.Delay(100)
		}
		p("MutationRate", k.mutationRate)
	} else if keyState[sdl.SCANCODE_W] != 0 {
		pf("\nsnake 1: %1.2f\n\n", k.snakes[0].Brain.Weights)
		sdl.Delay(1000)
	}
}

func (k *LocalParam) startPopulation(numMut int) {
	k.karta.kletki[0][0] = 5
	k.karta.kletki[0][1] = 5
	k.karta.kletki[1][1] = 5
	for i := 0; i < k.population; i++ {
		k.snakes[i].head.X, k.snakes[i].head.Y = 5+(5*i), 15+15*(i%3)
		k.snakes[i].alive = true
		k.karta.kletki[5+(5*i)][15+15*(i%3)] = 4
		// Хвост появляется под змеёй
		// for j := 0; j < k.lenSnakeStart; j++ {
		// 	k.snakes[i].tail[j].X, k.snakes[i].tail[j].Y = k.snakes[i].head.X, k.snakes[i].head.Y-1-j
		// 	k.karta.kletki[k.snakes[i].tail[j].X][k.snakes[i].tail[j].Y-1-j] = 3
		// }
		k.snakes[i].Mutation(numMut, k.mutationRate)
	}
}

func main() {
	// Важные переменные
	this := LocalParam{
		karta:         MAP{width: mapWidth, height: mapHeight, maxApples: 50, apples: 0, eaten: 0},
		restart:       false,
		speed:         128,
		mutationRate:  1,
		population:    12,
		numLeaders:    3,
		lenSnakeStart: 3,
		numMut:        3,
		generation:    1,
	}

	// запуск таймера
	timeLimit := time.Second * 15
	timer := time.Now().Add(timeLimit)

	// Если есть сохранения то загружает
	op, err := os.Open("snakeLeaders.json")
	if err != nil {
		op.Close()
		this.snakes = make([]Snake, this.population)
		for i := 0; i < this.population; i++ {
			this.snakes[i].tail = make([]Possition, this.lenSnakeStart)
		}
		this.startPopulation(50)
		this.createApple()

	} else {
		stat, _ := op.Stat()
		b1 := make([]byte, stat.Size())

		n, _ := op.Read(b1)
		_ = json.Unmarshal(b1, &this.snakes)

		op.Close()
		p("Прочитаны байты:", n)

		pf("snake 1: %1.2f\n", this.snakes[0].Brain.Weights)
		this.generation = this.snakes[0].Generation
		this.bestResultApple = this.snakes[0].ApplesEaten
		p("Generation loaded:", this.generation)
		this.snakes = make([]Snake, this.population)
		for i := 0; i < this.population; i++ {
			this.snakes[i].tail = make([]Possition, this.lenSnakeStart)
			this.snakes[i].alive = false
		}

		this.placeSnake(3)
		this.createApple()
	}
	// Графика
	pixels := make([]byte, winWidth*winHeight*4)

	err = sdl.Init(sdl.INIT_EVERYTHING)
	check(err)
	defer sdl.Quit()

	window, err := sdl.CreateWindow("Snake - Generation algorithm", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, int32(winWidth), int32(winHeight), sdl.WINDOW_SHOWN)
	check(err)
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	check(err)
	defer renderer.Destroy()

	tex, err := renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, int32(winWidth), int32(winHeight))
	check(err)
	defer tex.Destroy()

	// Считывание клавы
	keyState := sdl.GetKeyboardState()

	// ACTION!!!
	running := true
	for running {
		// Закрытие окна и сохранение
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				running = false

				db1, _ := json.Marshal(this.snakes[:3])
				f, _ := os.Create("snakeLeaders.json")
				kama, _ := f.Write(db1)
				f.Close()
				fmt.Println("Файл сохранён. Байт записано:", kama)

				break
			}
		}

		// Движение
		// clear(&this.karta, 1)
		alives := len(this.snakes)
		for i := 0; i < len(this.snakes); i++ {
			if this.snakes[i].alive {
				this.snakes[i].Move(&this)
				// pf("|%d %v; %d: %d %d (%d)", i, snakes[i].alive, snakes[i].Brain.desicion, snakes[i].head.X, snakes[i].head.Y, len(snakes[i].tail)) //0-left, 1-right, 2-up, 3-down
			} else {
				alives--
			}
		}
		// pf("\n(--%d--)\n", alives)

		// New Generation
		if alives <= 0 || this.restart || !time.Now().Before(timer) {
			clear(&this.karta, 0)
			this.karta.apples = 0
			this.createApple()
			// p("karta", karta)
			// p("main:", this.snakes[4].ApplesEaten)
			this.Selection(&this.snakes)

			this.NewPopulation(&this.snakes)
			this.generation++
			pf(" Gen: %d; Eaten: %d, Time: %v\n", this.generation, this.karta.eaten, time.Now().Sub(timer))
			this.karta.eaten = 0
			this.restart = false
			timer = time.Now().Add(timeLimit)
			sdl.Delay(100)
		}

		//Фон
		for y := 0; y < winHeight; y++ {
			for x := 0; x < winWidth; x++ {
				setPixel(x, y, color{180, 180, 190}, pixels)
			}
		}

		// Решетка
		for y := 0; y < winHeight; y += 10 {
			for x := 0; x < winWidth; x++ {
				setPixel(x, y, color{160, 160, 255}, pixels)
			}
		}
		for y := 0; y < winHeight; y++ {
			for x := 0; x < winWidth; x += 10 {
				setPixel(x, y, color{160, 160, 255}, pixels)
			}
		}

		// Графика
		for y := 0; y < mapHeight; y++ {
			for x := 0; x < mapWidth; x++ {
				if this.karta.kletki[x][y] != 0 {
					paintSquare(x, y, this.karta.kletki[x][y], pixels)
				}
			}
		}

		// Управление
		this.update(keyState)

		// Вывод на экран
		tex.Update(nil, pixels, winWidth*4)
		renderer.Copy(tex, nil, nil)
		renderer.Present()

		sdl.Delay(this.speed)
	}

}

func check(err error) {
	if err != nil {
		p(err)
	}
	return
}

// ПЕРЕДЕЛАТЬ систему карты
func clear(karta *MAP, apple int) {
	for y := 0; y < mapHeight; y++ {
		for x := 0; x < mapWidth; x++ {
			if (*karta).kletki[x][y] > (0 + apple) {
				(*karta).kletki[x][y] = 0
			}
		}
	}
}
