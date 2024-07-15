package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
	"github.com/sumit23802380/slack-gobot/slackbot/config"
	"github.com/sumit23802380/slack-gobot/slackbot/handlers"
)

func main() {
	godotenv.Load(".env")
	config.SetupConfig()
	// db.Init()
	// defer db.Close()
	
	token := os.Getenv("SLACK_OAUTH_TOKEN")
	appToken := os.Getenv("SLACK_APP_TOKEN")
	client := slack.New(token, slack.OptionDebug(true), slack.OptionAppLevelToken(appToken))

	socketClient := socketmode.New(
		client,
		socketmode.OptionDebug(true),
		socketmode.OptionLog(log.New(os.Stdout, "socketmode: ", log.Lshortfile|log.LstdFlags)),
	)

	_, cancel := context.WithCancel(context.Background())
	defer cancel()
	go handlers.SlackHandler(client, socketClient)

	socketClient.Run()
}
