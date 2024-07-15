package handlers

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
	"github.com/sumit23802380/slack-gobot/slackbot/services"
)

func SlackHandler(client *slack.Client, socketClient *socketmode.Client) {
	for {
		select {
		case event := <-socketClient.Events:
			switch event.Type {
			case socketmode.EventTypeEventsAPI:
				eventsAPI, ok := event.Data.(slackevents.EventsAPIEvent)
				if !ok {
					log.Printf("Could not type cast the event to the EventsAPI: %v\n", event)
					continue
				}

				socketClient.Ack(*event.Request)

				err := HandleEventMessage(eventsAPI, client)
				if err != nil {
					// Replace with actual err handeling
					log.Fatal(err)
				}
				if eventsAPI.Type == slackevents.CallbackEvent {
					innerEvent := eventsAPI.InnerEvent
					switch ev := innerEvent.Data.(type) {
					case *slackevents.MessageEvent:
						text := ev.Text
						channelID := ev.Channel
						log.Print(text)
						// Save message to the database
						if err := services.SaveMessage(text); err != nil {
							log.Printf("Failed to save message: %v\n", err)
							continue
						}

						// Reply with the same message
						_, _, err := client.PostMessage(channelID, slack.MsgOptionText(text, false))
						if err != nil {
							log.Printf("Failed to post message: %v\n", err)
						}
					}
				}
			}
		}
	}
}

func HandleAppMentionEventToBot(event *slackevents.AppMentionEvent, client *slack.Client) error {

	user, err := client.GetUserInfo(event.User)
	if err != nil {
		return err
	}

	text := strings.ToLower(event.Text)

	attachment := slack.Attachment{}

	if strings.Contains(text, "hello") || strings.Contains(text, "hi") {
		attachment.Text = fmt.Sprintf("Hello %s", user.Name)
		attachment.Color = "#4af030"
	} else if strings.Contains(text, "weather") {
		attachment.Text = fmt.Sprintf("Weather is sunny today. %s", user.Name)
		attachment.Color = "#4af030"
	} else {
		attachment.Text = fmt.Sprintf("I am good. How are you %s?", user.Name)
		attachment.Color = "#4af030"
	}
	_, _, err = client.PostMessage(event.Channel, slack.MsgOptionAttachments(attachment))
	if err != nil {
		return fmt.Errorf("failed to post message: %w", err)
	}
	return nil
}

func HandleEventMessage(event slackevents.EventsAPIEvent, client *slack.Client) error {
	switch event.Type {

	case slackevents.CallbackEvent:

		innerEvent := event.InnerEvent

		switch evnt := innerEvent.Data.(type) {
		case *slackevents.AppMentionEvent:
			err := HandleAppMentionEventToBot(evnt, client)
			if err != nil {
				return err
			}
		}
	default:
		return errors.New("unsupported event type")
	}
	return nil
}
