package main

import (
	"scienzabot/consts"
	"scienzabot/database"
	"scienzabot/utils"
	"strconv"
	"strings"
	"time"

	tba "github.com/go-telegram-bot-api/telegram-bot-api"
)

// The callbackQueryHandler.go file contains the code behind a callback query from an
//	inline keyboard

func callbackQueryRoute(ctx *Context) {
	message := ctx.Update.CallbackQuery
	var (
		user           database.User
		isAdmin        bool
		userExists     bool
		messageInGroup bool
		err            error
		locale         string
	)

	locale = consts.DefaultLocale

	messageInGroup = message.Message != nil && message.Message.Chat.IsGroup() || message.Message.Chat.IsSuperGroup()

	if userExists = ctx.Database.UserExists(message.From.ID); userExists {
		user, err = ctx.Database.GetUser(message.From.ID)
		isAdmin = utils.HasPermission(int(user.Permissions), consts.UserPermissionAdmin)
		locale, _ = ctx.Database.GetUserLocale(ctx.Update.CallbackQuery.From.ID)
	}

	if message.From.LanguageCode != "" {
		locale = message.From.LanguageCode
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
										Text: ctx.Database.GetBotStringValueOrDefaultNoError("callbackQueryAnswerSuccess", locale)})

									//lists, _ := ctx.Database.GetAvailableLists(message.Message.Chat.ID, message.From.ID, consts.MaximumInlineKeyboardRows+1, 0)
									lists, _ := ctx.Database.GetUserGroupListsWithLimits(int64(message.From.ID), message.Message.Chat.ID, consts.MaximumInlineKeyboardRows+1, 0)

									if len(lists) == 0 {
										editInlineMessageDBWithCloseButton(ctx, "noSubscription")
										return
									}

									rows := make([][]tba.InlineKeyboardButton, 0)
									paginationPresent := false
									locale, _ := ctx.Database.GetUserLocale(message.Message.From.ID)
									for i, lst := range lists {
										if i+2 > consts.MaximumInlineKeyboardRows {
											rows = append(rows, []tba.InlineKeyboardButton{
												//tba.NewInlineKeyboardButtonData("‌‌ ", "ignore"),
												tba.NewInlineKeyboardButtonData(ctx.Database.GetBotStringValueOrDefaultNoError("closeMessageText", locale), "delme-"),
												tba.NewInlineKeyboardButtonData("➡️", "uo-"+strconv.Itoa(consts.MaximumInlineKeyboardRows-1))})
											paginationPresent = true
											break
										}
										rows = append(rows, []tba.InlineKeyboardButton{tba.NewInlineKeyboardButtonData(lst.Name, "unsub-"+strconv.Itoa(int(lst.ID)))})
									}
									if !paginationPresent {
										rows = append(rows, []tba.InlineKeyboardButton{
											tba.NewInlineKeyboardButtonData(ctx.Database.GetBotStringValueOrDefaultNoError("closeMessageText", locale), "delme-"),
											tba.NewInlineKeyboardButtonData("‌‌ ", "ignore")})
									}

									editInlineMessageInlineKeyboard(ctx, tba.InlineKeyboardMarkup{InlineKeyboard: rows})

								} else {
									ctx.Bot.AnswerCallbackQuery(tba.CallbackConfig{CallbackQueryID: message.ID, ShowAlert: true,
										Text: ctx.Database.GetBotStringValueOrDefaultNoError("callbackQueryAnswerError", locale)})
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
								leftOffset := offset - (consts.MaximumInlineKeyboardRows - 1)
								if leftOffset <= 0 {
									leftOffset = 0
								}
								leftBtn := tba.NewInlineKeyboardButtonData("⬅️", "uo-"+strconv.Itoa(leftOffset))
								closeBtn := tba.NewInlineKeyboardButtonData(ctx.Database.GetBotStringValueOrDefaultNoError("closeMessageText", locale), "delme-")
								if offset-leftOffset < consts.MaximumInlineKeyboardRows-1 {
									leftBtn = closeBtn
								}

								for i, lst := range lists {
									if i+2 > consts.MaximumInlineKeyboardRows {
										rows = append(rows, []tba.InlineKeyboardButton{
											//tba.NewInlineKeyboardButtonData("‌‌ ", "ignore"),
											leftBtn,
											tba.NewInlineKeyboardButtonData("➡️", "uo-"+strconv.Itoa(offset+consts.MaximumInlineKeyboardRows-1))})
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

								editInlineMessageInlineKeyboard(ctx, tba.InlineKeyboardMarkup{InlineKeyboard: rows})
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
										Text: ctx.Database.GetBotStringValueOrDefaultNoError("callbackQueryAnswerSuccess", locale)})

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
												tba.NewInlineKeyboardButtonData(ctx.Database.GetBotStringValueOrDefaultNoError("closeMessageText", locale), "delme-"),
												tba.NewInlineKeyboardButtonData("➡️", "jo-"+strconv.Itoa(consts.MaximumInlineKeyboardRows-1))})
											paginationPresent = true
											break
										}
										rows = append(rows, []tba.InlineKeyboardButton{tba.NewInlineKeyboardButtonData(lst.Name, "sub-"+strconv.Itoa(int(lst.ID)))})
									}
									if !paginationPresent {
										rows = append(rows, []tba.InlineKeyboardButton{
											tba.NewInlineKeyboardButtonData(ctx.Database.GetBotStringValueOrDefaultNoError("closeMessageText", locale), "delme-"),
											tba.NewInlineKeyboardButtonData("‌‌ ", "ignore")})
									}

									editInlineMessageInlineKeyboard(ctx, tba.InlineKeyboardMarkup{InlineKeyboard: rows})

								} else {
									ctx.Bot.AnswerCallbackQuery(tba.CallbackConfig{CallbackQueryID: message.ID, ShowAlert: true,
										Text: ctx.Database.GetBotStringValueOrDefaultNoError("callbackQueryAnswerError", locale)})
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
								leftOffset := offset - (consts.MaximumInlineKeyboardRows - 1)
								if leftOffset <= 0 {
									leftOffset = 0
								}
								leftBtn := tba.NewInlineKeyboardButtonData("⬅️", "jo-"+strconv.Itoa(leftOffset))
								closeBtn := tba.NewInlineKeyboardButtonData(ctx.Database.GetBotStringValueOrDefaultNoError("closeMessageText", locale), "delme-")
								rightBtn := tba.NewInlineKeyboardButtonData("‌‌ ", "ignore")
								if offset-leftOffset < consts.MaximumInlineKeyboardRows-1 {
									leftBtn = closeBtn
								}

								for i, lst := range lists {
									if len(lists) > consts.MaximumInlineKeyboardRows && i+2 > consts.MaximumInlineKeyboardRows {
										rows = append(rows, []tba.InlineKeyboardButton{
											//tba.NewInlineKeyboardButtonData("‌‌ ", "ignore"),
											leftBtn,
											tba.NewInlineKeyboardButtonData("➡️", "jo-"+strconv.Itoa(offset+consts.MaximumInlineKeyboardRows-1))})
										paginationPresent = true
										break
									}
									rows = append(rows, []tba.InlineKeyboardButton{tba.NewInlineKeyboardButtonData(lst.Name, "sub-"+strconv.Itoa(int(lst.ID)))})
								}
								if !paginationPresent {
									rows = append(rows, []tba.InlineKeyboardButton{
										leftBtn,
										rightBtn})
								}

								editInlineMessageInlineKeyboard(ctx, tba.InlineKeyboardMarkup{InlineKeyboard: rows})

							}
						}
					}
				}
			}
			break
			//Add error handler
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
							ctx.Database.GetBotStringValueOrDefaultNoError("closeMessageText", locale), "delme-")))

				replymessage.ReplyMarkup = rm
				replymessage.ParseMode = tba.ModeMarkdown
				_, err = ctx.Bot.Send(replymessage)
				if err != nil {
					ctx.Bot.AnswerCallbackQuery(tba.CallbackConfig{CallbackQueryID: message.ID,
						Text: ctx.Database.GetBotStringValueOrDefaultNoError("callbackQueryAnswerTagError", locale)})
				} else {
					ctx.Bot.AnswerCallbackQuery(tba.CallbackConfig{CallbackQueryID: message.ID,
						Text: ctx.Database.GetBotStringValueOrDefaultNoError("callbackQueryAnswerSuccess", locale)})
				}

			}

			break

		case "bgo":
			//bookmark group offset - shows groups
			//bgo-<offset>
			if userExists && !messageInGroup && len(args) == 2 {

				offset, err := strconv.Atoi(args[1])
				if err == nil {
					bms, err := ctx.Database.GetUserBookmarks(message.From.ID)
					if err != nil {
						return
					}
					if len(bms) > 0 {
						//If there are bookmarks
						lastGroupID := int64(-1)
						groups := make([]int64, 0)
						for _, b := range bms {
							if b.GroupID != lastGroupID {
								lastGroupID = b.GroupID
								groups = append(groups, lastGroupID)
							}
						}
						rows := make([][]tba.InlineKeyboardButton, 0)
						paginationPresent := false
						leftOffset := offset - (consts.MaximumInlineKeyboardRows - 1)
						if leftOffset <= 0 {
							leftOffset = 0
						}
						leftBtn := tba.NewInlineKeyboardButtonData("⬅️", "bgo-"+strconv.Itoa(leftOffset))
						closeBtn := tba.NewInlineKeyboardButtonData(ctx.Database.GetBotStringValueOrDefaultNoError("closeMessageText", locale), "delme-")
						rightBtn := tba.NewInlineKeyboardButtonData("‌‌ ", "ignore")
						if offset-leftOffset < consts.MaximumInlineKeyboardRows-1 {
							leftBtn = closeBtn
						}

						for i, g := range groups {
							//Skip the groups we already passed
							if offset > i {
								continue
							}
							if (i-offset)+2 > consts.MaximumInlineKeyboardRows {
								//If we are, we add as final row the pagination, to delete the message or show the next page
								rows = append(rows, []tba.InlineKeyboardButton{

									tba.NewInlineKeyboardButtonData(ctx.Database.GetBotStringValueOrDefaultNoError("closeMessageText", ctx.Update.Message.From.LanguageCode), "delme-"),
									//bookamrks groups offset
									tba.NewInlineKeyboardButtonData("➡️", "bgo-"+strconv.Itoa(consts.MaximumInlineKeyboardRows-1+offset))})
								//Then we set the bool to true to say that we added the pagination
								paginationPresent = true
								//And interrupt the loop
								break
							}
							b, err := ctx.Database.GetGroup(g)
							if err == nil {
								rows = append(rows, []tba.InlineKeyboardButton{tba.NewInlineKeyboardButtonData(b.Title, "bk-"+strconv.FormatInt(g, 10)+"-0")})
							}

						}
						//If the pagination was not added in the loop, we add it here, without adding the button to see the next page
						if !paginationPresent {
							rows = append(rows, []tba.InlineKeyboardButton{
								leftBtn,
								rightBtn})
						}
						//replyMessageDBWithInlineKeyboard(ctx, "bookmarksGroups", tba.InlineKeyboardMarkup{InlineKeyboard: rows})
						//editInlineMessageWithInlineKeyboard(ctx, tba.InlineKeyboardMarkup{InlineKeyboard: rows})
						editInlineMessageDBWithInlineKeyboard(ctx, "bookmarksMessage", tba.InlineKeyboardMarkup{InlineKeyboard: rows})
						//replyMessageDBWithInlineKeyboard(ctx, "bookmarksMessage", tba.InlineKeyboardMarkup{InlineKeyboard: rows})
					} else {
						//there are no bookmarks
					}

				}

			}
			break

		case "bk":
			//bk-<group>-<offset>
			if userExists && !messageInGroup && len(args) == 4 {
				offset := 0
				groupID := int64(0)
				if offset, err = strconv.Atoi(args[3]); err != nil {
					return
				}
				if groupID, err = strconv.ParseInt(args[2], 10, 64); err != nil {
					return
				}
				groupID *= -1

				if err == nil {
					bms, err := ctx.Database.GetUserGroupBookmarks(message.From.ID, groupID)
					if err != nil {
						return
					}
					if offset > len(bms)-1 {
						return
					}
					if len(bms) > 0 {
						//If there are bookmarks

						rows := make([][]tba.InlineKeyboardButton, 0)
						leftOffset := offset - 1
						if leftOffset < 0 {
							leftOffset = 0
						}
						leftBtn := tba.NewInlineKeyboardButtonData("‌‌ ", "ignore")
						if offset > 0 {
							leftBtn = tba.NewInlineKeyboardButtonData("⬅️", "bk-"+strconv.FormatInt(groupID, 10)+"-"+strconv.Itoa(leftOffset))
						}

						backBtn := tba.NewInlineKeyboardButtonData(ctx.Database.GetBotStringValueOrDefaultNoError("backText", locale), "bgo-0")
						if offset == 0 {
							leftBtn = backBtn
						}
						deleteBookmarkBtn := tba.NewInlineKeyboardButtonData("🗑", "bkd-"+strconv.FormatInt(groupID, 10)+"-"+strconv.Itoa(offset))
						tagMessageBookmarkBtn := tba.NewInlineKeyboardButtonData("📳", "tag-"+strconv.FormatInt(groupID, 10)+"-"+strconv.FormatInt(bms[offset].MessageID, 10))
						rightBtn := tba.NewInlineKeyboardButtonData("‌‌ ", "ignore")
						if len(bms)-1 > offset {
							rightBtn = tba.NewInlineKeyboardButtonData("➡️‌‌️️️️️️️", "bk-"+strconv.FormatInt(groupID, 10)+"-"+strconv.Itoa(offset+1))
						}

						rows = append(rows, []tba.InlineKeyboardButton{leftBtn, deleteBookmarkBtn, tagMessageBookmarkBtn, rightBtn})

						bookmarkedMessageSender, err := ctx.Database.GetUser(int(bms[offset].UserID))
						messageBody := "<b>" + ctx.Database.GetBotStringValueOrDefaultNoError("bookmarkFromText", locale) + "</b>: <a href=\"tg://user?id=" + strconv.Itoa(int(bms[offset].UserID)) + "\">"
						if err == nil {
							messageBody += bookmarkedMessageSender.Nickname
						} else {
							messageBody += strconv.Itoa(int(bms[offset].UserID))
						}
						messageBody += "</a>\n"
						if bms[offset].Alias != "" {
							messageBody += "<b>" + ctx.Database.GetBotStringValueOrDefaultNoError("bookmarkAliasText", locale) + "</b>: " + bms[offset].Alias + "\n"
						}
						location, _ := time.LoadLocation("Europe/Rome")
						messageBody += "<b>" + ctx.Database.GetBotStringValueOrDefaultNoError("bookmarkSavedonText", locale) + "</b>: " + bms[offset].CreationDate.In(location).Format("02/01/2006 15:04") + "\n"
						messageBody += "<b>" + ctx.Database.GetBotStringValueOrDefaultNoError("bookmarkContentText", locale) + "</b>: " + bms[offset].MessageContent

						editInlineMessageWithInlineKeyboard(ctx, messageBody, tba.InlineKeyboardMarkup{InlineKeyboard: rows}, tba.ModeHTML)
						//replyMessageDBWithInlineKeyboard(ctx, "bookmarksGroups", tba.InlineKeyboardMarkup{InlineKeyboard: rows})
						//editInlineMessageDBWithInlineKeyboard(ctx, tba.InlineKeyboardMarkup{InlineKeyboard: rows})
					} else {
						//there are no bookmarks
					}

				} //fi err==nil

			} //fi base args check

			break
		case "bkd":
			//bkd-<group>-<offset>
			//deletes a bookmark
			if userExists && !messageInGroup && len(args) == 4 {
				offset := 0
				groupID := int64(0)
				if offset, err = strconv.Atoi(args[3]); err != nil {
					return
				}
				if groupID, err = strconv.ParseInt(args[2], 10, 64); err != nil {
					return
				}
				groupID *= -1

				if err == nil {
					bms, err := ctx.Database.GetUserGroupBookmarks(message.From.ID, groupID)
					if err != nil {
						return
					}
					if len(bms) > 0 {
						//If there are bookmarks
						deleted := false
						for i := range bms {
							if i == offset {
								err = ctx.Database.DeleteBookmark(int(bms[i].ID))
								if err == nil {
									deleted = true
								}
							}
						}
						if deleted {

							ctx.Bot.AnswerCallbackQuery(tba.CallbackConfig{CallbackQueryID: message.ID,
								Text: ctx.Database.GetBotStringValueOrDefaultNoError("callbackQueryAnswerSuccess", locale)})
						} else {
							ctx.Bot.AnswerCallbackQuery(tba.CallbackConfig{CallbackQueryID: message.ID,
								Text: ctx.Database.GetBotStringValueOrDefaultNoError("callbackQueryAnswerError", locale)})

						}

						//If there are at least 2 (we removed 1)
						if len(bms)-2 > 0 {
							//Redirect to offset-1
							if offset-1 < 1 {
								offset = 0
							} else {
								offset--
							}
							ctx.Update.CallbackQuery.Data = "bk-" + strconv.FormatInt(bms[0].GroupID, 10) + "-" + strconv.Itoa(offset)
							callbackQueryRoute(ctx)
						} else if len(bms)-1 == 1 {
							ctx.Update.CallbackQuery.Data = "bk-" + strconv.FormatInt(bms[0].GroupID, 10) + "-0"
							callbackQueryRoute(ctx)
							//Redurect to group lists
						} else if len(bms)-1 == 0 {
							ctx.Update.CallbackQuery.Data = "bgo-0"
							callbackQueryRoute(ctx)
							//Redurect to group lists
						}
					} else {
						//there are no bookmarks
					}
				} //fi err==nil

			} //fi base args check

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
	locale, _ := ctx.Database.GetUserLocale(ctx.Update.CallbackQuery.From.ID)
	messageToSend := tba.NewEditMessageText(ctx.Update.CallbackQuery.Message.Chat.ID, ctx.Update.CallbackQuery.Message.MessageID, ctx.Database.GetBotStringValueOrDefaultNoError(key, locale))
	rm := tba.NewInlineKeyboardMarkup(
		tba.NewInlineKeyboardRow(
			tba.NewInlineKeyboardButtonData(
				ctx.Database.GetBotStringValueOrDefaultNoError("closeMessageText", locale), "delme-")))
	messageToSend.ReplyMarkup = &rm
	ctx.Bot.Send(messageToSend)
}
func editInlineMessageWithInlineKeyboard(ctx *Context, messageBody string, ikm tba.InlineKeyboardMarkup, parseMode string) {
	//locale, _ := ctx.Database.GetUserLocale(ctx.Update.CallbackQuery.From.ID)
	messageToSend := tba.NewEditMessageText(ctx.Update.CallbackQuery.Message.Chat.ID, ctx.Update.CallbackQuery.Message.MessageID, messageBody)
	rm := ikm
	messageToSend.ReplyMarkup = &rm
	messageToSend.ParseMode = parseMode
	ctx.Bot.Send(messageToSend)
}

func editInlineMessageDBWithInlineKeyboard(ctx *Context, dbKey string, ikm tba.InlineKeyboardMarkup) {
	locale, _ := ctx.Database.GetUserLocale(ctx.Update.CallbackQuery.From.ID)
	messageToSend := tba.NewEditMessageText(ctx.Update.CallbackQuery.Message.Chat.ID, ctx.Update.CallbackQuery.Message.MessageID,
		ctx.Database.GetBotStringValueOrDefaultNoError(dbKey, locale))
	messageToSend.ReplyMarkup = &ikm
	ctx.Bot.Send(messageToSend)
}

func editInlineMessageInlineKeyboard(ctx *Context, ikm tba.InlineKeyboardMarkup) {
	messageToSend := tba.NewEditMessageReplyMarkup(ctx.Update.CallbackQuery.Message.Chat.ID, ctx.Update.CallbackQuery.Message.MessageID, ikm)
	ctx.Bot.Send(messageToSend)
}
