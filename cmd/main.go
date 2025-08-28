package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/isaacp84/go_server_1/internal/game"
)

// var udp_running atomic.Bool

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
	fmt.Println("started connection")
	return conn, addr
}

func ReadUDPPackets(ctx context.Context, packetChan chan []byte, conn *net.UDPConn) {
	buffer := make([]byte, 512)

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Server context cancelled, stopping UDP listener.")
			conn.Close()
			return
		default:
			if ctx.Done() == nil {
				continue
			}

			readDeadline := time.Now().Add(50 * time.Millisecond)
			if err := conn.SetReadDeadline(readDeadline); err != nil {
				fmt.Println("Error setting read deadline:", err)
				return
			}
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

			// Send a copy of the received data to avoid data races
			data := make([]byte, n)
			copy(data, buffer[:n])
			packetChan <- data
			fmt.Printf("Received %d bytes from %s: %s\n", n, addr.String(), string(data))
		}

	}
}

func run(ctx context.Context, w io.Writer, args []string) error {
	ctx, cancel := context.WithCancel(ctx)

	var wg sync.WaitGroup

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	fmt.Println("Hello world")

	var pong game.Game

	pong.Players()

	// Channel to send received UDP packets
	packetChan := make(chan []byte)
	conn, _ := startUDPConn()
	// Close the connection when we're done
	defer conn.Close()

	pong.LinkUDP(conn)

	// Channel to signal shutdown
	// doneChan := make(chan struct{})

	// Accept incoming connections and handle them
	// Goroutine to read UDP packets
	wg.Go(func() {
		ReadUDPPackets(ctx, packetChan, conn)
	})

	// Game loop
	wg.Go(func() {
		pong.Run(ctx, packetChan)
	})

	<-sigChan
	fmt.Println("Received termination signal, shutting down...")
	cancel()

	select {
	case <-time.After(5 * time.Second):
		fmt.Println("Graceful shutdown timeout exceeded. Exiting forcefully.")
	case <-ctx.Done():
		wg.Wait()
		fmt.Println("Server gracefully shutdown.")
	}

	return nil
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := run(ctx, os.Stdout, os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		cancel()
		os.Exit(1)
	}

}

// func handleConnection(conn *net.UDPConn, addr *net.UDPAddr, buf []byte) {
// 	// Print the incoming data
// 	fmt.Print("> ", string(buf[0:]))
// 	// Write back the message over UPD
// 	conn.WriteToUDP([]byte("Hello UDP Client\n"), addr)
// }
