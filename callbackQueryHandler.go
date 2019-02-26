package main

import (
	"scienzabot/consts"
	"scienzabot/database"
	"scienzabot/utils"
	"strconv"
	"strings"

	tba "github.com/go-telegram-bot-api/telegram-bot-api"
)

func callbackQueryRoute(ctx *Context) {
	message := ctx.Update.CallbackQuery
	var (
		user           database.User
		isAdmin        bool
		userExists     bool
		messageInGroup bool
		err            error
	)

	messageInGroup = message.Message != nil && message.Message.Chat.IsGroup() || message.Message.Chat.IsSuperGroup()

	if userExists = ctx.Database.UserExists(message.From.ID); userExists {
		user, err = ctx.Database.GetUser(message.From.ID)
		isAdmin = utils.HasPermission(int(user.Permissions), consts.UserPermissionAdmin)
	}

	switch {
	case message.Data != "" && strings.Contains(message.Data, "-"):
		switch args := strings.Split(message.Data, "-"); args[0] {

		case "unsub":
			if message.Message != nil && userExists {
				if messageInGroup {
					if message.Message.ReplyToMessage != nil && message.Message.ReplyToMessage.From.ID == message.From.ID {
						if len(args) == 2 {
							listID, err := strconv.Atoi(args[1])
							if err == nil {
								//TODO: check if user in group
								err = ctx.Database.RemoveSubscriptionByListAndUserID(listID, message.From.ID)
								if err == nil {
									//callbackQueryAnswerSuccess

									ctx.Bot.AnswerCallbackQuery(tba.CallbackConfig{CallbackQueryID: message.ID,
										Text: ctx.Database.GetBotStringValueOrDefaultNoError("callbackQueryAnswerSuccess", message.From.LanguageCode)})

									//lists, _ := ctx.Database.GetAvailableLists(message.Message.Chat.ID, message.From.ID, consts.MaximumInlineKeyboardRows+1, 0)
									lists, _ := ctx.Database.GetUserGroupListsWithLimits(int64(message.From.ID), message.Message.Chat.ID, consts.MaximumInlineKeyboardRows+1, 0)

									if len(lists) == 0 {
										editInlineMessageDBWithCloseButton(ctx, "noSubscription")
										return
									}

									rows := make([][]tba.InlineKeyboardButton, 0)
									paginationPresent := false
									for i, lst := range lists {
										if i+2 > consts.MaximumInlineKeyboardRows {
											rows = append(rows, []tba.InlineKeyboardButton{
												//tba.NewInlineKeyboardButtonData("‌‌ ", "ignore"),
												tba.NewInlineKeyboardButtonData(ctx.Database.GetBotStringValueOrDefaultNoError("deleteMessageText", message.Message.From.LanguageCode), "delme-"),
												tba.NewInlineKeyboardButtonData("➡️", "uo-"+strconv.Itoa(consts.MaximumInlineKeyboardRows-1))})
											paginationPresent = true
											break
										}
										rows = append(rows, []tba.InlineKeyboardButton{tba.NewInlineKeyboardButtonData(lst.Name, "unsub-"+strconv.Itoa(int(lst.ID)))})
									}
									if !paginationPresent {
										rows = append(rows, []tba.InlineKeyboardButton{
											tba.NewInlineKeyboardButtonData(ctx.Database.GetBotStringValueOrDefaultNoError("closeMessageText", message.Message.From.LanguageCode), "delme-"),
											tba.NewInlineKeyboardButtonData("‌‌ ", "ignore")})
									}

									editInlineMessageDBWithInlineKeyboard(ctx, tba.InlineKeyboardMarkup{InlineKeyboard: rows})

								} else {
									ctx.Bot.AnswerCallbackQuery(tba.CallbackConfig{CallbackQueryID: message.ID, ShowAlert: true,
										Text: ctx.Database.GetBotStringValueOrDefaultNoError("callbackQueryAnswerError", message.From.LanguageCode)})
								}
							}
						}
					}
				}
			}

			break

		case "uo":
			if message.Message != nil && userExists {
				if messageInGroup {
					if message.Message.ReplyToMessage != nil && message.Message.ReplyToMessage.From.ID == message.From.ID {
						if len(args) == 2 {

							offset, err := strconv.Atoi(args[1])
							if err == nil {
								//lists, _ := ctx.Database.GetUserLists(message.Message.Chat.ID, message.From.ID, consts.MaximumInlineKeyboardRows+1, offset)
								lists, _ := ctx.Database.GetUserGroupListsWithLimits(int64(message.From.ID), message.Message.Chat.ID, consts.MaximumInlineKeyboardRows+1, offset)

								rows := make([][]tba.InlineKeyboardButton, 0)
								paginationPresent := false
								leftOffset := offset - consts.MaximumInlineKeyboardRows - 1
								if leftOffset < 0 {
									leftOffset = 0
								}
								leftBtn := tba.NewInlineKeyboardButtonData("⬅️", "uo-"+strconv.Itoa(leftOffset))
								closeBtn := tba.NewInlineKeyboardButtonData(ctx.Database.GetBotStringValueOrDefaultNoError("deleteMessageText", message.Message.From.LanguageCode), "delme-")
								if offset < consts.MaximumInlineKeyboardRows-1 {
									leftBtn = closeBtn
								}

								for i, lst := range lists {
									if i+2 > consts.MaximumInlineKeyboardRows {
										rows = append(rows, []tba.InlineKeyboardButton{
											//tba.NewInlineKeyboardButtonData("‌‌ ", "ignore"),
											leftBtn,
											tba.NewInlineKeyboardButtonData("➡️", "uo-"+strconv.Itoa(consts.MaximumInlineKeyboardRows-1))})
										paginationPresent = true
										break
									}
									rows = append(rows, []tba.InlineKeyboardButton{tba.NewInlineKeyboardButtonData(lst.Name, "unsub-"+strconv.Itoa(int(lst.ID)))})
								}
								if !paginationPresent {
									rows = append(rows, []tba.InlineKeyboardButton{
										leftBtn,
										tba.NewInlineKeyboardButtonData("‌‌ ", "ignore")})
								}

								editInlineMessageDBWithInlineKeyboard(ctx, tba.InlineKeyboardMarkup{InlineKeyboard: rows})
								return
							}
						}
					}
				}
			}
			break

		case "sub":
			if message.Message != nil && userExists {
				if messageInGroup {
					if message.Message.ReplyToMessage != nil && message.Message.ReplyToMessage.From.ID == message.From.ID {
						if len(args) == 2 {
							listID, err := strconv.Atoi(args[1])
							if err == nil {
								//TODO: check if user in group
								err = ctx.Database.AddSubscription(message.From.ID, listID)
								if err == nil {
									//callbackQueryAnswerSuccess

									ctx.Bot.AnswerCallbackQuery(tba.CallbackConfig{CallbackQueryID: message.ID,
										Text: ctx.Database.GetBotStringValueOrDefaultNoError("callbackQueryAnswerSuccess", message.From.LanguageCode)})

									lists, _ := ctx.Database.GetAvailableLists(message.Message.Chat.ID, message.From.ID, consts.MaximumInlineKeyboardRows+1, 0)

									if len(lists) == 0 {
										editInlineMessageDBWithCloseButton(ctx, "noListsLeft")
										return
									}

									rows := make([][]tba.InlineKeyboardButton, 0)
									paginationPresent := false
									for i, lst := range lists {
										if len(lists) > consts.MaximumInlineKeyboardRows && i+2 > consts.MaximumInlineKeyboardRows {
											rows = append(rows, []tba.InlineKeyboardButton{
												//tba.NewInlineKeyboardButtonData("‌‌ ", "ignore"),
												tba.NewInlineKeyboardButtonData(ctx.Database.GetBotStringValueOrDefaultNoError("deleteMessageText", message.Message.From.LanguageCode), "delme-"),
												tba.NewInlineKeyboardButtonData("➡️", "jo-"+strconv.Itoa(consts.MaximumInlineKeyboardRows-1))})
											paginationPresent = true
											break
										}
										rows = append(rows, []tba.InlineKeyboardButton{tba.NewInlineKeyboardButtonData(lst.Name, "sub-"+strconv.Itoa(int(lst.ID)))})
									}
									if !paginationPresent {
										rows = append(rows, []tba.InlineKeyboardButton{
											tba.NewInlineKeyboardButtonData(ctx.Database.GetBotStringValueOrDefaultNoError("closeMessageText", message.Message.From.LanguageCode), "delme-"),
											tba.NewInlineKeyboardButtonData("‌‌ ", "ignore")})
									}

									editInlineMessageDBWithInlineKeyboard(ctx, tba.InlineKeyboardMarkup{InlineKeyboard: rows})

								} else {
									ctx.Bot.AnswerCallbackQuery(tba.CallbackConfig{CallbackQueryID: message.ID, ShowAlert: true,
										Text: ctx.Database.GetBotStringValueOrDefaultNoError("callbackQueryAnswerError", message.From.LanguageCode)})
								}
							}
						}
					}
				}
			}

			break

		case "jo":
			if message.Message != nil && userExists {
				if messageInGroup {
					if message.Message.ReplyToMessage != nil && message.Message.ReplyToMessage.From.ID == message.From.ID {
						if len(args) == 2 {

							offset, err := strconv.Atoi(args[1])
							if err == nil {
								lists, _ := ctx.Database.GetAvailableLists(message.Message.Chat.ID, message.From.ID, consts.MaximumInlineKeyboardRows+1, offset)

								rows := make([][]tba.InlineKeyboardButton, 0)
								paginationPresent := false
								leftOffset := offset - consts.MaximumInlineKeyboardRows - 1
								if leftOffset < 0 {
									leftOffset = 0
								}
								leftBtn := tba.NewInlineKeyboardButtonData("⬅️", "jo-"+strconv.Itoa(leftOffset))
								closeBtn := tba.NewInlineKeyboardButtonData(ctx.Database.GetBotStringValueOrDefaultNoError("deleteMessageText", message.Message.From.LanguageCode), "delme-")
								if offset < consts.MaximumInlineKeyboardRows-1 {
									leftBtn = closeBtn
								}

								for i, lst := range lists {
									if len(lists) > consts.MaximumInlineKeyboardRows && i+2 > consts.MaximumInlineKeyboardRows {
										rows = append(rows, []tba.InlineKeyboardButton{
											//tba.NewInlineKeyboardButtonData("‌‌ ", "ignore"),
											leftBtn,
											tba.NewInlineKeyboardButtonData("➡️", "jo-"+strconv.Itoa(consts.MaximumInlineKeyboardRows-1))})
										paginationPresent = true
										break
									}
									rows = append(rows, []tba.InlineKeyboardButton{tba.NewInlineKeyboardButtonData(lst.Name, "sub-"+strconv.Itoa(int(lst.ID)))})
								}
								if !paginationPresent {
									rows = append(rows, []tba.InlineKeyboardButton{
										leftBtn,
										tba.NewInlineKeyboardButtonData("‌‌ ", "ignore")})
								}

								editInlineMessageDBWithInlineKeyboard(ctx, tba.InlineKeyboardMarkup{InlineKeyboard: rows})

							}
						}
					}
				}
			}
			break

		case "tag":
			//If the messge is not null AND the user is admin OR the bot is replaying to the message sent to the user that clicked the button
			if message.Message != nil && len(args) == 4 && args[1] == "" && message.From.UserName != "" {
				groupID, err := strconv.ParseInt(args[2], 10, 64)
				if err != nil {
					break
				}

				messageID, err := strconv.ParseInt(args[3], 10, 64)
				if err != nil {
					break
				}

				replymessage := tba.NewMessage(groupID*-1, "[@"+message.From.FirstName+" "+message.From.LastName+"](tg://user?id="+strconv.Itoa(message.From.ID)+")")
				replymessage.ReplyToMessageID = int(messageID)

				rm := tba.NewInlineKeyboardMarkup(
					tba.NewInlineKeyboardRow(
						tba.NewInlineKeyboardButtonData(
							ctx.Database.GetBotStringValueOrDefaultNoError("deleteMessageText", message.From.LanguageCode), "delme-")))

				replymessage.ReplyMarkup = rm
				replymessage.ParseMode = tba.ModeMarkdown
				ctx.Bot.Send(replymessage)

			}

			break

		case "delme":
			//If the messge is not null AND the user is admin OR the bot is replaying to the message sent to the user that clicked the button
			if message.Message != nil &&
				(isAdmin ||
					(message.Message.ReplyToMessage != nil &&
						message.Message.ReplyToMessage.From.ID == message.From.ID)) {
				ctx.Bot.DeleteMessage(tba.DeleteMessageConfig{ChatID: message.Message.Chat.ID, MessageID: message.Message.MessageID})
			}

			break

		case "del":
			var groupID int64
			groupID, err = strconv.ParseInt(args[1], 10, 64)
			msgToDelete, err2 := strconv.Atoi(args[2])
			if err != nil && err2 != nil && len(args) == 3 {
				if isAdmin {
					ctx.Bot.DeleteMessage(tba.DeleteMessageConfig{ChatID: groupID, MessageID: msgToDelete})
				}
			}
			break

		}
		break

	}

}

func editInlineMessageDBWithCloseButton(ctx *Context, key string) {
	messageToSend := tba.NewEditMessageText(ctx.Update.CallbackQuery.Message.Chat.ID, ctx.Update.CallbackQuery.Message.MessageID, ctx.Database.GetBotStringValueOrDefaultNoError(key, ctx.Update.CallbackQuery.Message.From.LanguageCode))
	ctx.Bot.Send(messageToSend)
}
func editInlineMessageDBWithInlineKeyboard(ctx *Context, ikm tba.InlineKeyboardMarkup) {
	messageToSend := tba.NewEditMessageReplyMarkup(ctx.Update.CallbackQuery.Message.Chat.ID, ctx.Update.CallbackQuery.Message.MessageID, ikm)
	messageToSend.ReplyMarkup = &ikm
	ctx.Bot.Send(messageToSend)
}
