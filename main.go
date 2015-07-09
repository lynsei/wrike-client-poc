package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/kelseyhightower/envconfig"
)

type specification struct {
	WrikeBearer       string
	WrikeClientID     string
	WrikeClientSecret string
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

func refreshToken() {
	log.Println("Refreshing token...")
	client := &http.Client{}

	body := makeRefreshBody()

	req, err := http.NewRequest("POST", "https://www.wrike.com/oauth2/token", strings.NewReader(body))

	resp, err := client.Do(req)
	defer resp.Body.Close()

	if err != nil {
		log.Println("Client Error on refresh: ", err)
		return
	}
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading response body during refresh of token: ", err)
	}
	log.Println(string(contents))
}

func makeRefreshBody() string {
	return "body"
}

func setup() {
	err := envconfig.Process("wrikeclientpoc", &s)
	if err != nil {
		log.Fatal("boo: ", err.Error())
	}

	log.Println("bearer is " + s.WrikeBearer)
	log.Println("client is " + s.WrikeClientID)
	log.Println("secret is " + s.WrikeClientSecret)
}
