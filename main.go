package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/kelseyhightower/envconfig"
)

type refreshResponse struct {
	refreshToken string
}

type specification struct {
	WrikeBearer       string
	WrikeClientID     string
	WrikeClientSecret string
	WrikeRefreshToken string // temp for testing: should get this from response later
}

var (
	s specification
)

func main() {
	setup()

	refreshAuthToken()

	exampleRequest()
}

func exampleRequest() {
	log.Println("Making a request")
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

func refreshAuthToken() {
	log.Println("Refreshing token...")
	var refreshJSON refreshResponse
	client := &http.Client{}

	body := makeRefreshBody()

	req, err := http.NewRequest("POST", "https://www.wrike.com/oauth2/token", strings.NewReader(body))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded") // this is what curl -d sends by default

	resp, err := client.Do(req)
	defer resp.Body.Close()

	if err != nil {
		log.Println("Client Error on refresh: ", err)
		return
	}

	if resp.StatusCode != 200 {
		log.Fatal("Got response code " + resp.Status)
	}

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading response body during refresh of token: ", err)
	}

	// needs to reassign refresh token
	// refreshToken = something from body
	log.Println(string(contents))

	err = json.Unmarshal([]byte(contents), &refreshJSON)

	if err != nil {
		log.Fatal("Error unmarshalling: ", err)
	}

	log.Println("json parsed[refreshToken] is: " + refreshJSON.refreshToken)
}

func makeRefreshBody() string {
	body := "client_id=" + s.WrikeClientID + "&client_secret=" + s.WrikeClientSecret
	body += "&grant_type=refresh_token&refresh_token=" + s.WrikeRefreshToken

	log.Println("Body is " + body)

	return body
}

func setup() {
	err := envconfig.Process("wrikeclientpoc", &s)
	if err != nil {
		log.Fatal("boo: ", err.Error())
	}

	log.Println("bearer is " + s.WrikeBearer)
	log.Println("client is " + s.WrikeClientID)
	log.Println("secret is " + s.WrikeClientSecret)
	log.Println("refresh is " + s.WrikeRefreshToken)
}
