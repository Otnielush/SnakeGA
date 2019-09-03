// Управление: верх,низ - скорость; лево,право - мутация; N - новое поколение; W - показать вес первой змеи

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
	winWidth, winHeight int = 1200, 600
	aRELU                   = 0.01 // LeakyRELU коэфициент
)

var p = fmt.Println
var pf = fmt.Printf

// Разные змеи: учёба, спаринг своя и чужая
type Mode struct {
	Snakes     []Snake `json:"Snakes"` // Учёба
	MySnake    Snake   // спаринг
	EnemySnake Snake   // спаринг

}

type Snake struct {
	head        Possition
	tail        []Possition
	Brain       Brain `json:"Brain"`
	ApplesEaten int   `json:"ApplesEaten"`
	Moves       int
	alive       bool
	color       color
	Generation  int
}
type Possition struct {
	X, Y int
}
type Brain struct {
	vision   [242]float64    //0-120 на яблоко, 121-241 на преграду
	turns    [4]float64      //0-left, 1-right, 2-up, 3-down
	desicion int             //left, right, up, down
	Weights  [242][4]float64 `json:"Weights"` //242*4=968
}
type color struct {
	r, g, b byte
}

// Структура для параметров
type LocalParam struct {
	speed           uint32
	mutationRate    float64
	numMut          int
	population      int
	lenSnakeStart   int
	numLeaders      int
	restart         bool
	bestResultApple int
	bestRes         int
	newGen          bool
	generation      int
}

// Создать карту
type MAP struct {
	kletki    [mapWidth][mapHeight]int //0-пусто, 1-яблоко, 2 , Преграды: 3 - змея, 4 - голова, 5 - преграда ,6...
	apples    int
	maxApples int
	eaten     int
}

