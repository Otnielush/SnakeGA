package main

import (
	"math/rand"
	"time"
)

//Оценивать результат будем по кол-ву собранных яблок
//Но можно добавить и время жизни

//Возвращает номера лидеров в присланном массиве
type vibor struct {
	nomer int
	value int
}

// func (s *Mode) Selection(k *LocalParam) {
// 	if k.population == 1 {
// 		return
// 	}
// 	// pf("sel: %d\n", s.Snakes[4].ApplesEaten)

// 	temp := make([]Snake, len(s.Snakes))
// 	InsertionSort(&temp)

// 	// Будем ли делать новое поколение
// 	if temp[0].ApplesEaten < k.bestResultApple {
// 		k.newGen = false
// 		pf("Тупые! (%d)-%d", k.bestResultApple, temp[0].ApplesEaten)
// 	} else {
// 		k.newGen = true
// 		copy(s.Snakes, temp)
// 		k.bestResultApple = temp[0].ApplesEaten

// 		for i := 0; i < k.population; i++ {
// 			s.Snakes[i] = Snake{}
// 			s.Snakes[i].tail = make([]Possition, k.lenSnakeStart)
// 		}
// 	}
// }
func (s *Mode) Selection(k *LocalParam) {
	if k.population == 1 {
		return
	}

	InsertionSort(&s.Snakes)

	for i := 0; i < k.population; i++ {
		s.Snakes[i] = Snake{}
		s.Snakes[i].tail = make([]Possition, k.lenSnakeStart)
	}

}

func InsertionSort(mass *[]Snake) {
	for i := 0; i < len(*mass); i++ {
		for k := i; k > 0 && (*mass)[k].ApplesEaten > (*mass)[k-1].ApplesEaten; k-- {
			(*mass)[k].ApplesEaten, (*mass)[k-1].ApplesEaten = (*mass)[k-1].ApplesEaten, (*mass)[k].ApplesEaten
			(*mass)[k].Brain.Weights, (*mass)[k-1].Brain.Weights = (*mass)[k-1].Brain.Weights, (*mass)[k].Brain.Weights
			(*mass)[k].Generation, (*mass)[k-1].Generation = (*mass)[k-1].Generation, (*mass)[k].Generation
			(*mass)[k].Moves, (*mass)[k-1].Moves = (*mass)[k-1].Moves, (*mass)[k].Moves
		}
	}
}

//ХЗ нужно или нет
func (k *LocalParam) Crossover() {
	//Новое поколение
	//Создать новых потомков к прошлым лидерам (num число в поколении)
	for i := k.numLeaders - 1; i < k.population; i++ {

	}
}

func (s *Mode) NewPopulation(p *LocalParam, k *MAP) {
	var numMut int
	if s.Snakes[0].ApplesEaten < p.bestResultApple {
		numMut = 1
		pf("Тупые! (%d)-%d", p.bestResultApple, s.Snakes[0].ApplesEaten)
		// fmt.Println("numMut", numMut)
	} else {
		numMut = p.numMut
		p.bestResultApple = s.Snakes[0].ApplesEaten
		// fmt.Println("numMut", numMut)
	}

	if p.population == 1 {
		s.Snakes[0].Mutation(numMut, p.mutationRate)
		s.Snakes[0].head.X, s.Snakes[0].head.Y = 5+(5*1), 15+15*1
		s.Snakes[0].alive = true
		s.Snakes[0].Generation++
		s.Snakes[0].tail = make([]Possition, p.lenSnakeStart)
		k.kletki[5+(5*1)][15+15*(1)] = 4

		for j := 0; j < len(s.Snakes[0].tail); j++ {
			s.Snakes[0].tail[j].X, s.Snakes[0].tail[j].Y = s.Snakes[0].head.X, s.Snakes[0].head.Y-1-j
			k.kletki[s.Snakes[0].tail[j].X][s.Snakes[0].tail[j].Y-1-j] = 3
		}
	}

	if p.numLeaders == 1 {
		for i := 1; i < p.population; i++ {
			s.Snakes[i].Generation = s.Snakes[0].Generation
			s.Snakes[i].Brain.Weights = s.Snakes[0].Brain.Weights
			s.Snakes[i].Mutation(numMut, p.mutationRate)
			s.Snakes[i].tail = make([]Possition, p.lenSnakeStart)
		}
	} else {
		for i := 3; i < p.population/2; i++ {
			s.Snakes[i].Generation = s.Snakes[0].Generation
			s.Snakes[i].Brain.Weights = s.Snakes[0].Brain.Weights
			s.Snakes[i].Mutation(numMut, p.mutationRate)
			s.Snakes[i].tail = make([]Possition, p.lenSnakeStart)
		}
		for i := p.population / 2; i < p.population*3/4; i++ {
			s.Snakes[i].Generation = s.Snakes[1].Generation
			s.Snakes[i].Brain.Weights = s.Snakes[1].Brain.Weights
			s.Snakes[i].Mutation(numMut, p.mutationRate)
			s.Snakes[i].tail = make([]Possition, p.lenSnakeStart)
		}
		for i := p.population * 3 / 4; i < p.population; i++ {
			s.Snakes[i].Generation = s.Snakes[2].Generation
			s.Snakes[i].Brain.Weights = s.Snakes[2].Brain.Weights
			s.Snakes[i].Mutation(numMut, p.mutationRate)
			s.Snakes[i].tail = make([]Possition, p.lenSnakeStart)
		}
	}

	//++Generation и обнуляем яблоки
	for i := 0; i < p.population; i++ {
		s.Snakes[i].Generation++
		s.Snakes[i].ApplesEaten = 0
	}

	// Раскидываем по карте и оживляем
	s.placeSnake(k, p.population)
}

