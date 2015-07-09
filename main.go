package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/kelseyhightower/envconfig"
)

type refreshResponse struct {
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
}

type wrikeTaskResponse struct {
	Tasks []wrikeTask `json:"data"`
}

type wrikeTask struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Status    string `json:"status"`
	Permalink string `json:"permalink"`
}

type specification struct {
	WrikeBearer       string
	WrikeClientID     string
	WrikeClientSecret string
	WrikeRefreshToken string // temp for testing: should get this from response later
	SlackURL          string
}

var (
	s          specification
	wrikeTasks wrikeTaskResponse
)

func main() {
	setup()

	// sendToSlack("check check mic check")

	refreshAuthToken()

	getRecentTasks()

	if 1 == 0 { // hilariously bad debug hack
		exampleRequest()
	}
}

func sendToSlack(slackMessage string) {
	data := url.Values{}
	payloadBody := "{\"channel\": \"#dev-null\", \"username\": \"hook it up\", \"text\": \"" + slackMessage + "\", \"icon_emoji\": \":ghost:\"}"
	data.Set("payload", payloadBody)

	log.Println("sending body: " + payloadBody)
	log.Println("sending encoded body: " + data.Encode())

	client := &http.Client{}
	req, err := http.NewRequest("POST", s.SlackURL, strings.NewReader(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded") // this is what curl -d sends by default

	resp, err := client.Do(req)
	defer resp.Body.Close()

	if err != nil {
		log.Println("Client Error on Slackchat: ", err)
		return
	}

	if resp.StatusCode != 200 {
		log.Fatal("Got response code " + resp.Status)
	}
}

func getRecentTasks() {
	log.Println("Getting recent tasks")
	client := &http.Client{}

	var url = "https://www.wrike.com/api/v3/tasks"
	url += "?createdDate={\"start\":\"2015-07-09T21:28:00Z\"}"

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", "bearer "+s.WrikeBearer)

	resp, err := client.Do(req)
	defer resp.Body.Close()

	if err != nil {
		log.Println("Client Error: ", err)
		return
	}

	if resp.StatusCode != 200 {
		log.Println("Didn't get a 200 on tasks requests, getting new token")
		// assume invalid auth token and refresh, I suppose
		// refreshAuthToken()
		// try again
	}

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading response body: ", err)
	}

	err = json.Unmarshal([]byte(contents), &wrikeTasks)

	if err != nil {
		log.Fatal("Error unmarshalling wrike tasks: ", err)
	}

	log.Println("hey found json! first title is ", wrikeTasks.Tasks[0].Title)

	log.Println("sending to slack: " + wrikeTasks.Tasks[0].Title)
	sendToSlack(wrikeTasks.Tasks[0].Title + " at " + wrikeTasks.Tasks[0].Permalink)
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

	err = json.Unmarshal([]byte(contents), &refreshJSON)

	if err != nil {
		log.Fatal("Error unmarshalling: ", err)
	}

	s.WrikeRefreshToken = refreshJSON.RefreshToken
	s.WrikeBearer = refreshJSON.AccessToken
}

func makeRefreshBody() string {
	body := "client_id=" + s.WrikeClientID + "&client_secret=" + s.WrikeClientSecret
	body += "&grant_type=refresh_token&refresh_token=" + s.WrikeRefreshToken

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
	log.Println("slack url is " + s.SlackURL)
}
