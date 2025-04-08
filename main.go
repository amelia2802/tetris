package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	w = 10
	h = 20

	// blockChar = "█"
	blockChar = "0"
)

var (
	currentPiece = pickPiece()

	pieces = map[int]piece{
		0: {
			id: 0,
			// |
			// |
			// |
			// |
			points: []*point{{x: 0, y: 0}, {x: 1, y: 0}, {x: 2, y: 0}, {x: 3, y: 0}},
		},
		1: {
			id: 1,
			// ----
			points: []*point{{x: 0, y: 0}, {x: 0, y: 1}, {x: 0, y: 2}, {x: 0, y: 3}},
		},
		2: {
			id: 2,
			// --
			// --
			points: []*point{{x: 0, y: 0}, {x: 0, y: 1}, {x: 1, y: 0}, {x: 1, y: 1}},
		},
		3: {
			id: 3,
			//  --
			// --
			points: []*point{{x: 1, y: 0}, {x: 1, y: 1}, {x: 0, y: 1}, {x: 0, y: 2}},
		},
		4: {
			id: 4,
			//  -
			// ---
			points: []*point{{x: 1, y: 0}, {x: 1, y: 1}, {x: 1, y: 2}, {x: 0, y: 1}},
		},
		5: {
			id: 5,
			// -
			// ---
			points: []*point{{x: 1, y: 0}, {x: 1, y: 1}, {x: 1, y: 2}, {x: 0, y: 0}},
		},
		6: {
			id: 6,
			//   -
			// ---
			points: []*point{{x: 1, y: 0}, {x: 1, y: 1}, {x: 1, y: 2}, {x: 0, y: 2}},
		},
		7: {
			id: 7,
			// --
			//  --
			points: []*point{{x: 0, y: 0}, {x: 0, y: 1}, {x: 1, y: 1}, {x: 1, y: 2}},
		},
		8: {
			id: 8,
			// |
			// | |
			//   |
			points: []*point{{x: 0, y: 0}, {x: 1, y: 0}, {x: 1, y: 1}, {x: 2, y: 1}},
		},
		9: {
			id: 9,
			// |
			// | |
			// |
			points: []*point{{x: 0, y: 0}, {x: 1, y: 0}, {x: 1, y: 1}, {x: 2, y: 0}},
		},
		10: {
			id: 10,
			// ---
			//  -
			points: []*point{{x: 0, y: 0}, {x: 0, y: 1}, {x: 0, y: 2}, {x: 1, y: 1}},
		},
		11: {
			id: 11,
			//  |
			// ||
			//  |
			points: []*point{{x: 0, y: 1}, {x: 1, y: 0}, {x: 1, y: 1}, {x: 2, y: 1}},
		},
		12: {
			id: 12,
			// ||
			// |
			// |
			points: []*point{{x: 0, y: 0}, {x: 0, y: 1}, {x: 1, y: 0}, {x: 2, y: 0}},
		},
		13: {
			id: 13,
			// ---
			//   _
			points: []*point{{x: 0, y: 0}, {x: 0, y: 1}, {x: 0, y: 2}, {x: 1, y: 2}},
		},
		14: {
			id: 14,
			//  |
			//  |
			// ||
			points: []*point{{x: 0, y: 1}, {x: 1, y: 1}, {x: 2, y: 0}, {x: 2, y: 1}},
		},
		15: {
			id: 15,
			// |
			// |
			// ||
			points: []*point{{x: 0, y: 0}, {x: 1, y: 0}, {x: 2, y: 0}, {x: 2, y: 1}},
		},
		16: {
			id: 16,
			// ---
			// -
			points: []*point{{x: 0, y: 0}, {x: 0, y: 1}, {x: 0, y: 2}, {x: 1, y: 0}},
		},

		17: {
			id: 17,
			// ||
			//  |
			//  |
			points: []*point{{x: 0, y: 0}, {x: 0, y: 1}, {x: 1, y: 1}, {x: 2, y: 1}},
		},
		18: {
			id: 18,
			//  |
			// ||
			// |
			points: []*point{{x: 0, y: 1}, {x: 1, y: 0}, {x: 1, y: 1}, {x: 2, y: 0}},
		},
	}

	ratations = map[int]int{
		0:  1,
		1:  0,
		3:  8,
		8:  3,
		4:  9,
		9:  10,
		10: 11,
		11: 4,
		5:  12,
		12: 13,
		13: 14,
		14: 5,
		6:  15,
		15: 16,
		16: 17,
		17: 6,
		7:  18,
		18: 7,
	}

	moves = 0
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
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return timeTick{}
	})
}