func (s *Mode) placeSnake(k *MAP, num int) {
	for i := 0; i < num; i++ {
		s.Snakes[i].head.X, s.Snakes[i].head.Y = 5+(5*i), 15+15*(i%3)
		s.Snakes[i].alive = true
		k.kletki[5+(5*i)][15+15*(i%3)] = 4

		// for j := 0; j < k.lenSnakeStart; j++ {
		// 	k.Snakes[i].tail[j].X, k.Snakes[i].tail[j].Y = k.Snakes[i].head.X, k.Snakes[i].head.Y-j
		// 	k.karta.kletki[k.Snakes[i].tail[j].X][k.Snakes[i].tail[j].Y] = 3
		// }
	}
}

func (z *Snake) Mutation(numMut int, MutationRate float64) {
	if numMut == 0 {
		return
	}
	s1 := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s1)
	for j := 0; j < 4; j++ {
		for i := 0; i < numMut; i++ {
			// от -0.1 до 0.1 шаг:0.01
			z.Brain.Weights[rand.Intn(50)][j] += (MutationRate * float64(r.Intn(10)) / 50) - 0.1*MutationRate
		}
	}
}

func NewSnake(lenghtTail int) Snake {
	temp := Snake{}
	temp.tail = make([]Possition, lenghtTail)

	// 0-left, 1-right, 2-up, 3-down
	//яблоки
	for i := 0; i < 25; i += 5 {
		temp.Brain.Weights[i][0] = 1
	}
	for i := 4; i < 25; i += 5 {
		temp.Brain.Weights[i][1] = 1
	}
	for i := 0; i < 5; i++ {
		temp.Brain.Weights[i][2] = 1
	}
	for i := 20; i < 25; i++ {
		temp.Brain.Weights[i][3] = 1
	}
	for i := 6; i < 17; i += 5 {
		temp.Brain.Weights[i][0] = 2
	}
	for i := 8; i < 19; i += 5 {
		temp.Brain.Weights[i][1] = 2
	}
	for i := 6; i < 9; i++ {
		temp.Brain.Weights[i][2] = 2
	}
	for i := 16; i < 19; i++ {
		temp.Brain.Weights[i][3] = 2
	}

	//преграды
	for i := 25; i < 49; i += 5 {
		temp.Brain.Weights[i][0] = -1
	}
	for i := 29; i < 50; i += 5 {
		temp.Brain.Weights[i][1] = -1
	}
	for i := 25; i < 30; i++ {
		temp.Brain.Weights[i][2] = -1
	}
	for i := 45; i < 50; i++ {
		temp.Brain.Weights[i][3] = -1
	}
	for i := 21; i < 42; i += 5 {
		temp.Brain.Weights[i][0] = -2
	}
	for i := 33; i < 44; i += 5 {
		temp.Brain.Weights[i][1] = -2
	}
	for i := 31; i < 34; i++ {
		temp.Brain.Weights[i][2] = -2
	}
	for i := 41; i < 44; i++ {
		temp.Brain.Weights[i][3] = -2
	}

	return temp
}
