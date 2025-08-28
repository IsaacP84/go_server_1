package game

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"
)

type Game struct {
	players []Player
	Ball    Ball
	Paddles []Paddle
	conn    *net.UDPConn
}

func (g *Game) LinkUDP(c *net.UDPConn) {
	g.conn = c
}

func (g *Game) Run(ctx context.Context, packetChan chan []byte) error {
	if g.conn == nil {
		return nil
	}

	tickDuration := time.Second / 10
	ticker := time.NewTicker(tickDuration)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Server context cancelled, stopping game routine.")
			return nil
		case packet := <-packetChan:
			g.handleUDPPacket(packet)

		case <-ticker.C:
			// Game code
			g.Update(g.conn)

			// fmt.Println("Game tick.")
		}
	}
}

func (g *Game) handleUDPPacket(packet []byte) {
	fmt.Printf("Processing packet: %s\n", string(packet))
	switch string(packet[:4]) {
	case "JN  ":
		name := string(packet[4:19])
		p := Player{Name: name}
		fmt.Println(p)
		g.AddPlayer(p)
		log.Printf("Player joined: %s", name)

	case "LV  ":
	default:
	}
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