// Кнопощке
type interfaceGame struct {
	buttons  []button
	trainig  bool
	fighting bool
	menu     bool
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

func (k *LocalParam) update(keyState []uint8, s *Mode) {
	switch {
	case keyState[sdl.SCANCODE_UP] != 0 && k.speed >= 4:
		k.speed /= 2
		p("Speed", 1000/k.speed, "per second")
		sdl.Delay(200)
	case keyState[sdl.SCANCODE_DOWN] != 0:
		k.speed *= 2
		p("Speed", 1000/k.speed, "per second")
		sdl.Delay(200)
	case keyState[sdl.SCANCODE_N] != 0:
		k.restart = true
		sdl.Delay(1000)
	case keyState[sdl.SCANCODE_LEFT] != 0:
		if k.mutationRate < -10 {
			k.mutationRate = -10
			sdl.Delay(200)
		} else {
			k.mutationRate -= 0.02
			sdl.Delay(200)
		}
		p("MutationRate", k.mutationRate)
	case keyState[sdl.SCANCODE_RIGHT] != 0:
		if k.mutationRate > 10 {
			k.mutationRate = 10
			sdl.Delay(200)
		} else {
			k.mutationRate += 0.02
			sdl.Delay(100)
		}
		p("MutationRate", k.mutationRate)
	case keyState[sdl.SCANCODE_W] != 0:
		pf("\nИсправление эволюции (+/-)\n")
		p("1.", s.Snakes[0].Brain.Weights[60][2])
		for j := 0; j < 4; j++ {
			for i := 0; i < 121; i++ {
				if s.Snakes[0].Brain.Weights[i][j] < 0 {
					s.Snakes[0].Brain.Weights[i][j] *= -1
				}
			}
			for i := 121; i < 242; i++ {
				if s.Snakes[0].Brain.Weights[i][j] > 0 {
					s.Snakes[0].Brain.Weights[i][j] *= -1
				}
			}
		}
		p("2.", s.Snakes[0].Brain.Weights[60][2])
		sdl.Delay(1000)
	case keyState[sdl.SCANCODE_T] != 0:
		// s.Snakes[0].Save("1")
		// s.Snakes[1].Save("2")
		// s.Snakes[2].Save("3")

	}
}

func (s *Mode) startPopulation(k *MAP, p LocalParam) {
	k.kletki[0][0] = 5
	k.kletki[0][1] = 5
	k.kletki[1][1] = 5
	for i := 0; i < p.population; i++ {
		s.Snakes[i] = NewSnake(p.lenSnakeStart)
		s.Snakes[i].head.X, s.Snakes[i].head.Y = 5+(5*i), 5+5*(i%3)
		for j := 0; j < p.lenSnakeStart; j++ {
			s.Snakes[i].tail[j].X, s.Snakes[i].tail[j].Y = s.Snakes[i].head.X, s.Snakes[i].head.Y
		}
		s.Snakes[i].alive = true
		k.kletki[5+(5*i)][5+5*(i%3)] = 4

		s.Snakes[i].Mutation(50, p.mutationRate)
	}
}

//
func main() {
	// Важные переменные
	param := LocalParam{
		restart:       false,
		speed:         128,
		mutationRate:  1,
		population:    10,
		numLeaders:    3,
		lenSnakeStart: 3,
		numMut:        24, //из 242
		newGen:        true,
		generation:    1,
	}
	// Режим и змеи
	learn := Mode{}

	//------ Графика ---------------
	pixels := make([]byte, winWidth*winHeight*4)

	var err error
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
	//------------Графика----------------------

	// Меню и интерфейс
	inter := interfaceGame{trainig: false, fighting: false, menu: true}
	inter.buttons = make([]button, 3)
	// Кнопочке и табло
	inter.buttons[0] = button{name: "chooseMod", pos: Possition{X: 5, Y: winHeight - 30}, width: 70, height: 25, color: color{r: 255, g: 204, b: 0}}
	inter.buttons[1] = button{name: "fight", pos: Possition{X: winWidth/2 - 35, Y: winHeight / 2}, width: 70, height: 25, color: color{r: 255, g: 204, b: 0}}
	inter.buttons[2] = button{name: "trainig", pos: Possition{X: winWidth/2 - 35, Y: winHeight/2 + 25}, width: 70, height: 25, color: color{r: 255, g: 153, b: 51}}

Menu:
	for inter.menu {
		// Закрытие окна и сохранение
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				inter.menu = false
				break
			}
		}

		//Фон
		for y := 0; y < winHeight; y++ {
			for x := 0; x < winWidth; x++ {
				setPixel(x, y, color{180, 180, 190}, pixels)
			}
		}

		// Кнопочка
		//inter.buttons[0].draw(pixels)
		inter.buttons[1].draw(pixels)
		inter.buttons[2].draw(pixels)

		// Управление
		// param.update(keyState, &learn)
		mouseX, mouseY, mouseState := sdl.GetMouseState()
		inter.mouseMenu(mouseX, mouseY, mouseState)

		// Вывод на экран
		tex.Update(nil, pixels, winWidth*4)
		renderer.Copy(tex, nil, nil)
		renderer.Present()

		sdl.Delay(param.speed)

	}

	// TRAINIG+++++++++++++++++++++++++++++++++++++++++++++++++++++++

	// Карта
	karta := MAP{maxApples: 25, apples: 0, eaten: 0}

	// запуск таймера
	timeLimit := time.Second * 30
	timer := time.Now().Add(timeLimit)

	//Загрузка
	op, err := os.Open("snakeLeader1.json")
	// Если нет сохранений
	if err != nil {
		op.Close()
		learn.Snakes = make([]Snake, param.population)
		for i := 0; i < param.population; i++ {
			learn.Snakes[i].tail = make([]Possition, param.lenSnakeStart)
		}

		learn.startPopulation(&karta, param)
		karta.createApple()

	} else {
		// Если есть сохранения то загружает

		learn.Snakes = make([]Snake, param.population)
		// for i := 0; i < 3; i++ {
		// 	learn.Snakes[i].Mutation(50, 0)
		// }

		for i := 0; i < param.population; i++ {
			learn.Snakes[i].tail = make([]Possition, param.lenSnakeStart)
			learn.Snakes[i].alive = false
		}

		learn.Snakes[0].Load(op)

		// pf("snake 1: %1.2f\n", learn.Snakes[0].Brain.Weights)
		param.generation = learn.Snakes[0].Generation
		param.bestResultApple = learn.Snakes[0].ApplesEaten
		p("Generation loaded:", param.generation)
		op.Close()

		op, err = os.Open("snakeLeader2.json")
		learn.Snakes[1].Load(op)
		op.Close()

		op, err = os.Open("snakeLeader3.json")
		learn.Snakes[2].Load(op)
		op.Close()

		learn.placeSnake(&karta, 3)
		karta.createApple()
	}
	// TRAINING===============================================

	// ACTION!!!

	for inter.trainig {
		// Закрытие окна и сохранение
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				inter.trainig = false

				learn.Snakes[0].Save("1")
				// pf("ne op op %2.3f\n", learn.Snakes[0].Brain.Weights)
				learn.Snakes[1].Save("2")
				learn.Snakes[2].Save("3")

				break
			}
		}

		// Движение
		// clear(&this.karta, 1)
		alives := len(learn.Snakes)
		for i := 0; i < len(learn.Snakes); i++ {
			if learn.Snakes[i].alive {
				learn.Snakes[i].Move(&karta)
				// pf("|%d %v; %d: %d %d (%d)", i, Snakes[i].alive, Snakes[i].Brain.desicion, Snakes[i].head.X, Snakes[i].head.Y, len(Snakes[i].tail)) //0-left, 1-right, 2-up, 3-down
			} else {
				alives--
			}
		}
		// pf("\n(--%d--)\n", alives)

		// New Generation
		if alives <= 0 || param.restart || !time.Now().Before(timer) {
			clear(&karta, 0)
			karta.apples = 0
			karta.createApple()
			// p("karta", karta)
			// p("main:", learn.Snakes[4].ApplesEaten)
			learn.Selection(&param)

			learn.NewPopulation(&param, &karta)
			param.generation++
			pf("Gen: %d; Eaten: %d, Time: %v\n", param.generation, karta.eaten, time.Now().Sub(timer))
			karta.eaten = 0
			param.restart = false
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
		ystep := winHeight / mapHeight
		xstep := winWidth / mapWidth
		for y := 0; y < winHeight; y += ystep {
			for x := 0; x < winWidth; x++ {
				setPixel(x, y, color{160, 160, 255}, pixels)
			}
		}
		for y := 0; y < winHeight; y++ {
			for x := 0; x < winWidth; x += xstep {
				setPixel(x, y, color{160, 160, 255}, pixels)
			}
		}

		// Графика
		for y := 0; y < mapHeight; y++ {
			for x := 0; x < mapWidth; x++ {
				if karta.kletki[x][y] != 0 {
					paintSquare(x, y, karta.kletki[x][y], pixels)
				}
			}
		}

		// Кнопочка

		for y := mapHeight - 3; y < mapHeight; y++ {
			for x := 0; x < 7; x++ {
				karta.kletki[x][y] = 5
				paintSquare(x, y, karta.kletki[x][y], pixels)

			}
		}
		inter.buttons[0].draw(pixels)
		// inter.buttons[1].draw(pixels)
		// inter.buttons[2].draw(pixels)

		// Управление
		param.update(keyState, &learn)
		mouseX, mouseY, mouseState := sdl.GetMouseState()
		inter.mouseTraining(mouseX, mouseY, mouseState)

		// Вывод на экран
		tex.Update(nil, pixels, winWidth*4)
		renderer.Copy(tex, nil, nil)
		renderer.Present()

		sdl.Delay(param.speed)

		if inter.menu {
			learn.Snakes[0].Save("1")
			learn.Snakes[1].Save("2")
			learn.Snakes[2].Save("3")
			goto Menu
		}
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

// Load
func (z *Mode) Load2(op *os.File) {
	stat, _ := op.Stat()
	p("__", stat, "__")
	b1 := make([]byte, stat.Size())

	n, _ := op.Read(b1)
	_ = json.Unmarshal(b1, z)

	op.Close()
	p("Прочитаны байты:", n)
	p("Веса:", z.Snakes[0].Brain.Weights)
}

func (z *Snake) Load(op *os.File) {
	stat, _ := op.Stat()
	b1 := make([]byte, stat.Size())

	// p("1.", z.Moves)
	// p("2.", &z.Moves)
	n, _ := op.Read(b1)
	_ = json.Unmarshal(b1, z)

	// p("3.", z.Moves)
	// p("4.", &z.Moves)
	op.Close()
	p("Змея загружена. Прочитаны байты:", n)
	// p("Веса:", z.Brain.Weights)
}

func (s *Snake) Save(a string) {

	db1, _ := json.Marshal(s)
	// if a == "1" {
	// 	p("op op", s.Brain.Weights)
	// }
	ar := fmt.Sprint("snakeLeader", a, ".json")
	f, _ := os.Create(ar)
	kama, _ := f.Write(db1)
	f.Close()
	fmt.Println("Файл сохранён. Байт записано:", kama)
}
