package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"github.com/deferpanic/deferclient/deferstats"
)


func mlbPing(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "mlb")

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
	fmt.Fprintf(w, string(body))
}

func esPing(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "es")

	resp, err := http.Get("http://212.47.234.190:9200")
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

func main() {
	dps := deferstats.NewClient("z57z3xsEfpqxpr0dSte0auTBItWBYa1c")
	go dps.CaptureStats()

	http.HandleFunc("/es", dps.HTTPHandlerFunc(esPing))
	http.HandleFunc("/mlb", dps.HTTPHandlerFunc(mlbPing))
	http.ListenAndServe(":3000", nil)
}
