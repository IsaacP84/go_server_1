package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/isaacp84/go_server_1/cmd/server"
	"github.com/isaacp84/go_server_1/internal/game"
)

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
			return
		default:
			readDeadline := time.Now().Add(50 * time.Millisecond)
			if err := conn.SetReadDeadline(readDeadline); err != nil {
				fmt.Println("Error setting read deadline:", err)
				return
			}
			// Read incoming data
			n, remoteAddr, err := conn.ReadFromUDP(buffer)

			if err != nil {
				// Handle temporary errors or log persistent ones
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					// fmt.Println("Read timeout, still listening...")
					continue // It's a timeout, try again
				}

				if err == os.ErrClosed || strings.Contains(err.Error(), "use of closed network connection") {
					fmt.Println("UDP connection closed, stopping read loop.")
					break // Exit the loop if the connection is closed
				}

				log.Printf("Error reading from UDP: %v", err)
				break
			}

			// Send a copy of the received data to avoid data races
			data := make([]byte, n)
			copy(data, buffer[:n])

			packetChan <- data
			fmt.Printf("Received %d bytes from %s: %s\n", n, remoteAddr.String(), string(data))

			response := []byte("ACK: " + string(data[:n]))
			_, err = conn.WriteToUDP(response, remoteAddr)
			if err != nil {
				fmt.Printf("Error writing to UDP: %v\n", err)
			}
		}
	}
}

func run(p_ctx context.Context, w io.Writer, args []string) error {
	ctx, cancel := context.WithCancel(p_ctx)

	var wg sync.WaitGroup

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	fmt.Println("Hello world")

	var pong game.Game

	pong.Players()

	srv := server.NewServer()

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: srv,
	}

	go func() {
		log.Printf("listening on %s\n", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "error listening and serving: %s\n", err)
		}

	}()
	wg.Go(func() {
		<-ctx.Done()
		shutdownCtx := context.Background()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			fmt.Fprintf(os.Stderr, "error shutting down http server: %s\n", err)
		}
		fmt.Println("HTTP server stopped.")
	})

	// Channel to send received UDP packets
	packetChan := make(chan []byte)
	conn, _ := startUDPConn()
	// Close the connection when we're done
	// defer conn.Close()

	pong.LinkUDP(conn)

	// Channel to signal shutdown
	// doneChan := make(chan struct{})

	// Accept incoming connections and handle them
	// Goroutine to read UDP packets
	go func() {
		ReadUDPPackets(ctx, packetChan, conn)
	}()
	wg.Go(func() {
		<-ctx.Done()
		conn.Close()
	})

	// Game loop
	wg.Go(func() {
		<-ctx.Done()
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
	ctx := context.Background()
	if err := run(ctx, os.Stdout, os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

}
