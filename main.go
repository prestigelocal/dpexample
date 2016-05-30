package main

import (
	"fmt"
	"github.com/deferpanic/deferclient/deferstats"
	"io/ioutil"
	"net/http"
	"time"
)

// fast test
func fastHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "this request is fast")

	resp, err := http.Get("http://10.0.0.2:3000")
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
		fmt.Fprintf(w, "downloaded")
	}
}

// slow test
func slowHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(6 * time.Second)
	fmt.Fprintf(w, "this request is slow")
}

// panic test
func panicHandler(w http.ResponseWriter, r *http.Request) {
	panic("There is no need for panic")
	fmt.Fprintf(w, "this request is panic")
}

func main() {
	dps := deferstats.NewClient("z57z3xsEfpqxpr0dSte0auTBItWBYa1c")

	go dps.CaptureStats()

	// no need to change these?
	http.HandleFunc("/fast", dps.HTTPHandlerFunc(fastHandler))
	http.HandleFunc("/slow", dps.HTTPHandlerFunc(slowHandler))
	http.HandleFunc("/panic", dps.HTTPHandlerFunc(panicHandler))

	http.ListenAndServe(":3000", nil)
}
