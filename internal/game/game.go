package game

import (
	"fmt"
)

type Game struct {
	players []Player
}

func (g *Game) Players() []Player {
	return g.players
}

func (g *Game) AddPlayer(p Player) {
	g.players = append(g.players, p)
}

func (g *Game) ListPlayers() {
	for i := 0; i < len(g.players); i++ {
		fmt.Println(g.players[i].Name())
	}

}

type Player struct {
	name string
	loc  Location
}

func (p *Player) Name() string {
	return p.name
}

type Location struct {
	x, y float32
}
