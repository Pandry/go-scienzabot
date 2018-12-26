package main

import (
	"log"
	"strings"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func route(ctx *Context) {

	//Return messager if there's no update
	if ctx.Update == nil {
		return
	}

	if ctx.Update.Message != nil {
		//General message

		message := ctx.Update.Message

		if message.Text != "" {
			//Text message
			if message.IsCommand() {
				//Command
				args := strings.Split(message.Text, " ")
				switch args[0] {
				case "/start":
					break

				case "/help":
					msg := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, "Help text!")
					msg.ReplyToMessageID = ctx.Update.Message.MessageID

					ctx.Bot.Send(msg)

					break
				}

			}

		} else if message.Sticker != nil {
			//Sticker

		} else if message.Photo != nil {
			//Photo

		} else if message.Voice != nil {
			//Voice audio file

		} else if message.Document != nil {
			//Document

		} else if message.Location != nil {

		} else if message.Contact != nil {

		} else if message.NewChatMembers != nil {
			//New user(s)
			newGroup := false
			log.Println("User joined!\n Users' Nickname: ")
			for _, user := range *ctx.Update.Message.NewChatMembers {
				log.Println(user.UserName)
			}
			if ctx.Update.Message.GroupChatCreated {
				log.Println("This is a group just created! ")
			}

			UserJoined(ctx)
		} else if message.VideoNote != nil {
			//Video circoare

		} else if message.Video != nil {

		} else if message.Venue != nil {
			//NO IDEA
		} else if message.LeftChatMember != nil {
			//User removed (could be the bot)
		} else if message.PinnedMessage != nil {

		} else if message.NewChatPhoto != nil {

		} else if message.NewChatTitle != "" {

		} else if message.MigrateToChatID != 0 {

		}

	} else if ctx.Update.EditedMessage != nil {
		//Edited text message
	} else if ctx.Update.CallbackQuery != nil {
		//Callback query
	} else if ctx.Update.InlineQuery != nil {
		//Inline query
		return
	} else if ctx.Update.ChannelPost != nil {
		//Channel post
		return
	} else if ctx.Update.EditedChannelPost != nil {
		//Edited channel post
		return
	} else if ctx.Update.PreCheckoutQuery != nil {
		//Pre checkoput query - useless for now
		return
	} else if ctx.Update.ShippingQuery != nil {
		//Pre shipping query
		return
	} else if ctx.Update.ChosenInlineResult != nil {
		//Chosen inline result -> Chosen inline element?
		return
	}
	return
}
