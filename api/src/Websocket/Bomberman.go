package ws

import (
	"math/rand"
	"strconv"
	"time"

	utils "go-api-test.kayn.ooo/src/Utils"
)

type BombermanPlayer struct {
	*Player
	Direction   uint8     `json:"-"` // 0, 1, 2, 3 = top, right, bottom, left
	MoveRoutine chan bool `json:"-"`
	Pos         Position  `json:"pos"`
	CanMove     bool      `json:"-"`
	Speed       uint16    `json:"speed"` // delay in ms until next move
	Power       uint8     `json:"power"`
	Dead        bool      `json:"dead"`
	MaxBomb     uint8     `json:"maxbombs"`
}

type Bomb struct {
	Pos    Position         `json:"pos"`
	Player *BombermanPlayer `json:"-"`
}

type Explosion struct {
	Pos   Position `json:"pos"`
	Dir   string   `json:"dir"` // h | v
	From  uint8    `json:"from"`
	Size  uint8    `json:"size"`
	Start bool     `json:"start"`
	End   bool     `json:"end"`
}

func (e *Explosion) GetPoses() []Position {
	poses := make([]Position, 0)
	if !e.Start || e.End {
		return poses
	}

	dx := uint8(0)
	dy := uint8(0)

	if e.Dir == "h" {
		dx = 1
	} else {
		dy = 1
	}

	for i := range e.Size {
		poses = append(poses, e.Pos.Movement(int(i*dx), int(i*dy)))
	}
	return poses
}

type Bomberman struct {
	Board      [][]int8           `json:"board"` // 0 = empty, 1 = stone, 2 = box, 3 = speed+, 4 = bomb+, 5 = power+
	Players    []*BombermanPlayer `json:"players"`
	Bombs      []*Bomb            `json:"bombs"`
	Explosions []*Explosion       `json:"explosions"`
	Add        chan any           `json:"-"`
	Width      uint8              `json:"-"`
	Height     uint8              `json:"-"`
}

func (b *Bomberman) Init() {
	board := make([][]int8, b.Width)
	for i := range board {
		board[i] = make([]int8, b.Height)
		for j := range board[i] {
			board[i][j] = 0

			if (i > 1 && i < int(b.Width)-2) || (j > 1 && j < int(b.Height)-2) {
				if rand.Float64() < 0.90 {
					board[i][j] = 2
				}
			}

			if i*j%2 != 0 {
				board[i][j] = 1
			}
		}
	}
	b.Board = board

	b.Players = make([]*BombermanPlayer, 0)
	b.Bombs = make([]*Bomb, 0)
	b.Explosions = make([]*Explosion, 0)
}

func (b *Bomberman) GetBombsPositions() []Position {
	var positions []Position
	for _, bomb := range b.Bombs {
		positions = append(positions, bomb.Pos)
	}
	return positions
}

func (b *Bomberman) CountPlayerBombs(player *BombermanPlayer) int {
	count := 0
	for _, bomb := range b.Bombs {
		if bomb.Player == player {
			count++
		}
	}
	return count
}

func (b *Bomberman) GetPlayerIndex(token string) (*BombermanPlayer, int) {
	for i, player := range b.Players {
		if player.Token == token {
			return player, i
		}
	}
	return nil, -1
}

func (b *Bomberman) MovePlayer(player *BombermanPlayer) bool {
	dx, dy := 0, 0

	switch player.Direction {
	case 0:
		dy = -1
	case 1:
		dx = 1
	case 2:
		dy = 1
	case 3:
		dx = -1
	default:
		return false
	}

	newPos := player.Pos.Movement(dx, dy)
	if newPos.Out(0, 0, len(b.Board)-1, len(b.Board[0])-1) {
		return false
	}

	brick := b.Board[newPos.X][newPos.Y]

	if brick == 1 || brick == 2 {
		return false
	}

	if newPos.Some(b.GetBombsPositions()) {
		return false
	}

	player.Pos = newPos
	if b.CheckPlayerDeath(player) {
		return false
	}

	if brick == 3 {
		player.Speed -= 25
	}

	if brick == 4 {
		player.MaxBomb++
	}

	if brick == 5 {
		player.Power += 2
	}

	_, i := b.GetPlayerIndex(player.Token)
	if brick == 3 || brick == 4 || brick == 5 {
		b.Board[newPos.X][newPos.Y] = 0
		b.Add <- Update{"update", "data.board", b.Board}
		b.Add <- Update{"update", "data.players." + strconv.Itoa(i), b.Board}
	}

	return true
}

