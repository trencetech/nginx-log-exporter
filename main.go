package main

import (
	"crypto/subtle"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/songjiayang/nginx-log-exporter/collector"
	"github.com/songjiayang/nginx-log-exporter/config"
)

func main() {
	var listenAddress, configFile string
	var username, password string
	var placeholderReplace bool

	flag.StringVar(&listenAddress, `web.listen-address`, `:9999`, `Address to listen on for the web interface and API.`)
	flag.StringVar(&configFile, `config.file`, `config.yml`, `Nginx log exporter configuration file name.`)
	flag.StringVar(&username, `web.username`, `jogodo`, `Nginx log exporter username for basic auth.`)
	flag.StringVar(&password, `web.password`, `jegedee`, `Nginx log exporter password for basic auth.`)
	flag.BoolVar(&placeholderReplace, `placeholder.replace`, false, `Enable placeholder replacement when rewriting the request path.`)
	flag.Parse()

	cfg, err := config.LoadFile(configFile)
	if err != nil {
		log.Panic(err)
	}

	var options config.Options
	options.SetPlaceholderReplace(placeholderReplace)

	for _, app := range cfg.App {
		go collector.NewCollector(app, options).Run()
	}

	fmt.Printf("running HTTP server on address %s\n", listenAddress)

	metricsHandler := promhttp.HandlerFor(
		prometheus.DefaultGatherer,
		promhttp.HandlerOpts{},
	)
	protectedHandler := basicAuth(metricsHandler, username, password, "Metrics")
	http.Handle("/metrics", protectedHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		fmt.Fprintf(w, "Welcome to the home page!")
	})
	if err := http.ListenAndServe(listenAddress, nil); err != nil {
		log.Fatalf("start server with error: %v\n", err)
	}
}

func basicAuth(handler http.Handler, username, password, realm string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()

		if !ok || subtle.ConstantTimeCompare([]byte(user), []byte(username)) != 1 ||
			subtle.ConstantTimeCompare([]byte(pass), []byte(password)) != 1 {
			w.Header().Set("WWW-Authenticate", `Basic realm="`+realm+`"`)
			w.WriteHeader(401)
			w.Write([]byte("Unauthorized.\n"))
			return
		}

		handler.ServeHTTP(w, r)
	})
}
