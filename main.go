package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/kelseyhightower/envconfig"
)

type specification struct {
	WrikeBearer string
}

var s specification

func main() {
	setup()
	exampleRequest()
}

func exampleRequest() {
	client := &http.Client{}

	req, err := http.NewRequest("GET", "https://www.wrike.com/api/v3/accounts", nil)
	req.Header.Add("Authorization", "bearer "+s.WrikeBearer)

	resp, err := client.Do(req)
	defer resp.Body.Close()

	if err != nil {
		log.Println("Client Error: ", err)
		return
	}
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading response body: ", err)
	}
	log.Println(string(contents))
}

func setup() {
	err := envconfig.Process("wrikeclientpoc", &s)
	if err != nil {
		log.Fatal("boo: ", err.Error())
	}
}
