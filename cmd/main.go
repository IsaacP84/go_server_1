package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/isaacp84/go_server_1/internal/game"
)

// var udp_running atomic.Bool

func run(ctx context.Context, w io.Writer, args []string) error {
	sigChan, _ := signal.NotifyContext(ctx, os.Interrupt)

	fmt.Println("Hello world")

	// http.Handle("/", routers.Root{})
	// log.Fatal(http.ListenAndServe(":8080", nil))

	var my_game game.Game
	tickDuration := time.Second / 10
	ticker := time.NewTicker(tickDuration)
	defer ticker.Stop()

	my_game.Players()

	conn, _ := startUDPConn()
	// Close the connection when we're done
	defer conn.Close()

	// Channel to send received UDP packets
	packetChan := make(chan []byte)
	// Channel to signal shutdown
	// doneChan := make(chan struct{})

	// Accept incoming connections and handle them
	// Goroutine to read UDP packets
	go func() {

		buffer := make([]byte, 512)
		for {
			select {
			case <-ctx.Done():
				fmt.Println("Server context cancelled, stopping UDP listener.")
				return
			// case <-doneChan:
			// 	fmt.Println("UDP reader goroutine exiting.")
			// 	conn.Close()
			// 	return
			default:
				// Read incoming data
				n, addr, err := conn.ReadFromUDP(buffer)
				if err != nil {
					// Handle temporary errors or log persistent ones
					if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
						continue // It's a timeout, try again
					}
					log.Printf("Error reading from UDP: %v", err)
					continue
				}

				// handleConnection(conn, addr, buffer)

				// Send a copy of the received data to avoid data races
				data := make([]byte, n)
				copy(data, buffer[:n])
				packetChan <- data
				fmt.Printf("Received %d bytes from %s: %s\n", n, addr.String(), string(data))
			}
		}
	}()

	// Game loop
	// go func() error {
	// 	for {
	// 		select {
	// 		case <-ctx.Done():
	// 			fmt.Println("Server context cancelled, stopping game routine.")
	// 			return nil
	// 		case packet := <-packetChan:
	// 			fmt.Printf("Processing packet: %s\n", string(packet))
	// 			switch string(packet[:4]) {
	// 			case "JN  ":
	// 				name := string(packet[4:19])
	// 				p := game.Player{Name: name}
	// 				fmt.Println(p)
	// 				my_game.AddPlayer(p)
	// 				log.Printf("Player joined: %s", name)

	// 			case "LV  ":
	// 			default:
	// 			}

	// 		case <-ticker.C:
	// 			// Game code
	// 			my_game.Update(conn)

	// 			// fmt.Println("Game tick.")
	// 		}
	// 	}
	// }()

	<-sigChan.Done()
	fmt.Println("Received termination signal, shutting down...")
	// close(doneChan)
	// cancel()

	return nil
}

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := run(ctx, os.Stdout, os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	select {
	case <-time.After(5 * time.Second):
		fmt.Println("Graceful shutdown timeout exceeded. Exiting forcefully.")
	case <-ctx.Done():
		fmt.Println("Server gracefully shutdown.")
	}
}

func startUDPConn() (*net.UDPConn, *net.UDPAddr) {
	// Resolve the string address to a UDP address
	addr, err := net.ResolveUDPAddr("udp", ":27015")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Listen for incoming connections on port
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return nil, nil
	}

	// udp_running.Store(true)
	fmt.Println("started connection")
	return conn, addr

}

// func handleConnection(conn *net.UDPConn, addr *net.UDPAddr, buf []byte) {
// 	// Print the incoming data
// 	fmt.Print("> ", string(buf[0:]))
// 	// Write back the message over UPD
// 	conn.WriteToUDP([]byte("Hello UDP Client\n"), addr)
// }
