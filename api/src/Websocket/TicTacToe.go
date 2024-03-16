package ws

import (
	"strconv"
)

type TicTacToe struct {
	P1            *Player  `json:"p1"`
	P2            *Player  `json:"p2"`
	CurrentPlayer int8     `json:"currentPlayer"`
	Board         [][]int8 `json:"board"`
	State         string   `json:"state"`  // 'pending' | 'p1' | 'p2' | 'tie'
	Size          int8     `json:"size"`   // size of the board
	Length        int8     `json:"length"` // length to win
}

func (t *TicTacToe) Init() {
	if t.Size == 0 {
		t.Size = 3
	}
	if t.Length == 0 {
		t.Length = 3
	}

	board := make([][]int8, t.Size)
	for i := range board {
		board[i] = make([]int8, t.Size)
		for j := range board[i] {
			board[i][j] = 0
		}
	}

	t.Board = board

	if t.CurrentPlayer == 2 {
		t.CurrentPlayer = 2
	} else {
		t.CurrentPlayer = 1
	}

	t.State = "pending"
}

func (t *TicTacToe) Play(u *Player, x, y int) int8 {
	if t.Board[y][x] != 0 || t.State != "pending" {
		return 0
	}

	if t.CurrentPlayer == 1 {
		if t.P1 == nil || t.P1.Token != u.Token {
			return 0
		}
	}

	if t.CurrentPlayer == 2 {
		if t.P2 == nil || t.P2.Token != u.Token {
			return 0
		}
	}

	t.Board[y][x] = t.CurrentPlayer

	defer func() {
		if t.CurrentPlayer == 1 {
			t.CurrentPlayer = 2
		} else {
			t.CurrentPlayer = 1
		}
	}()

	t.CheckWin()

	return t.CurrentPlayer
}

func (t *TicTacToe) CheckWin() {
	for i := 0; i < 3; i++ {
		if t.Board[i][0] != 0 && t.Board[i][0] == t.Board[i][1] && t.Board[i][1] == t.Board[i][2] {
			t.State = "p" + strconv.Itoa(int(t.Board[i][0]))
			return
		}
	}

	for j := 0; j < 3; j++ {
		if t.Board[0][j] != 0 && t.Board[0][j] == t.Board[1][j] && t.Board[1][j] == t.Board[2][j] {
			t.State = "p" + strconv.Itoa(int(t.Board[0][j]))
			return
		}
	}

	if t.Board[0][0] != 0 && t.Board[0][0] == t.Board[1][1] && t.Board[1][1] == t.Board[2][2] {
		t.State = "p" + strconv.Itoa(int(t.Board[0][0]))
		return
	}

	if t.Board[0][2] != 0 && t.Board[0][2] == t.Board[1][1] && t.Board[1][1] == t.Board[2][0] {
		t.State = "p" + strconv.Itoa(int(t.Board[0][2]))
		return
	}

	tie := true
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if t.Board[i][j] == 0 {
				tie = false
				break
			}
		}
		if !tie {
			break
		}
	}
	if tie {
		t.State = "tie"
		return
	}

	t.State = "pending"
	return
}
