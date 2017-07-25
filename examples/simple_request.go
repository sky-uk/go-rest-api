package main

import (
	"fmt"
	"github.com/sky-uk/go-rest-api"
	"net/http"
)

func main() {
	client := rest.Client{URL: "http://www.example.com/"}
	api := rest.NewBaseAPI(http.MethodGet, "/", nil, new(string), nil)
	err := client.Do(api)
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
	resp := *api.ResponseObject().(*string)
	fmt.Println("Response:\n", resp)

}
