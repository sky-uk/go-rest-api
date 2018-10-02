package main

import (
	"flag"
	"fmt"
	"github.com/sky-uk/go-rest-api"
	"net/http"
	"strconv"
)

func main() {
	var server string
	var path string
	var port int

	flag.StringVar(&server, "server", "localhost", "the toxiproxy server IP or FQDN")
	flag.IntVar(&port, "port", 8474, "the toxiproxy server port")
	flag.StringVar(&path, "path", "/", "the query path")

	flag.Parse()

	url := "http://" + server + ":" + strconv.Itoa(port)

	client := rest.Client{URL: url}
	api := rest.NewBaseAPI(http.MethodGet, path, nil, new(string), nil)
	err := client.Do(api)
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
	resp := *api.ResponseObject().(*string)
	fmt.Println("Response:\n", resp)
}
