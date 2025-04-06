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

var currentPiece = pickPiece()

// piece is a game piece that can be moved across the board until it is
// emprinted which then becomes part of the board.
type piece struct {
	id     string
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

// canMoveRight returns true if the piece can be moved right based on the
// curretn configuration of the board.
func (p *piece) canMoveRight(b board) bool {
	for _, point := range p.points {
		if !point.canMoveRight(b) {
			return false
		}
	}

	return true
}

func (p *piece) canMoveLeft(b board) bool {
	for _, point := range p.points {
		if !point.canMoveLeft(b) {
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

func (p *piece) moveRight() {
	for _, point := range p.points {
		point.y++
	}
}

func (p *piece) moveLeft() {
	for _, point := range p.points {
		point.y--
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

func (p *piece) String() string {
	return fmt.Sprintf("%v", p.points)
}

type point struct {
	x, y int
}

// eq returns true if two point has the same coordinates.
func (p point) eq(other point) bool {
	return p.x == other.x && p.y == other.y
}

func (p point) String() string {
	return fmt.Sprintf("{%v, %v}", p.x, p.y)
}

// canMoveDown return true if the point can be moved down.
func (p point) canMoveDown(b board) bool {
	if p.x+1 < len(b.m) && b.m[p.x+1][p.y] != 1 { // Move down allowed.
		return true
	}

	return false
}

func (p *point) canMoveRight(b board) bool {
	if p.y+1 < len(b.m[0]) && b.m[p.x][p.y+1] != 1 { // Move right allowed.
		return true
	}

	return false
}

func (p *point) canMoveLeft(b board) bool {
	if p.y > 0 && b.m[p.x][p.y-1] != 1 { // Move left is allowed.
		return true
	}

	return false
}

type model struct {
	board *board
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
			if currentPiece.canMoveLeft(*m.board) {
				currentPiece.moveLeft()
			}

		case "right":
			if currentPiece.canMoveRight(*m.board) {
				currentPiece.moveRight()
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if currentPiece.canMoveDown(*m.board) {
				currentPiece.moveDown()
			}

			if m.board.emprint(*currentPiece) {
				currentPiece = pickPiece()
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

// emprint writes 1' in the board as the points indicate.
func (b *board) emprint(piece piece) bool {
	if piece.canMoveDown(*b) {
		return false
	}

	for _, p := range piece.points {
		b.m[p.x][p.y] = 1
	}

	b.removeFillRows()

	return true
}

func (b *board) removeFillRows() {
	for i := 0; i < len(b.m); i++ {
		var sum int
		for _, j := range b.m[i] {
			sum += j
		}
		if sum == len(b.m[i]) { // the entire row is filled.
			// b.m[i] = make([]int, w)
			b.m = append(b.m[:i], b.m[i+1:]...)
			b.m = append(b.m, make([]int, w))
			i--
		}
	}
}

const (
	w = 10
	h = 20
)

// initBoard creates an empty board.
func initBoard() *board {
	b := &board{
		m: make([][]int, 0),
	}

	for i := 0; i < h; i++ {
		var row []int
		for j := 0; j < w; j++ {
			row = append(row, 0)
		}
		b.m = append(b.m, row)
	}

	return b
}

// View generates a string representing the current state of the board.
func (m model) View() string {
	var board string

	for i := range h {
		var row string
		for j := range w {
			if currentPiece.isIn(point{i, j}) {
				row += "1"
			} else {
				row += fmt.Sprintf("%v", m.board.m[i][j])
			}
		}
		board += row + "\n"
	}

	return board
}

// gamePieces defines the possible pieces that the game can show.
var gamePieces = []piece{
	{
		id: "1",
		// ----
		points: []*point{{x: 0, y: 0}, {x: 0, y: 1}, {x: 0, y: 2}, {x: 0, y: 3}},
	},
	{
		id: "2",
		// --
		// --
		points: []*point{{x: 0, y: 0}, {x: 0, y: 1}, {x: 1, y: 0}, {x: 1, y: 1}},
	},
	{
		id: "3",
		//  --
		// --
		points: []*point{{x: 1, y: 0}, {x: 1, y: 1}, {x: 0, y: 1}, {x: 0, y: 2}},
	},
	{
		id: "4",
		//  -
		// ---
		points: []*point{{x: 1, y: 0}, {x: 1, y: 1}, {x: 1, y: 2}, {x: 0, y: 1}},
	},
	{
		id: "5",
		// -
		// ---
		points: []*point{{x: 1, y: 0}, {x: 1, y: 1}, {x: 1, y: 2}, {x: 0, y: 0}},
	},
	{
		id: "6",
		//   -
		// ---
		points: []*point{{x: 1, y: 0}, {x: 1, y: 1}, {x: 1, y: 2}, {x: 0, y: 2}},
	},
	{
		id: "7",
		// --
		//  --
		points: []*point{{x: 0, y: 0}, {x: 0, y: 1}, {x: 1, y: 1}, {x: 1, y: 2}},
	},
}

// pickPiece returns a random piece for the game.
func pickPiece() *piece {
	// rnd := rand.Intn(len(gamePieces))
	picked := gamePieces[1]

	p := &piece{}
	for _, pp := range picked.points {
		p.points = append(p.points, &point{x: pp.x, y: pp.y})
	}

	return p
}
