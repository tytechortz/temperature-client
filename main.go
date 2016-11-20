package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	input   = flag.String("input", "/sys/bus/w1/devices/28-00000620a282/w1_slave", "path to temperature file")
	address = flag.String("address", "http://127.0.0.1:8080", "temperature server address")
)

func init() {
	flag.Parse()
	http.DefaultClient.Timeout = 15 * time.Second
}

func main() {
	for {
		t := getTemperature(*input)
		saveTemperature(t)
		time.Sleep(10 * time.Second)
	}
}

func saveTemperature(t float64) {
	m := map[string]float64{"temperature": t}
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(m); err != nil {
		log.Println("couldn't marshal data", err)
		return
	}
	resp, err := http.Post(fmt.Sprintf("%s/save", *address), "application/json", &buf)
	if err != nil {
		log.Println("couldn't POST data", err)
		return
	}
	resp.Body.Close()
}

func getTemperature(pth string) float64 {
	d, err := ioutil.ReadFile(pth)
	if err != nil {
		log.Println("error, couldn't open file:", err)
		return 0.0
	}
	i := bytes.Index(d, []byte("t="))
	t, err := strconv.ParseInt(strings.TrimSpace(string(d[i+2:])), 10, 64)
	if err != nil {
		log.Println("error, couldn't parse data:", string(d), err)
		return 0.0
	}
	return float64(t) / 1000.0
}