// piece is a game piece that can be moved across the board until it is
// emprinted which then becomes part of the board.
type piece struct {
	id     int
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

// canMoveLeft returns true if the piece can be moved left based on the
// current configuration of the board.
func (p *piece) canMoveLeft(b board) bool {
	for _, point := range p.points {
		if !point.canMoveLeft(b) {
			return false
		}
	}

	return true
}

// moveDown moves the piece down.
func (p *piece) moveDown() {
	moves++
	for _, point := range p.points {
		point.x++
	}
}

// moveRight moves the piece right.
func (p *piece) moveRight() {
	for _, point := range p.points {
		point.y++
	}
}

// moveLeft moves the piece left.
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
	var str string
	for _, pp := range p.points {
		str += fmt.Sprintf("(%d,%d) ", pp.x, pp.y)
	}
	return str
}

// point is 1x1 block where a collection of points is a piece of the game.
type point struct {
	x, y int
}

// eq returns true if two point has the same coordinates.
func (p point) eq(other point) bool {
	return p.x == other.x && p.y == other.y
}

// canMoveDown return true if the point can be moved down.
func (p point) canMoveDown(b board) bool {
	if p.x+1 < len(b.m) && b.m[p.x+1][p.y] != 1 { // Move down allowed.
		return true
	}

	return false
}

// canMoveRight returns true if the point can be moved right.
func (p *point) canMoveRight(b board) bool {
	if p.y+1 < len(b.m[0]) && b.m[p.x][p.y+1] != 1 { // Move right allowed.
		return true
	}

	return false
}

// canMoveLeft returns true if the point can be left.
func (p *point) canMoveLeft(b board) bool {
	if p.y > 0 && b.m[p.x][p.y-1] != 1 { // Move left is allowed.
		return true
	}

	return false
}

type model struct {
	board *board

	ticker timeTick
}

func initialModel() model {
	m := model{board: initBoard()}

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

	for i := range h {
		var row string
		for j := range w {
			if currentPiece.isIn(point{i, j}) {
				row += blockChar
			} else if m.board.m[i][j] == 1 {
				row += blockChar // fmt.Sprintf("%v", m.board.m[i][j])
			} else {
				row += " "
			}
		}
		board += "|" + row + "|" + "\n"
	}

	bottom := ""
	for range w + 2 {
		bottom += "_"
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

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			currentPiece = currentPiece.rotate()

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
		}

	case timeTick:
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
	if currentPiece.canMoveDown(*m.board) {
		currentPiece.moveDown()
	} else {
		// if the pieces just showed up and it can't move down, then the game is over.
		if moves == 0 {
			return tea.Quit
		}
	}

	if m.board.emprint(*currentPiece) {
		currentPiece = pickPiece()
	}

	return nil
}

// board represent the current state of the game.
//
// The current moving piece is overlay on top of the board until it is emprinted
// on the board.
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
			b.m = append(b.m[:i], b.m[i+1:]...)
			b.m = append([][]int{make([]int, w)}, b.m...)
			i--
		}
	}
}

// initBoard creates an empty board.
func initBoard() *board {
	b := &board{
		m: make([][]int, 0),
	}

	for range h {
		var row []int
		for range w {
			row = append(row, 0)
		}
		b.m = append(b.m, row)
	}

	return b
}

// pickPiece returns a random piece for the game.
func pickPiece() *piece {
	rnd := rand.Intn(len(pieces))
	picked := pieces[rnd]

	p := &piece{id: picked.id}
	for _, pp := range picked.points {
		p.points = append(p.points, &point{x: pp.x, y: pp.y})
	}

	p.move(0, w/2-1) // center the piece in the board.

	moves = 0
	return p
}

func (p *piece) rotate() *piece {
	if r, ok := ratations[p.id]; ok {
		rr := pieces[r]
		rotated := cpPiece(rr)
		rotated.move(p.points[0].x, p.points[0].y)
		return rotated
	}

	return p
}

func cpPiece(p piece) *piece {
	cp := &piece{id: p.id}
	for _, pp := range p.points {
		cp.points = append(cp.points, &point{x: pp.x, y: pp.y})
	}

	return cp
}

func (p *piece) move(x, y int) {
	for _, pp := range p.points {
		pp.x += x
		pp.y += y
	}
}
