package main

import (
	"fmt"
	"os"
	"time"

	"com.github.anicolaspp/tetris/tetris"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	w = 10
	h = 20

	// blockChar = "█"
	blockChar = "0"

	menu = "\np - pause, q - quit, space - drop\n"
)

var (
	currentPiece = tetris.PickPiece(w)

	score = 0
	level = 1
	speed = 1000 * time.Millisecond

	paused = false
)

func main() {
	fmt.Println("Hello Tetris")

	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

// timeTick is a message sent every 1 second.
type timeTick struct{}

// run generates a timeTick command.
func (t timeTick) run() tea.Cmd {
	return tea.Tick(speed, func(t time.Time) tea.Msg {
		return timeTick{}
	})
}

type model struct {
	board *tetris.Board

	ticker timeTick
}

func initialModel() model {
	m := model{board: tetris.NewBoard(w, h)}

	return m
}

func (m model) Init() tea.Cmd {
	// Execute the first time tick command.
	return m.ticker.run()
}

// View generates a string representing the current state of the board with the
// current piece overlay on top.
func (m model) View() string {
	var board string
	board += menu + "\n"

	board += fmt.Sprintf("Score: %d, level: %d\n", score, level)

	if paused {
		board += "GAME PAUSED\n"
	}

	top := ""
	for range w + 2 {
		top += "\\"
	}
	board += top + "\n"

	for i := range h {
		var row string
		for j := range w {
			if currentPiece.IsIn(tetris.MakePoint(i, j)) {
				row += blockChar
			} else if m.board.At(i, j) == 1 {
				row += blockChar // fmt.Sprintf("%v", m.board.m[i][j])
			} else {
				row += " "
			}
		}
		board += "|" + row + "|" + "\n"
	}

	bottom := ""
	for range w + 2 {
		bottom += "\\"
	}

	board += bottom
	return board
}

// Update updates the model as a response to a IO change.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		case " ":
			// Drop the piece.
			if paused {
				return m, nil
			}

			for currentPiece.CanMoveDown(*m.board) {
				currentPiece.MoveDown()
			}

		case "p":
			// Pause the game.
			paused = !paused
			if !paused {
				return m, m.ticker.run()
			}

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if paused {
				return m, nil
			}

			currentPiece = currentPiece.Rotate(w)

		case "left":
			if paused {
				return m, nil
			}

			if currentPiece.CanMoveLeft(*m.board) {
				currentPiece.MoveLeft()
			}

		case "right":
			if paused {
				return m, nil
			}

			if currentPiece.CanMoveRight(*m.board) {
				currentPiece.MoveRight()
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if paused {
				return m, nil
			}

			if currentPiece.CanMoveDown(*m.board) {
				currentPiece.MoveDown()
			}
		}

	case timeTick:
		if paused {
			return m, nil
		}

		if m.moveDown() != nil {
			return m, tea.Quit
		}

		return m, m.ticker.run()
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m model) moveDown() tea.Cmd {
	if currentPiece.CanMoveDown(*m.board) {
		currentPiece.MoveDown()
	} else {
		// if the pieces just showed up and it can't move down, then the game is over.
		if currentPiece.Moves() == 0 {
			return tea.Quit
		}
	}

	if cnt, ok := m.board.Emprint(*currentPiece); ok {
		currentPiece = tetris.PickPiece(w)

		score += 10 * cnt
		tmp := level
		for l, v := range levels {
			if score > v {
				tmp = l
			}
		}
		level = tmp
	}

	speed = tetris.Speed(level)

	return nil
}

var levels = map[int]int{
	1: 0,
	2: 100,
	3: 200,
	4: 300,
	5: 400,
	6: 500,
	7: 600,
	8: 700,
}
