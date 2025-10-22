package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"time"
)

const (
	width  = 40
	height = 20
)

type Point struct {
	x, y int
}

type Direction int

const (
	Up Direction = iota
	Down
	Left
	Right
)

type Game struct {
	snake     []Point
	food      Point
	direction Direction
	score     int
	gameOver  bool
}

func clearScreen() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func (g *Game) init() {
	g.snake = []Point{{width / 2, height / 2}}
	g.direction = Right
	g.score = 0
	g.gameOver = false
	g.spawnFood()
}

func (g *Game) spawnFood() {
	for {
		g.food = Point{rand.Intn(width), rand.Intn(height)}
		valid := true
		for _, segment := range g.snake {
			if segment == g.food {
				valid = false
				break
			}
		}
		if valid {
			break
		}
	}
}

func (g *Game) update() {
	if g.gameOver {
		return
	}

	head := g.snake[0]
	var newHead Point

	switch g.direction {
	case Up:
		newHead = Point{head.x, head.y - 1}
	case Down:
		newHead = Point{head.x, head.y + 1}
	case Left:
		newHead = Point{head.x - 1, head.y}
	case Right:
		newHead = Point{head.x + 1, head.y}
	}

	// Check wall collision
	if newHead.x < 0 || newHead.x >= width || newHead.y < 0 || newHead.y >= height {
		g.gameOver = true
		return
	}

	// Check self collision
	for _, segment := range g.snake {
		if segment == newHead {
			g.gameOver = true
			return
		}
	}

	g.snake = append([]Point{newHead}, g.snake...)

	// Check food collision
	if newHead == g.food {
		g.score += 10
		g.spawnFood()
	} else {
		g.snake = g.snake[:len(g.snake)-1]
	}
}

func (g *Game) render() {
	clearScreen()

	// Create grid
	grid := make([][]rune, height)
	for i := range grid {
		grid[i] = make([]rune, width)
		for j := range grid[i] {
			grid[i][j] = ' '
		}
	}

	// Place food
	grid[g.food.y][g.food.x] = '●'

	// Place snake
	for i, segment := range g.snake {
		if i == 0 {
			grid[segment.y][segment.x] = '█'
		} else {
			grid[segment.y][segment.x] = '▓'
		}
	}

	// Draw border and grid
	fmt.Println("┌" + string(make([]rune, width*2)) + "┐")
	for i := 0; i < width*2; i++ {
		fmt.Print("─")
	}
	fmt.Println()

	for _, row := range grid {
		fmt.Print("│")
		for _, cell := range row {
			fmt.Print(string(cell) + " ")
		}
		fmt.Println("│")
	}

	fmt.Print("└")
	for i := 0; i < width*2; i++ {
		fmt.Print("─")
	}
	fmt.Println("┘")

	fmt.Printf("\nScore: %d | Length: %d\n", g.score, len(g.snake))
	fmt.Println("Controls: W=Up, S=Down, A=Left, D=Right, Q=Quit")

	if g.gameOver {
		fmt.Println("\n🎮 GAME OVER! Press Q to quit or R to restart")
	}
}

func getInput(inputChan chan rune) {
	// Disable input buffering
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
	defer exec.Command("stty", "-F", "/dev/tty", "echo").Run()

	var b []byte = make([]byte, 1)
	for {
		os.Stdin.Read(b)
		inputChan <- rune(b[0])
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	game := &Game{}
	game.init()

	inputChan := make(chan rune)
	go getInput(inputChan)

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			game.update()
			game.render()

		case input := <-inputChan:
			switch input {
			case 'w', 'W':
				if game.direction != Down {
					game.direction = Up
				}
			case 's', 'S':
				if game.direction != Up {
					game.direction = Down
				}
			case 'a', 'A':
				if game.direction != Right {
					game.direction = Left
				}
			case 'd', 'D':
				if game.direction != Left {
					game.direction = Right
				}
			case 'r', 'R':
				if game.gameOver {
					game.init()
				}
			case 'q', 'Q':
				clearScreen()
				fmt.Println("Thanks for playing!")
				return
			}
		}
	}
}
