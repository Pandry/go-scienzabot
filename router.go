package main

import (
	"log"
	"strconv"
	"strings"

	tba "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (ctx *Context) route() {

	//Return message if there's no update
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

				case "/exists":
					msg := "You do "
					if !ctx.Database.UserExists(ctx.Update.Message.From.ID) {
						msg += "not "
					}
					msg += "exist."
					ctx.Bot.Send(tba.NewMessage(message.Chat.ID, msg))
					break

				case "/help":
				case "/aiuto":
				case "/aiutami":
					ctx.SendHelpMessage()
					break

				case "/info":
				case "/informazioni":
				case "/about":
				case "/github":

					break

				case "/version":
				case "/v":
					val, err := ctx.Database.GetBotSettingValue("test")
					if err != nil {

						msg := "error bla bla committing..."
						ctx.Bot.Send(tba.NewMessage(message.Chat.ID, msg))
						if ctx.Bot.Debug {
							log.Println("No read blabla...", err)
						}

						err := ctx.Database.SetBotSettingValue("test", strconv.Itoa(ctx.Update.Message.From.ID))
						if err != nil {
							msg = "dafuq, error..."
							if ctx.Bot.Debug {
								log.Println("Error doing commit...", err)
							}
						} else {
							msg = "done..."
						}
						ctx.Bot.Send(tba.NewMessage(message.Chat.ID, msg))
					} else {
						ctx.Bot.Send(tba.NewMessage(message.Chat.ID, val))
					}
					break

				case "/gdpr":

					break

				case "/registrazione":
				case "/registra":
				case "/registrami":
				case "/signup":

					break

				case "/iscrivi":
				case "/iscrivimi":
				case "/join":
				case "/iscrizione":
				case "/entra":
				case "/sottoscrivi":

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
			log.Println("User joined!\n Users' Nickname: ")
			for _, user := range *ctx.Update.Message.NewChatMembers {
				log.Println(user.UserName)
			}
			if ctx.Update.Message.GroupChatCreated {
				log.Println("This is a group just created! ")
			}

			//User joined
		} else if message.VideoNote != nil {
			//Video circolare

		} else if message.Video != nil {

		} else if message.Venue != nil {
			//NO IDEA     <--- non capisco (Bhez)
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
