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

// piece is a game piece that can be moved across the board until it is
// emprinted which then becomes part of the board.
type piece struct {
	points []*point
}

// canMoveDown returns true if the piece can be moved down based on the current
// configuration of the board.
func (p *piece) canMoveDown(b board) bool {
	for _, point := range p.points {
		if !point.canMoveDown(b) {
			return false
		}
	}

	return true
}

func (p *piece) moveDown() {
	for _, point := range p.points {
		point.x++
	}
}

// isIn return true if the given point is part of the piece.
func (p *piece) isIn(point point) bool {
	for _, pp := range p.points {
		if pp.eq(point) {
			return true
		}
	}
	return false
}

type point struct {
	x, y int
}

// eq returns true if two point has the same coordinates.
func (p point) eq(other point) bool {
	return p.x == other.x && p.y == other.y
}

// canMoveDown return true if the point can be moved down.
func (p point) canMoveDown(b board) bool {
	if p.y+1 < len(b.m) && b.m[p.y+1][p.x] != 1 { // Move down allowed.
		return true
	}

	return false
}

type model struct {
	board *board

	pos point

	// The current pience being moved in the board.
	piece piece
}

func initialModel() model {
	return model{
		board: initBoard(),

		piece: pickPiece(),
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
			if m.pos.x < len(m.board.m)-1 {
				m.pos.x++
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.piece.canMoveDown(*m.board) {
				m.piece.moveDown()
				break
			}

			m.board.emprint(m.piece)
			m.piece = pickPiece()

			// if m.pos.y+1 < len(m.board.m) && m.board.m[m.pos.y+1][m.pos.x] != 1 {
			// 	m.pos.y++
			// 	break
			// }

			// m.board.m[m.pos.y][m.pos.x] = 1
			// m.pos = point{}
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

type board struct {
	m [][]int
}

// emprint writes 1' in the board as the points indicate.
func (b *board) emprint(piece piece) {
	for _, p := range piece.points {
		b.m[p.y][p.x] = 1
	}
}

func (b *board) P(y, x int) string {
	v := b.m[y][y]

	if v == 0 {
		return "0"
	}

	return "1"
}

// initBoard creates an empty board.
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

	for i := range 33 {
		var row string
		for j := range 16 {
			if m.piece.isIn(point{i, j}) {
				row += "1"
			} else {
				row += fmt.Sprintf("%v", m.board.m[i][j])
			}

			// if i == m.pos.y && j == m.pos.x {
			// 	row += "1"
			// } else {
			// 	row += fmt.Sprintf("%v", m.board.m[i][j])
			// }
		}
		board += row + "\n"
	}

	return board
}

// gamePieces defines the possible pieces that the game can show.
var gamePieces = []piece{
	{
		points: []*point{{x: 0, y: 0}},
	},
}

// pickPiece returns a random piece for the game.
func pickPiece() piece {
	return gamePieces[0]
}
