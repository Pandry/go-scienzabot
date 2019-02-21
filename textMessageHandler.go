package main

import (
	"strings"

	tba "github.com/go-telegram-bot-api/telegram-bot-api"
)

func textMessageRoute(ctx *Context) {
	message := ctx.Update.Message
	var (
		err         error
		messageBody string
	)

	userInDB := ctx.Database.UserExists(message.From.ID)
	messageInGroup := message.Chat.IsGroup() || message.Chat.IsSuperGroup()

	if userInDB {
		ctx.Database.UpdateUserLastSeen(message.From.ID, message.Time())
	}

	if message.IsCommand() {
		//Command

		switch args := strings.Split(message.Text, " "); args[0] {
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

		case "/help", "/aiuto", "/aiutami":
			if messageInGroup {
				messageBody = "onPrivateChatCommand"
			} else {
				messageBody = "helpCommand"
			}
			if messageBody, err = ctx.Database.GetBotStringValueOrDefault(messageBody, message.From.LanguageCode); err == nil {
				a := tba.NewInlineKeyboardMarkup(
					tba.NewInlineKeyboardRow(
						tba.NewInlineKeyboardButtonData(" ", "delme")))
				message := tba.NewMessage(message.Chat.ID, messageBody)
				message.ReplyMarkup = a
				ctx.Bot.Send(message)
			}
			break

		case "/info", "/informazioni", "/about", "/github":
			if messageInGroup {
				messageBody = "onPrivateChatCommand"
			} else {
				messageBody = "infoCommand"
			}

			if messageBody, err = ctx.Database.GetBotStringValueOrDefault(messageBody, message.From.LanguageCode); err == nil {
				ctx.Bot.Send(tba.NewMessage(message.Chat.ID, messageBody))
			}
			break

		case "/version", "/v":

			if val, err := ctx.Database.GetBotSettingValue("version"); err != nil {
				ctx.Bot.Send(tba.NewMessage(message.Chat.ID, val))
			}
			break

		case "/gdpr":

			break

		case "/registrazione", "/registra", "/registrami", "/signup":

			break

		case "/iscrivi", "/iscrivimi", "/join", "/iscrizione", "/entra", "/sottoscrivi":

			break

		case "/segnalibro", "/salva":

			break

		default:
			//Check if it exists in DB

		}
	}
}
