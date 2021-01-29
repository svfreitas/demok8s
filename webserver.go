package main

import (
	"html/template"
	"net/http"
	"os"
)

const version = "1.0"

type PageData struct {
	PageTitle string
	Hostname  string
	RemoteIP  string
	Version   string
}

const Layout string = `
<html>
    <head>
        <style>
        h2 {text-align: center;}
        h1 {text-align: center;background-color:green;"}
        </style>
    </head>
    <body>
        <h1>{{.PageTitle}}</h1>
        <h2>Hostname : {{.Hostname}}</h2>
    </body>
</html>`

func main() {
	tmpl, err := template.New("template").Parse(Layout)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		hostname, err := os.Hostname()
		if err != nil {
			panic(err)
		}

		data := PageData{
			PageTitle: "Containers Demo " + version,
			Hostname:  hostname,
		}
		tmpl.Execute(w, data)
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("ok"))
	})

	http.HandleFunc("/die", func(w http.ResponseWriter, r *http.Request) {
		os.Exit(3)
	})

	http.ListenAndServe(":80", nil)
}
