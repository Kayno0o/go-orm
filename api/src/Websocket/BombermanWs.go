package ws

import (
	"fmt"
	"strconv"
	"time"

	utils "go-api-test.kayn.ooo/src/Utils"
)

type BombermanWs struct {
	GenericWsI
}

func (ws *BombermanWs) Init() {
	Handle("bomberman", func(u *Player) *BombermanRoom {
		room := &BombermanRoom{
			Room: Room[Bomberman]{
				Data: Bomberman{
					Width:  15,
					Height: 15,
				},
			},
		}
		room.Data.Init()
		return room
	})
}

type BombermanRoom struct {
	Room[Bomberman]
}

func (r *BombermanRoom) Quit(u *Player) {
	player, i := r.Data.GetPlayerIndex(u.Token)
	if player == nil {
		return
	}

	utils.RemoveAtIndex(&r.Data.Players, i)
	r.Add <- Update{"delete", "data.players." + strconv.Itoa(i), nil}
}

func (r *BombermanRoom) UpdateUser(u *Player) {
}

func (r *BombermanRoom) HandleResponse(u *Player, message ClientMessage) {
	if r.Data.Add == nil {
		r.Data.Add = r.Add
	}

	switch message.Type {
	case "bomb":
		player, _ := r.Data.GetPlayerIndex(u.Token)
		if player == nil {
			fmt.Println("no player")
			return
		}

		if r.Data.CountPlayerBombs(player) >= int(player.MaxBomb) {
			return
		}

		r.Data.PlaceBomb(player)
	case "direction":
		player, _ := r.Data.GetPlayerIndex(u.Token)
		if player == nil {
			fmt.Println("no player")
			return
		}

		dir := message.Content

		switch dir {
		case "up":
			player.Direction = 0
		case "right":
			player.Direction = 1
		case "down":
			player.Direction = 2
		case "left":
			player.Direction = 3
		default:
			player.Direction = 255 // no movement
		}

		if player.CanMove {
			player.MoveRoutine <- true
		}
		break

	case "join":
		player, _ := r.Data.GetPlayerIndex(u.Token)
		if len(r.Data.Players) >= 4 || player != nil {
			return
		}

		poses := []Position{
			{0, 0},
			{int(r.Data.Width) - 1, int(r.Data.Height) - 1},
			{int(r.Data.Width) - 1, 0},
			{0, int(r.Data.Height) - 1},
		}

		i := len(r.Data.Players)

		player = &BombermanPlayer{Player: u, Speed: 300, Pos: poses[i], MoveRoutine: nil, Power: 3, Direction: 255, MaxBomb: 1}

		player.MoveRoutine = make(chan bool)
		go utils.SetInterval(time.Duration(player.Speed)*time.Millisecond, player.MoveRoutine, func(startup bool) {
			hasMoved := r.Data.MovePlayer(player)
			if hasMoved {
				r.Add <- Update{"update", "data.players." + strconv.Itoa(i), player}
				player.CanMove = false
			} else {
				player.CanMove = true
			}
		})

		r.Data.Players = append(r.Data.Players, player)
		r.Add <- Update{"push", "data.players", player}
		break

	case "quit":
		player, _ := r.Data.GetPlayerIndex(u.Token)
		if player != nil {
			player.MoveRoutine <- false
			close(player.MoveRoutine)
		}

		r.Quit(u)
		break
	}
}
