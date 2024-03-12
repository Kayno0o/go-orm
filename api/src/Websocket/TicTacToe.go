package ws

import (
	"encoding/json"
)

type TicTacToe struct {
	P1            *User            `json:"p1"`
	P2            *User            `json:"p2"`
	CurrentPlayer int8             `json:"currentPlayer"`
	Spectators    map[string]*User `json:"spectators"`
	Board         [][]int8         `json:"board"`
	State         string           `json:"state"` // 'pending' | 'p1' | 'p2' | 'tie'
}

func (t *TicTacToe) Init() {
	board := make([][]int8, 3)
	for i := range board {
		board[i] = make([]int8, 3)
		for j := range board[i] {
			board[i][j] = 0
		}
	}

	t.Board = board
	t.CurrentPlayer = 1
	t.State = "pending"
}

func (t *TicTacToe) Play(x, y int) {
	if t.Board[y][x] != 0 || t.State != "pending" {
		return
	}

	t.Board[y][x] = t.CurrentPlayer

	if t.CurrentPlayer == 1 {
		t.CurrentPlayer = 2
	} else {
		t.CurrentPlayer = 1
	}

	t.CheckWin()
}

func (t *TicTacToe) CheckWin() {
	for i := 0; i < 3; i++ {
		if t.Board[i][0] != 0 && t.Board[i][0] == t.Board[i][1] && t.Board[i][1] == t.Board[i][2] {
			if t.Board[i][0] == 1 {
				t.State = "p1"
			} else {
				t.State = "p2"
			}
			return
		}
	}

	for j := 0; j < 3; j++ {
		if t.Board[0][j] != 0 && t.Board[0][j] == t.Board[1][j] && t.Board[1][j] == t.Board[2][j] {
			if t.Board[0][j] == 1 {
				t.State = "p1"
			} else {
				t.State = "p2"
			}
			return
		}
	}

	if t.Board[0][0] != 0 && t.Board[0][0] == t.Board[1][1] && t.Board[1][1] == t.Board[2][2] {
		if t.Board[0][0] == 1 {
			t.State = "p1"
		} else {
			t.State = "p2"
		}
		return
	}

	if t.Board[0][2] != 0 && t.Board[0][2] == t.Board[1][1] && t.Board[1][1] == t.Board[2][0] {
		if t.Board[0][2] == 1 {
			t.State = "p1"
		} else {
			t.State = "p2"
		}
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
}

func (t *TicTacToe) SendState() {
	for _, s := range t.Spectators {
		_ = s.SendMessage(t.JSON())
	}
}

func (t *TicTacToe) JSON() string {
	str, err := json.Marshal(t)
	if err != nil {
		return ""
	}
	return string(str)
}
