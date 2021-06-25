// Run using:
// go run nflow-generator.go nflow_logging.go nflow_payload.go  -t 172.16.86.138 -p 9995
// Or:
// go build
// ./nflow-generator -t <ip> -p <port>
package main

import (
	"log"
	"math/rand"
	"net"
	"os"
	"time"
)

type Proto int

const (
	PROTO_NONE Proto = iota
	PROTO_FTP
	PROTO_SSH
	PROTO_DNS
	PROTO_HTTP
	PROTO_HTTPS
	PROTO_NTP
	PROTO_SNMP
	PROTO_IMAPS
	PROTO_MYSQL
	PROTO_HTTPS_ALT
	PROTO_P2P
	PROTO_BITTORRENT
)

func main() {

	collector_ip := os.Getenv("COLLECTOR_IP")
	collector_port := os.Getenv("COLLECTOR_PORT")

	if collector_ip == "" || collector_port == "" {
		log.Println("Error: both COLLECTOR_IP and COLLECTOR_PORT environment variables must be set")
		os.Exit(1)
	}

	collector := collector_ip + ":" + collector_port

	start_server(collector)
}

func generate(gen_ctx *generate_context_t) {
	log.Println("Starting netflow packet generation")
	gen_ctx.is_running = true
	udpAddr, err := net.ResolveUDPAddr("udp", gen_ctx.collector)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		log.Println("Error connecting to the target collector: " + err.Error())
		os.Exit(1)
	}
	defer conn.Close()

	for {
		select {
		case command := <-gen_ctx.command:
			if command == 0 {
				gen_ctx.is_running = false
				log.Println("Stopping netflow packet generation")
				return
			}
		default:
		}

		rand.Seed(time.Now().Unix())
		n := randomNum(50, 1000)
		// add spike data
		if gen_ctx.spike_proto != PROTO_NONE {
			GenerateSpike(gen_ctx.spike_proto)
		}
		if n > 900 {
			data := GenerateNetflow(8)
			buffer := BuildNFlowPayload(data)
			_, err := conn.Write(buffer.Bytes())
			if err != nil {
				log.Println("Error connecting to the target collector: " + err.Error())
				os.Exit(1)
			}
		} else {
			data := GenerateNetflow(16)
			buffer := BuildNFlowPayload(data)
			_, err := conn.Write(buffer.Bytes())
			if err != nil {
				log.Println("Error connecting to the target collector: " + err.Error())
				os.Exit(1)
			}
		}
		// add some periodic spike data
		if n < 150 {
			sleepInt := time.Duration(3000)
			time.Sleep(sleepInt * time.Millisecond)
		}
		sleepInt := time.Duration(n)
		time.Sleep(sleepInt * time.Millisecond)
	}
}

func randomNum(min, max int) int {
	return rand.Intn(max-min) + min
}
