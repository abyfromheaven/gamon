package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"time"

	"gamon/ping"
)

const banner = `
===========================================
   GAMON - Garda Monitoring v0.1
   Realtime Network Monitor (CLI)
===========================================`

func main() {
	fmt.Println(banner)
	fmt.Println()

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Masukkan IP Target: ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if net.ParseIP(input) == nil {
		// Coba resolve sebagai domain
		_, err := net.LookupHost(input)
		if err != nil {
			fmt.Printf("\n[ERROR] '%s' bukan IP address atau domain yang valid!\n", input)
			os.Exit(1)
		}
	}

	// Handle CTRL+C
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	fmt.Printf("\nMonitoring %s dimulai... (CTRL+C untuk berhenti)\n\n", input)

	seq := 0

	// Ping pertama langsung
	seq++
	result := ping.PingOnce(input)
	result.Seq = seq
	fmt.Println(ping.FormatOutput(result))

	for {
		select {
		case <-sigChan:
			fmt.Printf("\n\nMonitoring dihentikan. Total ping: %d\n", seq)
			os.Exit(0)
		case <-ticker.C:
			seq++
			result := ping.PingOnce(input)
			result.Seq = seq
			fmt.Println(ping.FormatOutput(result))
		}
	}
}
