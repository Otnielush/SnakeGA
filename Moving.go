//Схема нейросети:
//змея видит квадрат 5 на 5 = 25 для яблока и 25 для преграды
//(Но можно не делать разные входы, а разными цифрами вводить разные элементы на карте. Нужно тестить)

//50 входных нейронов, 4 средних (налево, направо, вверх, вниз) и на выход 1 вариант
//на вход 0 если ничего, 1 если есть яблоко/преграда

//Думаю на карте обозначать наличие на клетке чего либа разными цифрами:
//0-пусто, 1-яблоко, 2 , Преграды: 3 - змея, 4 - голова, 5 - преграда ,6...

package main

import "math"

//--------------------набросок----------------------

//--------------------набросок--------------------

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
	vision   [50]float64    //0-24 на яблоко, 25-50 на преграду
	turns    [4]float64     //0-left, 1-right, 2-up, 3-down
	desicion int            //left, right, up, down
	Weights  [50][4]float64 `json:"Weights"` //50*4=200
}

//For moving
func (z *Snake) Move(k *LocalParam) {
	// pf("___snakes: %1.2f\n", z.Brain.Weights)
	z.DesicionToMove(k)

	// Хвост на карте обнулили

	switch z.Brain.desicion {
	case 0: //left
		// Пошла в стенку/преграду/змею
		if z.head.X <= 0 || k.karta.kletki[z.head.X-1][z.head.Y] >= 3 {
			z.alive = false
			z.dead()
			return
		}
		k.karta.kletki[z.tail[len(z.tail)-1].X][z.tail[len(z.tail)-1].Y] = 0
		for t := len(z.tail) - 1; t > 0; t-- {
			z.tail[t].X, z.tail[t].Y = z.tail[t-1].X, z.tail[t-1].Y
			k.karta.kletki[z.tail[t].X][z.tail[t].Y] = 3 // потом убрать
		}
		z.tail[0].X, z.tail[0].Y = z.head.X, z.head.Y
		k.karta.kletki[z.tail[0].X][z.tail[0].Y] = 3
		z.head.X--
		z.Brain.desicion = 1
	case 1: //right
		if z.head.X >= mapWidth-1 || k.karta.kletki[z.head.X+1][z.head.Y] >= 3 {
			z.alive = false
			z.dead()
			return
		}
		k.karta.kletki[z.tail[len(z.tail)-1].X][z.tail[len(z.tail)-1].Y] = 0
		for t := len(z.tail) - 1; t > 0; t-- {
			z.tail[t].X, z.tail[t].Y = z.tail[t-1].X, z.tail[t-1].Y
			k.karta.kletki[z.tail[t].X][z.tail[t].Y] = 3
		}
		z.tail[0].X, z.tail[0].Y = z.head.X, z.head.Y
		k.karta.kletki[z.tail[0].X][z.tail[0].Y] = 3
		z.head.X++
		z.Brain.desicion = 0
	case 2: //up
		if z.head.Y <= 0 || k.karta.kletki[z.head.X][z.head.Y-1] >= 3 {
			z.alive = false
			z.dead()
			return
		}
		k.karta.kletki[z.tail[len(z.tail)-1].X][z.tail[len(z.tail)-1].Y] = 0
		for t := len(z.tail) - 1; t > 0; t-- {
			z.tail[t].X, z.tail[t].Y = z.tail[t-1].X, z.tail[t-1].Y
			k.karta.kletki[z.tail[t].X][z.tail[t].Y] = 3
		}
		z.tail[0].X, z.tail[0].Y = z.head.X, z.head.Y
		k.karta.kletki[z.tail[0].X][z.tail[0].Y] = 3
		z.head.Y--
		z.Brain.desicion = 3
	case 3: //down
		if z.head.Y >= mapHeight-1 || k.karta.kletki[z.head.X][z.head.Y+1] >= 3 {
			z.alive = false
			z.dead()
			return
		}
		k.karta.kletki[z.tail[len(z.tail)-1].X][z.tail[len(z.tail)-1].Y] = 0
		for t := len(z.tail) - 1; t > 0; t-- {
			z.tail[t].X, z.tail[t].Y = z.tail[t-1].X, z.tail[t-1].Y
			k.karta.kletki[z.tail[t].X][z.tail[t].Y] = 3
		}
		z.tail[0].X, z.tail[0].Y = z.head.X, z.head.Y
		k.karta.kletki[z.tail[0].X][z.tail[0].Y] = 3
		z.head.Y++
		z.Brain.desicion = 2
	}
	// Яблочки!
	if k.karta.kletki[z.head.X][z.head.Y] == 1 {
		z.ApplesEaten++
		// p("Food!", z.ApplesEaten)
		k.karta.eaten++
		k.karta.apples--
		k.createApple()
		//а можно прописать параметр длинну и не задействованые части хвоста держать вне карты
		z.tail = append(z.tail, Possition{X: z.tail[len(z.tail)-1].X, Y: z.tail[len(z.tail)-1].Y})
	}

	//Нанесение положения головы змеи на карту
	k.karta.kletki[z.head.X][z.head.Y] = 4
	z.Moves++
}

