package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Message struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type RequestBody struct {
	To       string    `json:"to"`
	Messages []Message `json:"messages"`
}

type Config struct {
	UserID             string `json:"USER_ID"`
	ChannelAccessToken string `json:"CHANNEL_ACCESS_TOKEN"`
}

const (
	ENDPOINT = "https://api.line.me/v2/bot/message/push"
)

func loadConfig() (*Config, error) {
	f, err := ioutil.ReadFile("./config.json")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var cfg Config

	json.Unmarshal(f, &cfg)
	return &cfg, err
}

func main() {
	env, err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}

	requestBody := RequestBody{
		To: env.UserID,
		Messages: []Message{
			{
				Type: "text",
				Text: "uooooo",
			},
		},
	}

	jsonString, err := json.Marshal(requestBody)
	if err != nil {
		panic("Error")
	}
	req, err := http.NewRequest("POST", ENDPOINT, bytes.NewBuffer(jsonString))
	if err != nil {
		panic("Error")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+env.ChannelAccessToken)

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		panic("Error")
	}
	defer resp.Body.Close()

	byteArray, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic("Error")
	}

	fmt.Printf("%#v\n", string(byteArray))
}
