package main

import (
	"fmt"
	tba "github.com/go-telegram-bot-api/telegram-bot-api"
	"time"
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

	//Let's check if the antibot is enabled
	botCheckerEnabled, err := ctx.Database.GetSettingValue("botCheckerEnabled", ctx.Update.Message.Chat.ID)
	if err == nil && botCheckerEnabled == "y" {
		//Let's block the users
		falseBool := false
		for _, usr := range members {
			//If a user is already subscribed to the bot, he's probably not a userbot
			_, err := ctx.Database.GetUser(usr.ID)
			if err == nil {
				continue
			}

			// We need to be sure that the user is not already restricted
			chatMemberResult, err := ctx.Bot.GetChatMember(tba.ChatConfigWithUser{ChatID: ctx.Update.Message.Chat.ID, UserID: usr.ID})

			//The user can't send messages, so he's already limited.
			//We don't care about him
			//Can be "member" or "restricted" so far (need to check doc)
			if chatMemberResult.Status == "restricted" {
				continue
			}

			resp, err := ctx.Bot.RestrictChatMember(tba.RestrictChatMemberConfig{
				CanSendMessages: &falseBool,
				ChatMemberConfig: tba.ChatMemberConfig{
					ChatID: ctx.Update.Message.Chat.ID,
					UserID: usr.ID}})
			//All the new users are now blocked
			//We then need to send a message to each of them to let them click the button to verify they are not a bot
			//gotta make this stateless...
			//
			//IDEAL CAPTCHA
			//Thinking about a "easy" number test
			//Questions like "How may days has the month of feb in the year 2021?" A: 28
			//Questions like "How may legs had the Napoleon's 3-legs horse?" A: 3
			//Questions like "How many engineers does it take to change a light bulb?" - A: 2
			//Questions like "30 + 1"
			// <message>
			//[1] [2] [3] [4] [5]
			//[6] [7] [8] [9] [10]
			//[11] [12] [13] [14] [15]
			//[16] [17] [18] [19] [20]
			//[21] [22] [23] [24] [25]
			//[26] [27] [28] [29] [30]
			//[31] [32] [33] [34] [35]
			//[36] [37] [38] [39] [40]

			//This is a temporary solution; in future it could be a good idea to change it
			if err == nil && resp.Ok {
				//TODO: A good idea could be a countdown to delete the message after a while

				waitDuration, _ := time.ParseDuration("10s")

				message := tba.NewMessage(ctx.Update.Message.Chat.ID, ctx.Database.GetBotStringValueOrDefaultNoError("captchaMessageText", usr.LanguageCode))
				message.ReplyToMessageID = ctx.Update.Message.MessageID
				message.ReplyMarkup = tba.NewInlineKeyboardMarkup(
					tba.NewInlineKeyboardRow(
						tba.NewInlineKeyboardButtonData(
							fmt.Sprintf("%.0f:%.0f:%.0f", waitDuration.Hours(), waitDuration.Minutes(), waitDuration.Seconds()), "lolnothing-")))
				m, err := ctx.Bot.Send(message)
				if err == nil {

					unlockButtonTimer, timerIsStopped := time.NewTimer(waitDuration), false
					//If we had no issues sending the message, we start an async function
					go func() {
						<-unlockButtonTimer.C
						timerIsStopped = true
						time.Sleep(500 * time.Millisecond)
						ctx.Bot.Send(
							tba.NewEditMessageReplyMarkup(m.Chat.ID, m.MessageID,
								tba.NewInlineKeyboardMarkup(
									tba.NewInlineKeyboardRow(
										tba.NewInlineKeyboardButtonData(
											ctx.Database.GetBotStringValueOrDefaultNoError("captchaVerifyButtonText", usr.LanguageCode), "verify-")))))
					}()

					for !timerIsStopped {
						//timer is running
						time.Sleep(1 * time.Second)
						waitDuration = (time.Duration)(waitDuration.Nanoseconds() - 1*time.Second.Nanoseconds())
						ctx.Bot.Send(
							tba.NewEditMessageReplyMarkup(m.Chat.ID, m.MessageID,
								tba.NewInlineKeyboardMarkup(
									tba.NewInlineKeyboardRow(
										tba.NewInlineKeyboardButtonData(
											fmt.Sprintf("%.0f:%.0f:%.0f", waitDuration.Hours(), waitDuration.Minutes(), waitDuration.Seconds()), "lolnothing-")))))
					}

					//Here we're waiting 10 seconds to put the button in the message...
					// Hopefully, userbots aren't clever enough to consider this...
				}

				// I gotta be sure a user cannot reproduce the click in case is banned via telegram or another method...
				// In a stateless way...
				// TODO: check that the bat can actually see if the user joined the group (for the first time)
			}

		}

	}
}
