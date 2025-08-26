package game

import (
	"fmt"
	"net"
)

type Game struct {
	players []Player
	Ball    Ball
	Paddles []Paddle
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

type Ball struct {
	Location *Location
}

type Paddle struct {
	Location *Location
}

type Player struct {
	Name string
	Addr *net.UDPAddr
}

type Location struct {
	x, y float32
}
