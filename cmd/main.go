package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/isaacp84/go_server_1/internal/game"
)

func main() {
	fmt.Println("Hello world")
	var game game.Game

	tickDuration := time.Second / 20
	ticker := time.NewTicker(tickDuration)
	defer ticker.Stop()

	game.Players()

	conn := startUDPConn()
	// Close the connection when we're done
	defer conn.Close()

	// Channel to send received UDP packets
	packetChan := make(chan []byte)
	// Channel to signal shutdown
	doneChan := make(chan struct{})

	// Accept incoming connections and handle them
	// Goroutine to read UDP packets
	go func() {
		buffer := make([]byte, 512)
		for {
			select {
			case <-doneChan:
				fmt.Println("UDP reader goroutine exiting.")
				return
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
	for {
		select {

		case packet := <-packetChan:
			fmt.Printf("Processing packet: %s\n", string(packet))

		case <-ticker.C:
			// Game code
			// fmt.Println("Game tick.")

		case <-time.After(3 * time.Second):
			fmt.Println("Game stopped.")
			close(doneChan) // Signal the reader goroutine to exit

			return
		}
	}

}

func startUDPConn() *net.UDPConn {
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
		return nil
	}

	fmt.Println("started connection")
	return conn

}

// func handleConnection(conn *net.UDPConn, addr *net.UDPAddr, buf []byte) {
// 	// Print the incoming data
// 	fmt.Print("> ", string(buf[0:]))
// 	// Write back the message over UPD
// 	conn.WriteToUDP([]byte("Hello UDP Client\n"), addr)
// }
