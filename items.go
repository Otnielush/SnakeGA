package main

import "math/rand"

func (k *MAP) createApple() {
	haser := k.maxApples - k.apples
	if haser > 0 {
		for i := 1; i <= haser; i++ {
		PustMesto:
			x, y := randXY()
			if k.kletki[x][y] == 0 {
				k.kletki[x][y] = 1
				// p("создано яблоко", x, y)
				k.apples++
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
