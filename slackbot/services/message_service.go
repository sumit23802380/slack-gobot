package services

import (
	"github.com/sumit23802380/slack-gobot/slackbot/db"
	"github.com/sumit23802380/slack-gobot/slackbot/models"
)

func SaveMessage(text string) error {
	message := models.Message{
		Text: text,
	}
	return db.DB.Create(&message).Error
}
