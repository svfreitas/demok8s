package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"sync"
	"time"
)

const version = "2.0"
const color = "blue"

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

var mu sync.RWMutex
var healthzIsBad bool

func main() {

	started := time.Now()
	healthzIsBad = false

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

		// simulate some process time, database access,...
		time.Sleep(25 * time.Millisecond)

		tmpl.Execute(w, data)
	})

	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		mu.RLock()
		if healthzIsBad {
			w.WriteHeader(500)
			w.Write([]byte("<html><h1>:-(</h1><html>"))
		} else {
			w.WriteHeader(200)
			w.Write([]byte("<html><h1>:-)</h1><html>"))
		}
		mu.RUnlock()
	})

	http.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
		duration := time.Since(started)
		if duration.Seconds() > 5 {
			w.WriteHeader(200)
			w.Write([]byte("<html><h1>I am ready!!!</h1></html>"))
		} else {
			w.WriteHeader(500)
			w.Write([]byte(fmt.Sprintf("<html><h1>Need to wait  %v seconds</h1></html>", 10-duration.Seconds())))
		}
	})

	http.HandleFunc("/degraded", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		healthzIsBad = true
		w.WriteHeader(200)
		w.Write([]byte(fmt.Sprintf("The %s container is  degraded", hostname)))
		go ResetDegradation()
		mu.Unlock()
	})

	http.HandleFunc("/crash", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(fmt.Sprintf("The %s container is  allocating memory", hostname)))
		a := make([]byte, 50*1024*1024)
		for i := 0; i < len(a); i++ {
			a[i] = 'x'
			if i%1024 == 0 {
				time.Sleep(500 * time.Millisecond)
			}
		}
	})

	http.ListenAndServe(":80", nil)
}

func ResetDegradation() {
	time.Sleep(3 * time.Second)
	mu.Lock()
	healthzIsBad = false
	mu.Unlock()
}