func (b *Bomberman) CanHaveExplosion(pos Position) bool {
	if pos.Some(b.GetBombsPositions()) {
		return false
	}
	if pos.Out(0, 0, len(b.Board)-1, len(b.Board[0])-1) {
		return false
	}

	cube := b.Board[pos.X][pos.Y]
	if cube == 1 {
		return false
	}

	return true
}

func (b *Bomberman) CreateExplosion(pos Position, dir string, size uint8) *Explosion {
	// top, right, bottom or left blocks outside bomb
	outside := (size - 1) / 2
	// top or left available outside
	tl := 0
	// bottom or right available outside
	br := 0

	// movement where to check for available placement
	dx := 0
	dy := 0

	if dir == "h" {
		dx = 1
	} else {
		dy = 1
	}

	for i := range int(outside) {
		checkPos := pos.Movement(dx*(i+1), dy*(i+1))
		if b.CanHaveExplosion(checkPos) {
			br++
		} else {
			break
		}
		if b.Board[checkPos.X][checkPos.Y] == 2 {
			// stop at first box
			break
		}
	}

	for i := range int(outside) {
		checkPos := pos.Movement(-dx*(i+1), -dy*(i+1))
		if b.CanHaveExplosion(checkPos) {
			tl++
		} else {
			break
		}
	}

	return &Explosion{Pos: pos.Movement(-dx*tl, -dy*tl), Size: uint8(1 + tl + br), Dir: dir, From: uint8(tl), Start: false, End: false}
}

func (b *Bomberman) BreakBoxes(expl *Explosion) bool {
	hasBreak := false
	poses := expl.GetPoses()
	for _, pos := range poses {
		if b.Board[pos.X][pos.Y] == 2 {
			if rand.Float64() < 0.07 {
				b.Board[pos.X][pos.Y] = 3
			} else if rand.Float64() < 0.07 {
				b.Board[pos.X][pos.Y] = 4
			} else if rand.Float64() < 0.07 {
				b.Board[pos.X][pos.Y] = 5
			} else {
				b.Board[pos.X][pos.Y] = 0
			}

			hasBreak = true
		}
	}

	return hasBreak
}

func (b *Bomberman) CheckPlayerDeath(p *BombermanPlayer) bool {
	_, i := b.GetPlayerIndex(p.Token)
	poses := make([]Position, 0)

	for _, explosion := range b.Explosions {
		poses = append(poses, explosion.GetPoses()...)
	}

	if p.Pos.Some(poses) {
		p.Dead = true
		b.Add <- Update{"update", "data.players." + strconv.Itoa(i), p}
		p.MoveRoutine <- false
		close(p.MoveRoutine)
		utils.RemoveFromArray(&b.Players, p)
		utils.SetTimeout(300*time.Millisecond, func() {
			b.Add <- Update{"update", "data.players", b.Players}
		})
		return true
	}
	return false
}

func (b *Bomberman) CheckPlayersDeath() {
	for _, p := range b.Players {
		b.CheckPlayerDeath(p)
	}
}

func (b *Bomberman) PlaceBomb(player *BombermanPlayer) {
	if player.Pos.Some(b.GetBombsPositions()) {
		return
	}

	bomb := &Bomb{
		Pos:    player.Pos,
		Player: player,
	}

	b.Bombs = append(b.Bombs, bomb)
	b.Add <- Update{"update", "data.bombs", b.Bombs}

	explH := b.CreateExplosion(player.Pos, "h", player.Power)
	explV := b.CreateExplosion(player.Pos, "v", player.Power)

	b.Explosions = append(b.Explosions, explH, explV)
	b.Add <- Update{"update", "data.explosions", b.Explosions}

	utils.SetTimeout(2500*time.Millisecond, func() {
		utils.RemoveFromArray(&b.Bombs, bomb)
		b.Add <- Update{"update", "data.bombs", b.Bombs}

		explH.Start = true
		explV.Start = true
		b.Add <- Update{"update", "data.explosions", b.Explosions}

		b.CheckPlayersDeath()

		hasBreak := b.BreakBoxes(explH)
		hasBreak = b.BreakBoxes(explV) || hasBreak
		if hasBreak {
			b.Add <- Update{"update", "data.board", b.Board}
		}

		utils.SetTimeout(1000*time.Millisecond, func() {
			explH.End = true
			explV.End = true
			b.Add <- Update{"update", "data.explosions", b.Explosions}

			utils.SetTimeout(300*time.Millisecond, func() {
				utils.RemoveFromArray(&b.Explosions, explH)
				utils.RemoveFromArray(&b.Explosions, explV)
			})
		})
	})
}
