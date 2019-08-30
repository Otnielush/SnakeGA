package main

import "math/rand"

func (k *LocalParam) createApple() {
	haser := k.karta.maxApples - k.karta.apples
	if haser > 0 {
		for i := 1; i <= haser; i++ {
		PustMesto:
			x, y := randXY()
			if k.karta.kletki[x][y] == 0 {
				k.karta.kletki[x][y] = 1
				// p("создано яблоко", x, y)
				k.karta.apples++
			} else {
				// p("Занято")
				goto PustMesto

			}
		}
	}
}
func randXY() (int, int) {
	return rand.Intn(mapWidth), rand.Intn(mapHeight)
}
