package main

import (
	"regexp"
	"scienzabot/consts"
	"scienzabot/database"
	"scienzabot/utils"
	"strconv"
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
		userExists       bool
		user             database.User
		userPermission   int
		userIsBotAdmin   bool
		userIsGroupAdmin bool
		groupStatus      int64
	)

	if userExists = ctx.Database.UserExists(message.From.ID); userExists {
		user, err = ctx.Database.GetUser(message.From.ID)
		ctx.Database.SetUserNickname(message.From.ID, message.From.UserName)
		ctx.Database.SetUserLocale(message.From.ID, message.From.LanguageCode)
		userIsBotAdmin = utils.HasPermission(int(user.Permissions), consts.UserPermissionAdmin)
		ctx.Database.UpdateUserLastSeen(message.From.ID, message.Time())
	}

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
			//Every 2000 messages reload chat admins
			if message.MessageID%2000 == 0 {
				reloadChatAdmins(ctx)
			}

			groupStatus, _ = ctx.Database.GetGroupStatus(message.Chat.ID)
			if !userIsBotAdmin && utils.HasPermission(int(groupStatus), consts.GroupBanned) {
				return
			}

			if userExists {
				userPermission, err = ctx.Database.GetPermission(int64(message.From.ID), message.Chat.ID)
				if err != nil {
					reloadChatAdmins(ctx)
					userPermission, _ = ctx.Database.GetPermission(int64(message.From.ID), message.Chat.ID)
				}
				userIsGroupAdmin = utils.HasPermission(userPermission, consts.UserPermissionGroupAdmin) || utils.HasPermission(userPermission, consts.UserPermissionAdmin)
				ctx.Database.IncrementMessageCount(int64(message.From.ID), message.Chat.ID)
			}
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
			replyMessageWithCloseButton(ctx, msg)
			break

		case "/help", "/aiuto", "/aiutami":
			if messageInGroup {
				messageBody = "onPrivateChatCommand"
			} else {
				messageBody = "helpCommand"
			}
			replyDbMessageWithCloseButton(ctx, messageBody)
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
				replyMessageWithCloseButton(ctx, val)
			}
			break

		case "/ping":
			if userIsBotAdmin {
				replyMessageWithCloseButton(ctx, "🏓 Pong!")
			}
			break

		case "/fulllist":
			if userIsBotAdmin {
				messageBody := ""
				if messageInGroup {
					lists, _ := ctx.Database.GetLists(message.Chat.ID)
					for _, lst := range lists {
						messageBody += lst.Name + "\n"
						users, _ := ctx.Database.GetSubscribedUsers(lst.ID)
						for i, usr := range users {
							user, _ := ctx.Database.GetUser(int(usr.UserID))
							if i == len(users)-1 {
								messageBody += "╚ "
							} else {
								messageBody += "╠ "
							}
							messageBody += user.Nickname + " [" + strconv.Itoa(int(user.ID)) + "]" + "\n"
						}
					}
					replyMessageWithCloseButton(ctx, messageBody)
				} else {
					replyDbMessageWithCloseButton(ctx, "notImplemented")
				}
			}
			break

		case "/gdpr":
			replyDbMessageWithCloseButton(ctx, "notImplemented")
			break

		case "/lists":
			if userExists {
				if !messageInGroup {
					grps, _ := ctx.Database.GetUserGroups(message.From.ID)
					messageBody := ""
					for _, group := range grps {
						messageBody += group.Title + "\n"
						lists, _ := ctx.Database.GetLists(group.ID)
						for i, lst := range lists {
							if i == len(lists)-1 {
								messageBody += "  ╚ " + lst.Name + "\n"
							} else {
								messageBody += "  ╠ " + lst.Name + "\n"
							}
						}
						messageBody += "\n"
					}
					messageToSend := tba.NewMessage(message.Chat.ID, messageBody)
					rm := tba.NewInlineKeyboardMarkup(
						tba.NewInlineKeyboardRow(
							tba.NewInlineKeyboardButtonData(ctx.Database.GetBotStringValueOrDefaultNoError("deleteMessageText", message.From.LanguageCode), "delme-")))
					//tba.NewInlineKeyboardButtonData(" ", "delme-")))
					messageToSend.ReplyMarkup = rm
					messageToSend.ReplyToMessageID = message.MessageID
					ctx.Bot.Send(messageToSend)
				} else {

					messageBody += message.Chat.Title + "\n"
					lists, _ := ctx.Database.GetLists(message.Chat.ID)
					for i, lst := range lists {
						if i == len(lists)-1 {
							messageBody += "  ╚ " + lst.Name + "\n"
						} else {
							messageBody += "  ╠ " + lst.Name + "\n"
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
			}
			break

		case "/newlist":

			if len(args) != 2 {
				replyDbMessageWithCloseButton(ctx, "newlistSyntaxError")
				return
			}
			listNameIsValid, _ := regexp.MatchString("^[a-z\\-_]{1,30}$", args[1])
			if !listNameIsValid {
				replyDbMessageWithCloseButton(ctx, "newlistSyntaxError")
				return
			}

			if messageInGroup {

				if userIsBotAdmin || userIsGroupAdmin || utils.HasPermission(userPermission, consts.UserPermissionCanCreateList) {

					err = ctx.Database.AddList(database.List{Name: args[1], GroupID: message.Chat.ID})
					if err == nil {
						replyDbMessageWithCloseButton(ctx, "listCreatedSuccessfully")
					}
				} else {
					replyDbMessageWithCloseButton(ctx, "notAuthorized")
				}

			} else {
				//TODO: implement group choosing where is admin
				replyDbMessageWithCloseButton(ctx, "notImplemented")
			}

			break

		case "/deletelist":

			if len(args) != 2 {
				replyDbMessageWithCloseButton(ctx, "deletelistSyntaxError")
				return
			}
			listNameIsValid, _ := regexp.MatchString("^[a-z\\-_]{1,30}$", args[1])
			if !listNameIsValid {
				replyDbMessageWithCloseButton(ctx, "deletelistSyntaxError")
				return
			}

			if messageInGroup {

				if userIsBotAdmin || userIsGroupAdmin || utils.HasPermission(userPermission, consts.UserPermissionCanRemoveList) {

					err = ctx.Database.DeleteListByName(message.Chat.ID, args[1])
					if err == nil {
						replyDbMessageWithCloseButton(ctx, "listDeletedSuccessfully")
					}
				} else {
					replyDbMessageWithCloseButton(ctx, "notAuthorized")
				}

			} else {
				//TODO: implement group choosing where is admin
				replyDbMessageWithCloseButton(ctx, "notImplemented")
			}

			break

		case "/del", "/deleteMessage":
			if message.ReplyToMessage != nil && (userIsBotAdmin || userIsGroupAdmin) {
				ctx.Bot.DeleteMessage(tba.NewDeleteMessage(message.Chat.ID, message.ReplyToMessage.MessageID))
			}
			break

		case "/registrazione", "/registra", "/registrati", "/registrami", "/signup":
			if !userExists {
				//We want registration to happen in private, not in public
				if messageInGroup {
					replyDbMessageWithCloseButton(ctx, "onPrivateChatCommand")
					//Warn about error?
					return
				}
				err = ctx.Database.AddUser(database.User{ID: int64(message.From.ID), Nickname: message.From.UserName, Status: consts.UserStatusActive, Locale: message.From.LanguageCode})
				if err != nil {
					replyDbMessageWithCloseButton(ctx, "generalError")
				} else {
					replyDbMessageWithCloseButton(ctx, "userAddedSuccessfully")
				}
			} else {
				replyDbMessageWithCloseButton(ctx, "userAlreadyRegistred")
			}
			break

		case "/iscrivi", "/iscrivimi", "/join", "/iscrizione", "/entra", "/sottoscrivi", "/subscribe":
			if userExists {
				//We want registration to happen in private, not in public
				if messageInGroup {
					//replyDbMessageWithCloseButton(ctx, "onPrivateChatCommand")

					lists, _ := ctx.Database.GetAvailableLists(message.Chat.ID, message.From.ID, consts.MaximumInlineKeyboardRows+1, 0)

					if len(lists) == 0 {
						replyDbMessageWithCloseButton(ctx, "noListsLeft")
						return
					}

					rows := make([][]tba.InlineKeyboardButton, 0)
					paginationPresent := false
					for i, lst := range lists {
						//if len(lists) > consts.MaximumInlineKeyboardRows && i+2 > consts.MaximumInlineKeyboardRows {
						if i+2 > consts.MaximumInlineKeyboardRows {
							rows = append(rows, []tba.InlineKeyboardButton{
								//tba.NewInlineKeyboardButtonData("‌‌ ", "ignore"),
								tba.NewInlineKeyboardButtonData(ctx.Database.GetBotStringValueOrDefaultNoError("closeMessageText", ctx.Update.Message.From.LanguageCode), "delme-"),
								tba.NewInlineKeyboardButtonData("➡️", "jo-"+strconv.Itoa(consts.MaximumInlineKeyboardRows-1))})
							paginationPresent = true
							break
						}
						rows = append(rows, []tba.InlineKeyboardButton{tba.NewInlineKeyboardButtonData(lst.Name, "sub-"+strconv.Itoa(int(lst.ID)))})
					}
					if !paginationPresent {
						rows = append(rows, []tba.InlineKeyboardButton{
							tba.NewInlineKeyboardButtonData(ctx.Database.GetBotStringValueOrDefaultNoError("closeMessageText", ctx.Update.Message.From.LanguageCode), "delme-"),
							tba.NewInlineKeyboardButtonData("‌‌ ", "ignore")})
					}

					replyMessageDBWithInlineKeyboard(ctx, "availableLists", tba.InlineKeyboardMarkup{InlineKeyboard: rows})
					return
				}
				//

				replyDbMessageWithCloseButton(ctx, "notImplemented")

				//err = ctx.Database.AddUser(database.User{ID: int64(message.From.ID), Nickname: message.From.UserName, Status: consts.UserStatusActive})
				/*
					if err != nil {
						replyDbMessageWithCloseButton(ctx, "generalError")
					} else {
						//replyDbMessageWithCloseButton(ctx, "userAddedSuccessfully")
						replyDbMessageWithCloseButton(ctx, "notImplemented")
					}
				*/

			} else {
				replyDbMessageWithCloseButton(ctx, "userNotRegistred")
			}
			break

		case "/unsubscribe", "/disicrivi", "/disicriviti":
			if userExists {
				//We want registration to happen in private, not in public
				if messageInGroup {
					//replyDbMessageWithCloseButton(ctx, "onPrivateChatCommand")

					//lists, _ := ctx.Database.GetUserLists()
					//message.Chat.ID, message.From.ID, consts.MaximumInlineKeyboardRows+1, 0
					lists, err := ctx.Database.GetUserGroupListsWithLimits(int64(message.From.ID), message.Chat.ID, consts.MaximumInlineKeyboardRows+1, 0)
					if err != nil {

					}

					if len(lists) == 0 {
						replyDbMessageWithCloseButton(ctx, "noSubscription")
						return
					}

					rows := make([][]tba.InlineKeyboardButton, 0)
					paginationPresent := false
					for i, lst := range lists {
						if i+2 > consts.MaximumInlineKeyboardRows {
							rows = append(rows, []tba.InlineKeyboardButton{
								//tba.NewInlineKeyboardButtonData("‌‌ ", "ignore"),
								tba.NewInlineKeyboardButtonData(ctx.Database.GetBotStringValueOrDefaultNoError("closeMessageText", ctx.Update.Message.From.LanguageCode), "delme-"),
								tba.NewInlineKeyboardButtonData("➡️", "uo-"+strconv.Itoa(consts.MaximumInlineKeyboardRows-1))})
							paginationPresent = true
							break
						}
						rows = append(rows, []tba.InlineKeyboardButton{tba.NewInlineKeyboardButtonData(lst.Name, "unsub-"+strconv.Itoa(int(lst.ID)))})
					}
					if !paginationPresent {
						rows = append(rows, []tba.InlineKeyboardButton{
							tba.NewInlineKeyboardButtonData(ctx.Database.GetBotStringValueOrDefaultNoError("closeMessageText", ctx.Update.Message.From.LanguageCode), "delme-"),
							tba.NewInlineKeyboardButtonData("‌‌ ", "ignore")})
					}

					replyMessageDBWithInlineKeyboard(ctx, "subscribedLists", tba.InlineKeyboardMarkup{InlineKeyboard: rows})
					return
				}
				//

				replyDbMessageWithCloseButton(ctx, "notImplemented")

				//err = ctx.Database.AddUser(database.User{ID: int64(message.From.ID), Nickname: message.From.UserName, Status: consts.UserStatusActive})
				/*
					if err != nil {
						replyDbMessageWithCloseButton(ctx, "generalError")
					} else {
						//replyDbMessageWithCloseButton(ctx, "userAddedSuccessfully")
						replyDbMessageWithCloseButton(ctx, "notImplemented")
					}
				*/

			} else {
				replyDbMessageWithCloseButton(ctx, "userNotRegistred")
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
	} else {
		//Check if list was invoked
		listPrefixes := []string{"@", "#", "!", "."}
		possibleLists := make([]string, 0)

		for _, prefix := range listPrefixes {
			if strings.Contains(message.Text, prefix) {
				words := strings.Split(message.Text, " ")
				for _, word := range words {
					if len(word) > 1 {
						if word[0] == prefix[0] {
							possibleLists = append(possibleLists, strings.ToLower(word[1:]))
						}
					}
				}
			}
		}

		if len(possibleLists) < 1 {
			return
		}

		groupLists, err := ctx.Database.GetLists(message.Chat.ID)
		if err != nil {
			replyDbMessageWithCloseButton(ctx, "generalError")
			return
		}

		lists := make([]database.List, 0)
		for _, plist := range possibleLists {
			for _, glist := range groupLists {
				if plist == glist.Name {
					lists = append(lists, glist)
				}
			}
		}

		if len(lists) < 1 {
			return
		}

		contactedUsers := make([]int64, 0)

		for _, list := range lists {
			subs, _ := ctx.Database.GetSubscribedUsers(list.ID)
			for _, sub := range subs {
				found := false
				for _, cUse := range contactedUsers {

					if sub.UserID == cUse {
						found = true

						break
					}

				}

				if !found {
					user, _ := ctx.Database.GetUser(int(sub.UserID))
					messageToSend := tba.NewMessage(sub.UserID, strings.Replace(strings.Replace(ctx.Database.GetBotStringValueOrDefaultNoError("tagNotification", user.Locale), "{{categoryName}}", list.Name, -1), "{{groupName}}", message.Chat.Title, -1))
					if message.Chat.IsSuperGroup() {
						if message.Chat.UserName != "" {
							ikm1 := tba.NewInlineKeyboardButtonURL(ctx.Database.GetBotStringValueOrDefaultNoError("tagNotificationGroupLink", user.Locale), "t.me/"+message.Chat.UserName)
							ikm2 := tba.NewInlineKeyboardButtonURL(ctx.Database.GetBotStringValueOrDefaultNoError("tagNotificationMessageLink", user.Locale), "t.me/"+message.Chat.UserName+"/"+strconv.Itoa(message.MessageID))
							ikm3 := tba.NewInlineKeyboardButtonData(ctx.Database.GetBotStringValueOrDefaultNoError("tagNotificationTag", user.Locale), "tag-"+strconv.FormatInt(message.Chat.ID, 10)+"-"+strconv.Itoa(message.MessageID))
							ikl := []tba.InlineKeyboardButton{ikm1, ikm2, ikm3}
							ikm := tba.NewInlineKeyboardMarkup(ikl)
							messageToSend.ReplyMarkup = ikm
						} else {
							ikm3 := tba.NewInlineKeyboardButtonData(ctx.Database.GetBotStringValueOrDefaultNoError("tagNotificationTag", user.Locale), "tag-"+strconv.FormatInt(message.Chat.ID, 10)+"-"+strconv.Itoa(message.MessageID))
							ikl := []tba.InlineKeyboardButton{ikm3}
							ikm := tba.NewInlineKeyboardMarkup(ikl)
							messageToSend.ReplyMarkup = ikm
						}
					} else {
						ikm3 := tba.NewInlineKeyboardButtonData(ctx.Database.GetBotStringValueOrDefaultNoError("tagNotificationTag", user.Locale), "tag-"+strconv.FormatInt(message.Chat.ID, 10)+"-"+strconv.Itoa(message.MessageID))
						ikl := []tba.InlineKeyboardButton{ikm3}
						ikm := tba.NewInlineKeyboardMarkup(ikl)
						messageToSend.ReplyMarkup = ikm
					}

					ctx.Bot.Send(messageToSend)

					contactedUsers = append(contactedUsers, sub.UserID)
				}
			}
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

func replyMessageDBWithInlineKeyboard(ctx *Context, keyString string, ikm tba.InlineKeyboardMarkup) {
	messageBody := ctx.Database.GetBotStringValueOrDefaultNoError(keyString, ctx.Update.Message.From.LanguageCode)
	messageToSend := tba.NewMessage(ctx.Update.Message.Chat.ID, messageBody)
	messageToSend.ReplyMarkup = ikm
	messageToSend.ReplyToMessageID = ctx.Update.Message.MessageID
	ctx.Bot.Send(messageToSend)
}

func replyMessageWithCloseButton(ctx *Context, messageBody string) {
	messageToSend := tba.NewMessage(ctx.Update.Message.Chat.ID, messageBody)
	rm := tba.NewInlineKeyboardMarkup(
		tba.NewInlineKeyboardRow(
			tba.NewInlineKeyboardButtonData(
				ctx.Database.GetBotStringValueOrDefaultNoError("deleteMessageText", ctx.Update.Message.From.LanguageCode), "delme-")))
	messageToSend.ReplyMarkup = rm
	messageToSend.ReplyToMessageID = ctx.Update.Message.MessageID
	ctx.Bot.Send(messageToSend)
}

func replyDbMessageWithCloseButton(ctx *Context, keyString string) {
	messageBody := ctx.Database.GetBotStringValueOrDefaultNoError(keyString, ctx.Update.Message.From.LanguageCode)
	messageToSend := tba.NewMessage(ctx.Update.Message.Chat.ID, messageBody)
	rm := tba.NewInlineKeyboardMarkup(
		tba.NewInlineKeyboardRow(
			tba.NewInlineKeyboardButtonData(
				ctx.Database.GetBotStringValueOrDefaultNoError("deleteMessageText", ctx.Update.Message.From.LanguageCode), "delme-")))
	messageToSend.ReplyMarkup = rm
	messageToSend.ReplyToMessageID = ctx.Update.Message.MessageID
	ctx.Bot.Send(messageToSend)
}
