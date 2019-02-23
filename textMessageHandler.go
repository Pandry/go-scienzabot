package main

import (
	"regexp"
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
		err              error
		messageBody      string
		user             database.User
		userPermission   int
		userIsBotAdmin   bool
		userIsGroupAdmin bool
		groupStatus      int64
	)

	if userExists := ctx.Database.UserExists(message.From.ID); userExists {
		user, err = ctx.Database.GetUser(message.From.ID)
		userIsBotAdmin = utils.HasPermission(int(user.Permissions), consts.UserPermissionAdmin)
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
				if message.Chat.InviteLink != "" {
					ref = message.Chat.InviteLink
				}
			}
			ctx.Database.AddGroup(database.Group{ID: message.Chat.ID, Title: message.Chat.Title, Ref: message.Chat.UserName})
			reloadChatAdmins(ctx)

		} else {
			groupStatus, _ = ctx.Database.GetGroupStatus(message.Chat.ID)
			if !userIsBotAdmin && utils.HasPermission(int(groupStatus), consts.GroupBanned) {
				return
			}
			userPermission, err = ctx.Database.GetPermission(int64(message.From.ID), message.Chat.ID)
			userIsGroupAdmin = utils.HasPermission(userPermission, consts.UserPermissionGroupAdmin) || utils.HasPermission(userPermission, consts.UserPermissionAdmin)
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
			messageToSend := tba.NewMessage(message.Chat.ID, msg)
			rm := tba.NewInlineKeyboardMarkup(
				tba.NewInlineKeyboardRow(
					tba.NewInlineKeyboardButtonData(
						ctx.Database.GetBotStringValueOrDefaultNoError("deleteMessageText", message.From.LanguageCode), "delme-")))
			messageToSend.ReplyMarkup = rm
			messageToSend.ReplyToMessageID = message.MessageID
			ctx.Bot.Send(messageToSend)
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
			if val, err := ctx.Database.GetBotSettingValue("version"); err == nil {
				messageToSend := tba.NewMessage(message.Chat.ID, val)
				rm := tba.NewInlineKeyboardMarkup(
					tba.NewInlineKeyboardRow(
						tba.NewInlineKeyboardButtonData(ctx.Database.GetBotStringValueOrDefaultNoError("deleteMessageText", message.From.LanguageCode), "delme-")))
				//tba.NewInlineKeyboardButtonData(" ", "delme-")))
				messageToSend.ReplyMarkup = rm
				messageToSend.ReplyToMessageID = message.MessageID
				ctx.Bot.Send(messageToSend)
			}
			break

		case "/ping":
			if userIsBotAdmin {
				messageToSend := tba.NewMessage(message.Chat.ID, "🏓 Pong!")
				rm := tba.NewInlineKeyboardMarkup(
					tba.NewInlineKeyboardRow(
						tba.NewInlineKeyboardButtonData(ctx.Database.GetBotStringValueOrDefaultNoError("deleteMessageText", message.From.LanguageCode), "delme-")))
				//tba.NewInlineKeyboardButtonData(" ", "delme-")))
				messageToSend.ReplyMarkup = rm
				messageToSend.ReplyToMessageID = message.MessageID
				ctx.Bot.Send(messageToSend)
			}
			break

		case "/gdpr":

			break

		case "/lists":
			if userInDB {
				grps, _ := ctx.Database.GetUserGroups(message.From.ID)
				messageBody := ""
				for _, group := range grps {
					messageBody += group.Title + "\n"
					lists, _ := ctx.Database.GetLists(group.ID)
					for _, lst := range lists {
						messageBody += "  " + lst.Name + "\n"

					}
				}
				messageToSend := tba.NewMessage(message.Chat.ID, messageBody)
				rm := tba.NewInlineKeyboardMarkup(
					tba.NewInlineKeyboardRow(
						tba.NewInlineKeyboardButtonData(ctx.Database.GetBotStringValueOrDefaultNoError("deleteMessageText", message.From.LanguageCode), "delme-")))
				//tba.NewInlineKeyboardButtonData(" ", "delme-")))
				messageToSend.ReplyMarkup = rm
				messageToSend.ReplyToMessageID = message.MessageID
				ctx.Bot.Send(messageToSend)
			}
			break

		case "/newlist":

			if len(args) != 2 {
				messageBody, err = ctx.Database.GetBotStringValueOrDefault("newlistSyntaxError", message.From.LanguageCode)
				messageToSend := tba.NewMessage(message.Chat.ID, messageBody)
				rm := tba.NewInlineKeyboardMarkup(
					tba.NewInlineKeyboardRow(
						tba.NewInlineKeyboardButtonData(
							ctx.Database.GetBotStringValueOrDefaultNoError("deleteMessageText", message.From.LanguageCode), "delme-")))
				messageToSend.ReplyMarkup = rm
				messageToSend.ReplyToMessageID = message.MessageID
				ctx.Bot.Send(messageToSend)
				return
			}
			listNameIsValid, _ := regexp.MatchString("^[a-z\\-_]{1,30}$", args[1])
			if !listNameIsValid {
				messageBody, err = ctx.Database.GetBotStringValueOrDefault("newlistSyntaxError", message.From.LanguageCode)
				messageToSend := tba.NewMessage(message.Chat.ID, messageBody)
				rm := tba.NewInlineKeyboardMarkup(
					tba.NewInlineKeyboardRow(
						tba.NewInlineKeyboardButtonData(
							ctx.Database.GetBotStringValueOrDefaultNoError("deleteMessageText", message.From.LanguageCode), "delme-")))
				messageToSend.ReplyMarkup = rm
				messageToSend.ReplyToMessageID = message.MessageID
				ctx.Bot.Send(messageToSend)
				return
			}

			if messageInGroup {

				if userIsBotAdmin || userIsGroupAdmin || utils.HasPermission(userPermission, consts.UserPermissionCanCreateList) {

					err = ctx.Database.AddList(database.List{Name: args[1], GroupID: message.Chat.ID})
					if err == nil {
						messageToSend := tba.NewMessage(message.Chat.ID, ctx.Database.GetBotStringValueOrDefaultNoError("listCreatedSuccessfully", message.From.LanguageCode))
						rm := tba.NewInlineKeyboardMarkup(
							tba.NewInlineKeyboardRow(
								tba.NewInlineKeyboardButtonData(ctx.Database.GetBotStringValueOrDefaultNoError("deleteMessageText", message.From.LanguageCode), "delme-")))
						//tba.NewInlineKeyboardButtonData(" ", "delme-")))
						messageToSend.ReplyMarkup = rm
						messageToSend.ReplyToMessageID = message.MessageID
						ctx.Bot.Send(messageToSend)
					}
				} else {
					messageBody, err = ctx.Database.GetBotStringValueOrDefault("notAuthorized", message.From.LanguageCode)
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
				//TODO: implement group choosing where is admin
				messageBody, err = ctx.Database.GetBotStringValueOrDefault("generalError", message.From.LanguageCode)
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

		case "/del", "/deleteMessage":
			if message.ReplyToMessage != nil && (userIsBotAdmin || userIsGroupAdmin) {
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
						ctx.Database.GetBotStringValueOrDefaultNoError("userAddedSuccessfully", message.From.LanguageCode), message.From.LanguageCode)
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
				messageBody, _ = ctx.Database.GetBotStringValueOrDefault("userNotRegistred", message.From.LanguageCode)
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

		case "/segnalibro", "/salva":

			break

		case "/reloadpermissions", "/ricarica", "/riavvia", "/restart":
			reloadChatAdmins(ctx)
			break

		default:
			//Check if it exists in DB

		}
	}
}

