package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
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
	f, err := os.ReadFile("/usr/local/alert_to_line/config.json")
	if err != nil {
		log.Fatal(err)
	}

	var cfg Config

	json.Unmarshal(f, &cfg)
	return &cfg, err
}

func sendMessage(message string) {
	env, err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}

	requestBody := RequestBody{
		To: env.UserID,
		Messages: []Message{
			{
				Type: "text",
				Text: message,
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

	byteArray, err := io.ReadAll(resp.Body)
	if err != nil {
		panic("Error")
	}

	fmt.Printf("%#v\n", string(byteArray))

}

func main() {
	flag.Parse()
	args := flag.Args()

	if args[0] == "test" {
		sendMessage("Test Message")
		return
	}

	alertType := args[0]
	var alertMsg string
	var stateIcon string

	if alertType == "SERVICE" {
		switch args[5] {
		case "OK":
			stateIcon = "✅"
		case "WARNING":
			stateIcon = "⚠️"
		case "CRITICAL":
			stateIcon = "🚫"
		case "UNKNOWN":
			stateIcon = "❓"
		default:
			stateIcon = "📛"
		}
	} else if alertType == "HOST" {
		switch args[3] {
		case "UP":
			stateIcon = "✅"
		case "DOWN":
			stateIcon = "🚫"
		case "UNREACHABLE":
			stateIcon = "❓"
		default:
			stateIcon = "📛"
		}
	} else {
		log.Fatal("first arg is not HOST or SERVICE")
	}

	switch alertType {
	case "HOST":
		alertMsg = fmt.Sprintf("%v %v %v\nHost: %v\n\n %v",args[0], args[3], stateIcon, args[2], args[5] )
	case "SERVICE":
		alertMsg = fmt.Sprintf("%v %v %v %v\nHost: %v\n\n%v",args[0], args[2], args[5], stateIcon, args[3], args[7])
	default:
		log.Fatal("first arg is not HOST or SERVICE")
	}

	sendMessage(alertMsg)
}
