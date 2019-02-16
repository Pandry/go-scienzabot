package main

import (
	"log"
	"strconv"
	"strings"

	tba "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (ctx *Context) route() {

	switch {
	//Return message if there's no update
	case ctx.Update == nil:
		return
	case ctx.Update.Message != nil:
		//General message

		message := ctx.Update.Message
		userIsInGroup := message.Chat.IsSuperGroup() || message.Chat.IsGroup()
		_, err := ctx.Database.GetUser(int64(message.From.ID))
		userIsRegistred := err == nil

		//Count last time the user was seen in a group
		if userIsInGroup && userIsRegistred {
			ctx.Database.UpdateUserLastSeen(message.From.ID, message.Time())
		}

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

		} else {
			switch {

			case message.Sticker != nil:
				//Sticker

				break

			case message.Photo != nil:
				//Photo

				break

			case message.Voice != nil:
				//Voice audio file

				break

			case message.Document != nil:
				//Document

				break

			case message.Location != nil:

				break

			case message.Contact != nil:

				break

			case message.NewChatMembers != nil:
				//New user(s)
				log.Println("User joined!\n Users' Nickname: ")
				for _, user := range *ctx.Update.Message.NewChatMembers {
					log.Println(user.UserName)
				}
				if ctx.Update.Message.GroupChatCreated {
					log.Println("This is a group just created! ")
				}

				//User joined
				break

			case message.VideoNote != nil:
				//Video circolare

				break

			case message.Video != nil:

				break

			case message.Venue != nil:
				//NO IDEA     <--- non capisco (Bhez)
				break

			case message.LeftChatMember != nil:
				//User removed (could be the bot)
				break

			case message.PinnedMessage != nil:

				break

			case message.NewChatPhoto != nil:

				break

			case message.NewChatTitle != "":

				break

			case message.MigrateToChatID != 0:

				break
			}
		}

	case ctx.Update.EditedMessage != nil:
		//Edited text message
		break

	case ctx.Update.CallbackQuery != nil:
		//Callback query
		break

	case ctx.Update.InlineQuery != nil:
		//Inline query
		return

	case ctx.Update.ChannelPost != nil:
		//Channel post
		return

	case ctx.Update.EditedChannelPost != nil:
		//Edited channel post
		return

	case ctx.Update.PreCheckoutQuery != nil:
		//Pre checkoput query - useless for now
		return

	case ctx.Update.ShippingQuery != nil:
		//Pre shipping query
		return

	case ctx.Update.ChosenInlineResult != nil:
		//Chosen inline result -> Chosen inline element?
		return

	}
}
