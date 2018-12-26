package main

import (
	"log"
	"os"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_TOEKN"))
	if err != nil {
		log.Panic(err)
	}

	ctx := Context{bot, nil}

	bot.Debug = true
	log.Printf("Authorized on account @%s", bot.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)
	for update := range updates {
		ctx.Update = &update
		route(&ctx)
	}
}