func reloadChatAdmins(ctx *Context) {
	admins, err := ctx.Bot.GetChatAdministrators(ctx.Update.Message.Chat.ChatConfig())

	if err != nil {
		return
	}

	currentAdms, err := ctx.Database.GetPrivilegedUsers(ctx.Update.Message.Chat.ID)

	if err != nil {
		return
	}

	//List of admins in the database
	privUsers := make([]database.Permission, 0)

	//First, we add to the db all the users that are not admin
	//  via telegram admin permission
	for _, usr := range currentAdms {
		//If the user has permissions other then the admin value, add it to the new list without the group admin permissions
		if utils.RemovePermission(int(usr.Permission), consts.UserPermissionGroupAdmin) != 0 {
			privUsers = append(privUsers, database.Permission{
				UserID:     usr.UserID,
				GroupID:    usr.GroupID,
				Permission: int64(utils.RemovePermission(int(usr.Permission), consts.UserPermissionGroupAdmin))})
		}
	}

	//Then we add all the admins who are admin of the group
	for _, tAdm := range admins {
		found := false

		//First we see if the user is already in the new list
		for i, nAdm := range privUsers {
			//If so, we just add the group admin permission
			if tAdm.User.ID == int(nAdm.UserID) {
				privUsers[i].Permission = int64(utils.SetPermission(int(nAdm.Permission), consts.UserPermissionGroupAdmin))
				found = true
				break
			}
		}

		if found {
			continue
		}

		privUsers = append(privUsers, database.Permission{
			UserID:     int64(tAdm.User.ID),
			GroupID:    ctx.Update.Message.Chat.ID,
			Permission: consts.UserPermissionGroupAdmin})
	}

	//We then remove all the permissions from the group
	ctx.Database.RemoveAllGroupPermissions(ctx.Update.Message.Chat.ID)
	//And readd them
	for _, p := range privUsers {
		ctx.Database.SetPermissions(p)
	}

}

func replyMessageWithCloseButton(ctx *Context, keyString string) {
	messageBody, _ := ctx.Database.GetBotStringValueOrDefault(
		ctx.Database.GetBotStringValueOrDefaultNoError(keyString, ctx.Update.Message.From.LanguageCode), ctx.Update.Message.From.LanguageCode)
	messageToSend := tba.NewMessage(ctx.Update.Message.Chat.ID, messageBody)
	rm := tba.NewInlineKeyboardMarkup(
		tba.NewInlineKeyboardRow(
			tba.NewInlineKeyboardButtonData(
				ctx.Database.GetBotStringValueOrDefaultNoError("deleteMessageText", ctx.Update.Message.From.LanguageCode), "delme-")))
	messageToSend.ReplyMarkup = rm
	messageToSend.ReplyToMessageID = ctx.Update.Message.MessageID
	ctx.Bot.Send(messageToSend)
}
