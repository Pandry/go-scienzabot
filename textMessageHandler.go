package main

import (
	"scienzabot/consts"
	"scienzabot/database"
	"scienzabot/utils"
	"strings"

	tba "github.com/go-telegram-bot-api/telegram-bot-api"
)

func textMessageRoute(ctx *Context) {
	message := ctx.Update.Message

	//For the moment we don't care about channels
	if message.Chat.IsChannel() {
		return
	}

	var (
		err         error
		messageBody string
		user        database.User
		isAdmin     bool
	)

	if userExists := ctx.Database.UserExists(message.From.ID); userExists {
		user, err = ctx.Database.GetUser(message.From.ID)
		isAdmin = utils.HasPermission(int(user.Permissions), consts.UserPermissionAdmin)
	}

	userInDB := ctx.Database.UserExists(message.From.ID)
	messageInGroup := message.Chat.IsGroup() || message.Chat.IsSuperGroup()
	if messageInGroup {
		if !ctx.Database.GroupExists(message.Chat.ID) {
			ref := message.Chat.InviteLink
			if ref == "" {
				if message.Chat.UserName != "" {
					ref = "https://t.me/" + message.Chat.UserName
				}
			}
			ctx.Database.AddGroup(database.Group{ID: message.Chat.ID, Title: message.Chat.Title, Ref: message.Chat.UserName})
		}
	}

	if userInDB {
		ctx.Database.UpdateUserLastSeen(message.From.ID, message.Time())
		if messageInGroup {
			ctx.Database.IncrementMessageCount(int64(message.From.ID), message.Chat.ID)
		}
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
				rm := tba.NewInlineKeyboardMarkup(
					tba.NewInlineKeyboardRow(
						tba.NewInlineKeyboardButtonData(" ", "delme-")))
				message := tba.NewMessage(message.Chat.ID, messageBody)
				message.ReplyMarkup = rm
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

		case "/del", "/deleteMessage":
			if isAdmin {
				ctx.Bot.DeleteMessage(tba.NewDeleteMessage(message.Chat.ID, message.ReplyToMessage.MessageID))
			}
			break

		case "/registrazione", "/registra", "/registrati", "/registrami", "/signup":
			if !userInDB {
				//We want registration to happen in private, not in public
				if messageInGroup {
					if messageBody, err = ctx.Database.GetBotStringValueOrDefault("onPrivateChatCommand", message.From.LanguageCode); err == nil {
						messageToSend := tba.NewMessage(message.Chat.ID, messageBody)
						rm := tba.NewInlineKeyboardMarkup(
							tba.NewInlineKeyboardRow(
								tba.NewInlineKeyboardButtonData(ctx.Database.GetBotStringValueOrDefaultNoError("deleteMessageText", message.From.LanguageCode), "delme-")))
						//tba.NewInlineKeyboardButtonData(" ", "delme-")))
						messageToSend.ReplyMarkup = rm
						messageToSend.ReplyToMessageID = message.MessageID
						ctx.Bot.Send(messageToSend)
					}
					//Warn about error?
					return
				}
				err = ctx.Database.AddUser(database.User{ID: int64(message.From.ID), Nickname: message.From.UserName, Status: consts.UserStatusActive})
				if err != nil {
					messageBody, err = ctx.Database.GetBotStringValueOrDefault("generalError", message.From.LanguageCode)
					messageToSend := tba.NewMessage(message.Chat.ID, messageBody)
					rm := tba.NewInlineKeyboardMarkup(
						tba.NewInlineKeyboardRow(
							tba.NewInlineKeyboardButtonData(
								ctx.Database.GetBotStringValueOrDefaultNoError("deleteMessageText", message.From.LanguageCode), "delme-")))
					messageToSend.ReplyMarkup = rm
					messageToSend.ReplyToMessageID = message.MessageID
					ctx.Bot.Send(messageToSend)
				} else {
					messageBody, _ = ctx.Database.GetBotStringValueOrDefault(
						ctx.Database.GetBotStringValueOrDefaultNoError("deleteMessageText", message.From.LanguageCode), message.From.LanguageCode)
					messageToSend := tba.NewMessage(message.Chat.ID, messageBody)
					rm := tba.NewInlineKeyboardMarkup(
						tba.NewInlineKeyboardRow(
							tba.NewInlineKeyboardButtonData(
								ctx.Database.GetBotStringValueOrDefaultNoError("deleteMessageText", message.From.LanguageCode), "delme-")))
					messageToSend.ReplyMarkup = rm
					messageToSend.ReplyToMessageID = message.MessageID
					ctx.Bot.Send(messageToSend)
				}
			} else {
				messageBody, _ = ctx.Database.GetBotStringValueOrDefault("userAlreadyRegistred", message.From.LanguageCode)
				messageToSend := tba.NewMessage(message.Chat.ID, messageBody)
				rm := tba.NewInlineKeyboardMarkup(
					tba.NewInlineKeyboardRow(
						tba.NewInlineKeyboardButtonData(
							ctx.Database.GetBotStringValueOrDefaultNoError("deleteMessageText", message.From.LanguageCode), "delme-")))
				messageToSend.ReplyMarkup = rm
				messageToSend.ReplyToMessageID = message.MessageID
				ctx.Bot.Send(messageToSend)
			}
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
