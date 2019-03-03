package main

import (
	tba "github.com/go-telegram-bot-api/telegram-bot-api"
)

func userJoinedRoute(ctx *Context) {
	members := *ctx.Update.Message.NewChatMembers
	locales := make([]string, 1)
	locales[0] = members[0].LanguageCode
	if len(members) > 1 {
		for _, usr := range members {
			localeIsPresent := false
			for _, loc := range locales {
				if loc == usr.LanguageCode {
					localeIsPresent = true
					break
				}
			}
			if !localeIsPresent {
				locales = append(locales, usr.LanguageCode)
			}
		}
	}
	for _, locale := range locales {
		if ctx.Database.StringExists("welcomeMessage", locale, ctx.Update.Message.Chat.ID) {
			welcomeMessageText, err := ctx.Database.GetStringValue("welcomeMessage", ctx.Update.Message.Chat.ID, locale)
			if err == nil {
				messageToSend := tba.NewMessage(ctx.Update.Message.Chat.ID, welcomeMessageText)
				if len(locales) == 1 {
					messageToSend.ReplyToMessageID = ctx.Update.Message.MessageID
				}
				messageToSend.DisableWebPagePreview = true
				messageToSend.ParseMode = tba.ModeHTML
				ctx.Bot.Send(messageToSend)
			}
		}
	}
}
