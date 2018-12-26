package main


import "github.com/go-telegram-bot-api/telegram-bot-api"

// Context is a struct that contains the current redis client, the bot and a update.
// It is used as a paramater, passed withafter the routing
type Context struct {
	Bot    *tgbotapi.BotAPI
	Update *tgbotapi.Update
}

// UpdateInfo is a struct that contains the current redis client, the bot and a update.
// It is used as a paramater, passed withafter the routing
/*
type UpdateInfo struct {
	Forworded bool
}
*/
