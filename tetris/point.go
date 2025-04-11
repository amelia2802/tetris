package tetris

// Point is 1x1 block where a collection of points is a piece of the game.
type Point struct {
	x, y int
}

func MakePoint(x, y int) Point {
	return Point{x: x, y: y}
}

// eq returns true if two point has the same coordinates.
func (p Point) eq(other Point) bool {
	return p.x == other.x && p.y == other.y
}

// canMoveDown return true if the point can be moved down.
func (p Point) canMoveDown(b Board) bool {
	if p.x+1 < len(b.m) && b.m[p.x+1][p.y] != 1 { // Move down allowed.
		return true
	}

	return false
}

// canMoveRight returns true if the point can be moved right.
func (p *Point) canMoveRight(b Board) bool {
	if p.y+1 < len(b.m[0]) && b.m[p.x][p.y+1] != 1 { // Move right allowed.
		return true
	}

	return false
}

// canMoveLeft returns true if the point can be left.
func (p *Point) canMoveLeft(b Board) bool {
	if p.y > 0 && b.m[p.x][p.y-1] != 1 { // Move left is allowed.
		return true
	}

	return false
}
