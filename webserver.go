package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"sync"
	"time"
)

const version = "1.0"
const color = "green"

type PageData struct {
	PageTitle string
	Hostname  string
	Color     string
	Version   string
}

const Layout string = `
<html>
    <head>
        <style>
        h1 {text-align: center;background-color:{{.Color}};}
        h2 {text-align: center;}
        </style>
    </head>
    <body>
        <h1>{{.PageTitle}}</h1>
        <h2>This request was processed by : {{.Hostname}}</h2>
    </body>
</html>
`

func main() {

	started := time.Now()

	var mu sync.RWMutex

	healthzIsBad := false

	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	tmpl, err := template.New("template").Parse(Layout)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		data := PageData{
			PageTitle: "Containers Demo " + version,
			Hostname:  hostname,
			Color:     color,
		}

		// simulate some process time
		time.Sleep(25 * time.Millisecond)

		tmpl.Execute(w, data)
	})

	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		mu.RLock()
		if healthzIsBad {
			w.WriteHeader(500)
			w.Write([]byte(":-("))
		} else {
			w.WriteHeader(200)
			w.Write([]byte(":-)"))
		}
		mu.RUnlock()
	})

	http.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
		duration := time.Since(started)
		if duration.Seconds() > 5 {
			w.WriteHeader(200)
			w.Write([]byte("I am ready!!!"))
		} else {
			w.WriteHeader(500)
			w.Write([]byte(fmt.Sprintf("Need to wait  %v seconds", 10-duration.Seconds())))
		}
	})

	http.HandleFunc("/degraded", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		healthzIsBad = !healthzIsBad
		w.WriteHeader(200)
		if healthzIsBad {
			w.Write([]byte(fmt.Sprintf("The %s container is  degraded", hostname)))
		} else {
			w.Write([]byte(fmt.Sprintf("The %s container is  restored", hostname)))
		}
		mu.Unlock()
	})

	http.HandleFunc("/kill", func(w http.ResponseWriter, r *http.Request) {
		os.Exit(1)
	})

	http.ListenAndServe(":80", nil)
}
