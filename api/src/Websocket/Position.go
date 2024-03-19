package ws

type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func (p Position) Eq(target Position) bool {
	return p.X == target.X && p.Y == target.Y
}

func (p Position) Some(targets []Position) bool {
	for _, target := range targets {
		if p.Eq(target) {
			return true
		}
	}
	return false
}

func (p Position) Movement(dx, dy int) Position {
	p.X += dx
	p.Y += dy
	return p
}

func (p Position) Out(x, y, w, h int) bool {
	return p.X < x || p.Y < y || p.X > x+w || p.Y > y+h
}
