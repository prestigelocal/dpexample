package main

import (
	"fmt"
	"github.com/deferpanic/deferclient/deferstats"
	"github.com/takama/daemon"
	"io/ioutil"
	"net/http"
	"log"
	"os"
	"os/signal"
	"syscall"
	"net"
)

const (
	port = ":9977"
)

var stdlog, errlog *log.Logger

type Service struct {
	daemon.Daemon
}

func fastHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "this request is fast")

	resp, err := http.Get("http://gd2.mlb.com/components/game/mlb/year_2016/month_08/day_07/master_scoreboard.json")
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	if len(body) != 0 {
		fmt.Fprintf(w, string(body))
	}
}

func (service *Service) Manage() (string, error) {
	usage := "Usage: myservice install | remove | start | stop | status"

	if len(os.Args) > 1 {
		command := os.Args[1]
		switch command {
		case "install":
			return service.Install()
		case "remove":
			return service.Remove()
		case "start":
			return service.Start()
		case "stop":
			return service.Stop()
		case "status":
			return service.Status()
		default:
			return usage, nil
		}
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill, syscall.SIGTERM)

	listener, err := net.Listen("tcp", port)
	if err != nil {
		return "Possibly was a problem with the port binding", err
	}

	listen := make(chan net.Conn, 100)
	go acceptConnection(listener, listen)

	for {
		select {
		case conn := <-listen:
			go handleClient(conn)
		case killSignal := <-interrupt:
			stdlog.Println("Got signal:", killSignal)
			stdlog.Println("Stoping listening on ", listener.Addr())
			listener.Close()
			if killSignal == os.Interrupt {
				return "Daemon was interruped by system signal", nil
			}
			return "Daemon was killed", nil
		}
	}

	// never happen, but need to complete code
	return usage, nil
}

func acceptConnection(listener net.Listener, listen chan<- net.Conn) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		listen <- conn
	}
}

func handleClient(client net.Conn) {
	for {
		buf := make([]byte, 4096)
		numbytes, err := client.Read(buf)
		if numbytes == 0 || err != nil {
			return
		}
		client.Write(buf[:numbytes])
	}
}

func init() {
	stdlog = log.New(os.Stdout, "", log.Ldate|log.Ltime)
	errlog = log.New(os.Stderr, "", log.Ldate|log.Ltime)
}


func main() {
	dps := deferstats.NewClient("z57z3xsEfpqxpr0dSte0auTBItWBYa1c")
	go dps.CaptureStats()

	service, err := daemon.New("name", "description")

	if err != nil {
		log.Fatal("Error: ", err)
	}
	status, err := service.Install()
	fmt.Println(status)

	http.HandleFunc("/health", dps.HTTPHandlerFunc(fastHandler))
	http.ListenAndServe(":3000", nil)
}
