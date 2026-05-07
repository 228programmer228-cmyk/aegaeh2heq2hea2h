package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
	addr := flag.String("addr", "", "Target address (ip:port)")
	mode := flag.String("mode", "client", "Mode")
	conns := flag.Int("c", 1, "Connections count")
	duration := flag.String("d", "3600s", "Duration")
	ddnet := flag.Bool("ddnet-info", false, "DDNet mode")
	
	flag.Int("n", 0, "unused")
	flag.Int("r", 0, "unused")
	flag.String("t", "1s", "unused")
	flag.Bool("open-loop", false, "unused")

	flag.Parse()

	if *addr == "" {
		fmt.Println("No address specified")
		os.Exit(1)
	}

	d, _ := time.ParseDuration(*duration)
	fmt.Printf("Starting in mode %s on %s for %s with %d conns (DDNet: %v)\n", *mode, *addr, *duration, *conns, *ddnet)

	deadline := time.Now().Add(d)

	for i := 0; i < *conns; i++ {
		go func() {
			conn, err := net.Dial("udp", *addr)
			if err != nil {
				return
			}
			defer conn.Close()

			payload := []byte("\xff\xff\xff\xffgetinfo") // Базовый пакет для DDNet
			if !*ddnet {
				payload = []byte("standard-udp-payload-test")
			}

			for time.Now().Before(deadline) {
				conn.Write(payload)
			}
		}()
	}

	time.Sleep(d)
}
