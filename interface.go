package main

type button struct {
	name   string
	pos    Possition
	width  int
	height int
	color  color
}

func (b *button) draw(pixels []byte) {
	// border
	for y := b.pos.Y + 1; y < b.pos.Y+2; y++ {
		for x := b.pos.X; x < b.pos.X+b.width; x++ {
			setPixel(x, y, color{r: 10, g: 10, b: 10}, pixels)
		}
	}
	for y := b.pos.Y + b.height - 2; y < b.pos.Y+b.height-1; y++ {
		for x := b.pos.X; x < b.pos.X+b.width; x++ {
			setPixel(x, y, color{r: 10, g: 10, b: 10}, pixels)
		}
	}
	for y := b.pos.Y + 2; y < b.pos.Y+b.height-2; y++ {
		for x := b.pos.X; x < b.pos.X+1; x++ {
			setPixel(x, y, color{r: 10, g: 10, b: 10}, pixels)
		}
		for x := b.pos.X + b.width - 1; x < b.pos.X+b.width; x++ {
			setPixel(x, y, color{r: 10, g: 10, b: 10}, pixels)
		}
	}
	//button
	for y := b.pos.Y + 2; y < b.pos.Y+b.height-2; y++ {
		for x := b.pos.X + 1; x < b.pos.X+b.width-1; x++ {
			setPixel(x, y, color{r: b.color.r, g: b.color.g, b: b.color.b}, pixels)
		}
	}
}

func (ifg *interfaceGame) mouseMenu(mouseX, mouseY int32, mouseState uint32) {
	// mouseState 1-LBM, 4-RBM, 2-wheel
	switch {
	case int(mouseX) > ifg.buttons[1].pos.X && int(mouseX) < ifg.buttons[1].pos.X+ifg.buttons[1].width && int(mouseY) > ifg.buttons[1].pos.Y && int(mouseY) < ifg.buttons[1].pos.Y+ifg.buttons[1].height:
		// ifg.buttons[1].color.r = 220
		// ifg.buttons[1].color.g = 180
		// ifg.buttons[1].color.b = 0

	case int(mouseX) > ifg.buttons[2].pos.X && int(mouseX) < ifg.buttons[2].pos.X+ifg.buttons[2].width && int(mouseY) > ifg.buttons[2].pos.Y && int(mouseY) < ifg.buttons[2].pos.Y+ifg.buttons[2].height:
		if mouseState == 1 {
			ifg.trainig = true
			ifg.menu = false
		}
		// default:
		// ifg.buttons[1].color = color{255, 204, 0}
	}
}

func (ifg *interfaceGame) mouseTraining(mouseX, mouseY int32, mouseState uint32) {
	// mouseState 1-LBM, 4-RBM, 2-wheel
	switch {
	case int(mouseX) > ifg.buttons[0].pos.X && int(mouseX) < ifg.buttons[0].pos.X+ifg.buttons[0].width && int(mouseY) > ifg.buttons[0].pos.Y && int(mouseY) < ifg.buttons[0].pos.Y+ifg.buttons[0].height:
		if mouseState == 1 {
			ifg.trainig = false
			ifg.menu = true
		}
	}
}