func (z *Snake) DesicionToMove(k *LocalParam) {

	// обнуление Зрения
	for i := 0; i < len(z.Brain.vision); i++ {
		z.Brain.vision[i] = 0
	}
	neiro := int(0) //Счётчик для массива входных нейронов
	//Сканирует карту 5 на 5 от головы
	for y := z.head.Y - 2; y <= z.head.Y+2; y++ {
		if y < 0 || y >= mapHeight {
			for x := z.head.X - 2; x <= z.head.X+2; x++ {
				z.Brain.vision[neiro+25] = 1
				neiro++
			}
			continue
		}
		for x := z.head.X - 2; x <= z.head.X+2; x++ {
			if x < 0 || x >= mapWidth {
				z.Brain.vision[neiro+25] = 1
			} else {
				switch k.karta.kletki[x][y] {
				case 0: // Обнулили всё
				case 1:
					z.Brain.vision[neiro] = 1
				case 2:
					z.Brain.vision[neiro] = 0.5
				case 3, 4, 5:
					z.Brain.vision[neiro+25] = 1
				}
			}
			neiro++
		}
	}

	//Вычисляем 4 варианта хода
	for i := 0; i < 4; i++ {
		var neiron float64 = 0
		for j := 0; j < 50; j++ {
			neiron += z.Brain.vision[j] * z.Brain.Weights[j][i]
		}
		z.Brain.turns[i] = LeakyRELU(neiron)
	}
	//Выбираем куда идти
	//Желаемые 2 направления 0-left, 1-right, 2-up, 3-down
	turn1, turn2 := maxs(&z.Brain.turns)
	//Это чтоб себя не ела. Ходит на желаемый из возможных направлений
	if z.Brain.desicion == turn1 {
		z.Brain.desicion = turn2
	} else {
		z.Brain.desicion = turn1
	}
}
func LeakyRELU(neiron float64) float64 {
	a := 0.01 // Объявить в main
	if neiron > 0 {
		return neiron
	} else {
		return neiron * a
	}
}
func Sigmoid(neiron float64) float64 {
	return 1 / (1 + math.Pow(math.E, -neiron))
}
func maxs(mass *[4]float64) (int, int) {
	max1, max2 := 0, 0
	for i := 1; i < 4; i++ {
		if (*mass)[i] > (*mass)[max1] {
			max2 = max1
			max1 = i
		}
	}
	if max1 == 0 {
		max2 = 1
	}
	for i := 0; i < 4; i++ {
		if i != max1 && (*mass)[i] > (*mass)[max2] {
			max2 = i
		}
	}
	return max1, max2
}

func (z *Snake) dead() {
	z.color = color{r: 0, g: 0, b: 0}
}
