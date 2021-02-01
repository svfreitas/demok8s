package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
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
		tmpl.Execute(w, data)
	})

	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		if healthzIsBad {
			w.WriteHeader(500)
			w.Write([]byte(":-("))
		} else {
			w.WriteHeader(200)
			w.Write([]byte(":-)"))
		}
	})

	http.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
		duration := time.Since(started)
		if duration.Seconds() > 10 {
			w.WriteHeader(200)
			w.Write([]byte("I am ready!!!"))
		} else {
			w.WriteHeader(500)
			w.Write([]byte(fmt.Sprintf("Need to wait  %v seconds", 10-duration.Seconds())))
		}
	})

	http.HandleFunc("/infected", func(w http.ResponseWriter, r *http.Request) {
		healthzIsBad = true
		w.WriteHeader(200)
		w.Write([]byte(fmt.Sprintf("The %s container was infected", hostname)))

	})

	http.ListenAndServe(":80", nil)
}
