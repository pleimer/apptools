package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type Grafana struct {
	address string
	apiBase string
	headers string
}

//NewGrafana grafana factory
func NewGrafana(address string) Grafana {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	return Grafana{
		address: address,
		apiBase: "/api/dashboards/uid/",
		headers: "Content-Type: aplication/json",
	}
}

func (g Grafana) getDashboard(uid string) ([]byte, error) {
	resp, err := http.Get("https://" + g.address + g.apiBase + uid)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s -uid <dashboard uid> [options]\n", os.Args[0])
		flag.PrintDefaults()
	}

	apiAddr := flag.String("address", "localhost:3000", "Grafana API server address")
	dbuid := flag.String("uid", "", "dashboard UID to pull")
	output := flag.String("output", "dashboard.yml", "Location for dashboard output")
	flag.Parse()

	if *dbuid == "" {
		flag.Usage()
		os.Exit(1)
	}

	grafana := NewGrafana(*apiAddr)

	dbData, err := grafana.getDashboard(*dbuid)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s", err)
	}

	jsonMap := make(map[string]interface{})

	err = json.Unmarshal(dbData, &jsonMap)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s", err)
	}

	dbJSON := jsonMap["dashboard"].(map[string]interface{})

	dbJSON["id"] = nil

	dbData, err = json.MarshalIndent(dbJSON, "", "    ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s", err)
	}

	ioutil.WriteFile(*output, dbData, 0644)
}
