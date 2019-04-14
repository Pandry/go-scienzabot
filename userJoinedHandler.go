package main

import (
	tba "github.com/go-telegram-bot-api/telegram-bot-api"
)

func userJoinedRoute(ctx *Context) {
	//The members who joined or were added to the groups
	members := *ctx.Update.Message.NewChatMembers
	//An array containing all the locales of the group who joined
	locales := make([]string, 1)
	//Since there is at least 1 guy, we take his locale and put it in the locales slice
	locales[0] = members[0].LanguageCode
	//If there are more than 1 person, iterate them to see if there is another locale to consider
	if len(members) > 1 {
		for _, usr := range members {
			localeIsPresent := false
			for _, loc := range locales {
				if loc == usr.LanguageCode {
					localeIsPresent = true
					break
				}
			}
			//If the locale of the current user is not present, add it to the locales slice
			if !localeIsPresent {
				locales = append(locales, usr.LanguageCode)
			}
		}
	}
	//For each locale, welcome the user with the correct locale
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
