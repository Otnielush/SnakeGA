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

func (k *LocalParam) Selection(snakes *[]Snake) {
	if k.population == 1 {
		return
	}
	// pf("sel: %d\n", (*snakes)[4].ApplesEaten)

	InsertionSort(snakes)

	for i := k.numLeaders - 1; i < k.population; i++ {
		k.snakes[i] = Snake{}
		k.snakes[i].tail = make([]Possition, k.lenSnakeStart)
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

func (k *LocalParam) NewPopulation(snakes *[]Snake) {
	numMut := k.numMut
	if (*snakes)[0].ApplesEaten <= k.bestResultApple {
		numMut = 1
		pf("Тупые! (%d)-%d", k.bestResultApple, (*snakes)[0].ApplesEaten)
	} else {
		k.bestResultApple = (*snakes)[0].ApplesEaten
	}

	if k.population == 1 {
		(*snakes)[0].Mutation(numMut, k.mutationRate)
		(*snakes)[0].head.X, (*snakes)[0].head.Y = 5+(5*1), 15+15*1
		(*snakes)[0].alive = true
		(*snakes)[0].Generation++
		(*snakes)[0].tail = make([]Possition, k.lenSnakeStart)
		k.karta.kletki[5+(5*1)][15+15*(1)] = 4

		for j := 0; j < len((*snakes)[0].tail); j++ {
			(*snakes)[0].tail[j].X, (*snakes)[0].tail[j].Y = (*snakes)[0].head.X, (*snakes)[0].head.Y-1-j
			k.karta.kletki[(*snakes)[0].tail[j].X][(*snakes)[0].tail[j].Y-1-j] = 3
		}
	}

	if k.numLeaders == 1 {
		for i := 1; i < k.population; i++ {
			(*snakes)[i].Generation = (*snakes)[0].Generation
			(*snakes)[i].Brain.Weights = (*snakes)[0].Brain.Weights
			(*snakes)[i].Mutation(numMut, k.mutationRate)
			(*snakes)[i].tail = make([]Possition, k.lenSnakeStart)
		}
	} else {
		for i := 0; i < k.population/2; i++ {
			(*snakes)[i].Generation = (*snakes)[0].Generation
			(*snakes)[i].Brain.Weights = (*snakes)[0].Brain.Weights
			(*snakes)[i].Mutation(numMut, k.mutationRate)
			(*snakes)[i].tail = make([]Possition, k.lenSnakeStart)
		}
		for i := k.population / 2; i < k.population*3/4; i++ {
			(*snakes)[i].Generation = (*snakes)[1].Generation
			(*snakes)[i].Brain.Weights = (*snakes)[1].Brain.Weights
			(*snakes)[i].Mutation(numMut, k.mutationRate)
			(*snakes)[i].tail = make([]Possition, k.lenSnakeStart)
		}
		for i := k.population * 3 / 4; i < k.population; i++ {
			(*snakes)[i].Generation = (*snakes)[2].Generation
			(*snakes)[i].Brain.Weights = (*snakes)[2].Brain.Weights
			(*snakes)[i].Mutation(numMut, k.mutationRate)
			(*snakes)[i].tail = make([]Possition, k.lenSnakeStart)
		}
	}

	//++Generation и обнуляем яблоки
	for i := 0; i < k.population-1; i++ {
		(*snakes)[i].Generation++
		(*snakes)[i].ApplesEaten = 0
	}

	// Раскидываем по карте и оживляем
	k.placeSnake(k.population)
}

func (k *LocalParam) placeSnake(num int) {
	for i := 0; i < num; i++ {
		k.snakes[i].head.X, k.snakes[i].head.Y = 5+(5*i), 15+15*(i%3)
		k.snakes[i].alive = true
		k.karta.kletki[5+(5*i)][15+15*(i%3)] = 4

		// for j := 0; j < k.lenSnakeStart; j++ {
		// 	k.snakes[i].tail[j].X, k.snakes[i].tail[j].Y = k.snakes[i].head.X, k.snakes[i].head.Y-j
		// 	k.karta.kletki[k.snakes[i].tail[j].X][k.snakes[i].tail[j].Y] = 3
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
