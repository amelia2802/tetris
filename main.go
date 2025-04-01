package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	fmt.Println("Hello Tetris")

	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

type point struct {
	x, y int
}

type model struct {
	board *board

	pos point
}

func initialModel() model {
	return model{
		board: initBoard(),
	}
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		// The "up" and "k" keys move the cursor up
		case "up", "k":

		case "left":
			if m.pos.x > 0 {
				m.pos.x--
			}

		case "right":
			if m.pos.x < len(m.board.m) {
				m.pos.x++
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.pos.y < len(m.board.m)-1 {
				m.pos.y++
			}
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

type board struct {
	m [][]int
}

func initBoard() *board {
	b := &board{
		m: make([][]int, 0),
	}

	for i := 0; i < 33; i++ {
		var row []int
		for j := 0; j < 16; j++ {
			row = append(row, 0)
		}
		b.m = append(b.m, row)
	}

	return b
}

func (m model) View() string {
	var board string

	for i := 0; i < 33; i++ {
		var row string
		for j := 0; j < 16; j++ {
			if i == m.pos.y && j == m.pos.x {
				row += "1"
			} else {
				row += "0"
			}
		}
		board += row + "\n"
	}

	return board
}
