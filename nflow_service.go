package main

import (
	"fmt"
	"log"
	"net/http"
)

//Listening IP for the web server
const ip string = "0.0.0.0"

//Listening port for web server
const port string = "8080"

var finished_generating bool = true

type generate_context_t struct {
	collector   string
	command     chan int
	is_running  bool
	spike_proto Proto
}

func start_server(collector string) {
	log.Println("Starting netflow service: " + collector)

	// Rest API handlers
	gen_ctx := &generate_context_t{collector: collector, command: make(chan int),
		is_running: false, spike_proto: PROTO_NONE}

	http.HandleFunc("/api/v1/generate/start", gen_ctx.generate_start)
	http.HandleFunc("/api/v1/generate/stop", gen_ctx.generate_stop)
	http.HandleFunc("/api/v1/status", gen_ctx.get_status)

	log.Printf("Starting web server on: %s:%s\n", ip, port)
	log.Fatal(http.ListenAndServe(ip+":"+port, nil))
}

func (g_ctx *generate_context_t) generate_start(writer http.ResponseWriter, request *http.Request) {
	if !g_ctx.is_running {
		go generate(g_ctx)
	}
	return
}

func (g_ctx *generate_context_t) generate_stop(writer http.ResponseWriter, request *http.Request) {
	if g_ctx.is_running {
		go func() {
			g_ctx.command <- int(0)
		}()
	}
}

func (g_ctx *generate_context_t) get_status(writer http.ResponseWriter, request *http.Request) {
	if g_ctx.is_running {
		fmt.Fprintln(writer, "running")
	} else {
		fmt.Fprintln(writer, "stopped")
	}

	return
}
