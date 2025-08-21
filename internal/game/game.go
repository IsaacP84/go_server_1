package game

import (
	"fmt"
	"net"
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
		fmt.Println(g.players[i].Name)
	}

}

func (g *Game) Update(conn *net.UDPConn) {
	for i := 0; i < len(g.players); i++ {
		conn.WriteToUDP([]byte("TK: "), g.players[i].Addr)
	}
}

type Player struct {
	Name string
	Addr *net.UDPAddr
	loc  *Location
}

// func (p *Player) Name() string {
// 	return p.name
// }

// func (p *Player) SetName(s string) {
// 	p.name = s
// }

type Location struct {
	x, y float32
}
