package main

import (
	"regexp"
	"scienzabot/consts"
	"scienzabot/database"
	"scienzabot/utils"
	"strconv"
	"strings"
	"time"

	tba "github.com/go-telegram-bot-api/telegram-bot-api"
)

// The textMessageHandler.go file is the route chosen for a text message sent

func textMessageRoute(ctx *Context) {
	message := ctx.Update.Message

	//For the moment we don't care about channels
	if message.Chat.IsChannel() {
		return
	}

	// Delcasing the variables
	var (
		//Err is a general error
		err error
		//messageBody is a variable that will contain an eventual body of a text message
		messageBody string
		//userExists tells if the user exists in the database
		userExists bool
		//If the user exists in the DB, the user variable will contain its representation in the database
		user database.User
		//userPermission is the permission of the user in the chat, if the message was sent in a group
		userPermission int
		//userIsBotAdmin tells if the user is a bot admin (and can do special things)
		userIsBotAdmin bool
		//userIsGroupAdmin is set if, IN THE DATABASE the user is admin in the group (via teelegram permissions)
		userIsGroupAdmin bool
		//groupStatus is the status of the group. it can indicate if a group is banned
		groupStatus int64
	)

	//If the user exists in the database, set the variable to true and update some values in the DB:
	if userExists = ctx.Database.UserExists(message.From.ID); userExists {
		//Get the user representation from the database
		user, err = ctx.Database.GetUser(message.From.ID)
		//Update the user nickname, in case he changed it
		ctx.Database.SetUserNickname(message.From.ID, message.From.UserName)
		//Update the user locale, in case he changed the language in use
		ctx.Database.SetUserLocale(message.From.ID, message.From.LanguageCode)
		//see if the user is admin of the bot
		userIsBotAdmin = utils.HasPermission(int(user.Permissions), consts.UserPermissionAdmin)
		//Update the last time the user was seen
		ctx.Database.UpdateUserLastSeen(message.From.ID, message.Time())
	}

	//See if the message was sent in a group or a supergroup
	messageInGroup := message.Chat.IsGroup() || message.Chat.IsSuperGroup()
	//If so
	if messageInGroup {
		//If the group does not exist, add it to the database
		if !ctx.Database.GroupExists(message.Chat.ID) {
			//Get the invite link from the chat, if it exists
			//TODO: the bot needs to generate the invite link by itself
			ref := message.Chat.InviteLink
			//If the invite link was not found
			if ref == "" {
				//If the chat has a public username, use it as a ref
				if message.Chat.UserName != "" {
					ref = "https://t.me/" + message.Chat.UserName
				}
			}
			//Add the group to the database
			ctx.Database.AddGroup(database.Group{ID: message.Chat.ID, Title: message.Chat.Title, Ref: message.Chat.UserName})
			//Reload the admins in the group
			reloadChatAdmins(ctx)

		}

		//Every 2000 messages reload chat admins
		if message.MessageID%2000 == 0 {
			reloadChatAdmins(ctx)
		}

		//Get the group status from the database
		groupStatus, _ = ctx.Database.GetGroupStatus(message.Chat.ID)
		//If the user who sent the message is not a bot admin and the group is banned, ignore the message
		if !userIsBotAdmin && utils.HasPermission(int(groupStatus), consts.GroupBanned) {
			return
		}

		//If the user exists in database
		if userExists {
			//Get its permission in the group
			userPermission, err = ctx.Database.GetPermission(int64(message.From.ID), message.Chat.ID)
			//If there's an error (like no result set), reload the admins in the group and retry
			if err != nil {
				reloadChatAdmins(ctx)
				userPermission, _ = ctx.Database.GetPermission(int64(message.From.ID), message.Chat.ID)
			}
			//See if the user is an administrator of the group
			userIsGroupAdmin = utils.HasPermission(userPermission, consts.UserPermissionGroupAdmin) || utils.HasPermission(userPermission, consts.UserPermissionAdmin)
			//Increment the message count of the user in the group
			ctx.Database.IncrementMessageCount(message.From.ID, message.Chat.ID)
			ctx.Database.UpdateLastSeen(message.From.ID, message.Chat.ID, message.Time())
		}

	}

	//If the message is a command
	if message.IsCommand() {
		//Command

		//Remove an eventual @botusername from the string
		// Then split the message using spaces as separators and use the switch to select the command, if exists
		switch args := strings.Split(strings.Replace(message.Text, "@"+ctx.Bot.Self.UserName, "", 1), " "); args[0] {
		case "/start":
			break

		//Exists just shows if a user exists in the database
		case "/exists":
			msg := "You do "
			if !userExists {
				msg += "not "
			}
			msg += "exist."
			replyMessageWithCloseButton(ctx, msg)
			break

		// Help message
		case "/help", "/aiuto", "/aiutami":
			if messageInGroup {
				messageBody = "onPrivateChatCommand"
			} else {
				messageBody = "helpCommand"
			}
			replyDbMessageWithCloseButton(ctx, messageBody)
			break

		//Info message
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

		//Get bot version
		case "/version", "/v":
			if val, err := ctx.Database.GetBotSettingValue("version"); err == nil {
				replyMessageWithCloseButton(ctx, val)
			}
			break

		//If the user is a bot admin, he can ping the bot
		case "/ping":
			if userIsBotAdmin {
				replyMessageWithCloseButton(ctx, "🏓 Pong!")
			}
			break

		//Fulllist returns all the user subscribed to all the lists
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
					//Get all the groups in the database
					groups, _ := ctx.Database.GetGroups()
					//For each group
					for _, group := range groups {
						//Get the lists in the group
						lists, _ := ctx.Database.GetLists(group.ID)
						//Add the name of the group to the message
						messageBody += group.Title + "\n"
						//for each list
						for _, lst := range lists {
							//Add to the message the name of the list
							messageBody += lst.Name + "\n"
							//And get all the users subscribed to the list
							users, _ := ctx.Database.GetSubscribedUsers(lst.ID)
							//for each user
							for i, usr := range users {
								//Get the user details
								user, _ := ctx.Database.GetUser(int(usr.UserID))
								//If the user is the last one in the list, prepend the conclusive char
								if i == len(users)-1 {
									messageBody += "╚ "
								} else {
									//Otherwise use the "normal" one
									messageBody += "╠ "
								}
								//Then append to the message body the message the user nickame and its ID
								messageBody += user.Nickname + " [" + strconv.Itoa(int(user.ID)) + "]" + "\n"
							} //End of user loop
						} //End of list loop
						//Add 2 carriage returns to the message body
						messageBody += "\n\n"
					} //end of group loop

					//Send the message
					replyMessageWithCloseButton(ctx, messageBody)
				}
			}
			break

		case "/gdpr":
			replyDbMessageWithCloseButton(ctx, "notImplemented")
			break

		case "/lists":
			//If the user is registered
			if userExists {
				//If the user is not in a group
				if !messageInGroup {
					//Get all the groups the user is in
					grps, _ := ctx.Database.GetUserGroups(message.From.ID)
					messageBody := ""
					//For each group the user is in
					for _, group := range grps {
						//Add the group name to the message body
						messageBody += group.Title + "\n"
						//Get all the lists in the group
						lists, _ := ctx.Database.GetLists(group.ID)
						//For each list in the group
						for i, lst := range lists {
							//if the list is the latest one, prepend the last char
							if i == len(lists)-1 {
								messageBody += "  ╚ "
							} else {
								messageBody += "  ╠ "
							}
							messageBody += lst.Name + "\n"
						} //end of list loop
						//Add a space between the groups
						messageBody += "\n"
					} //end of group loop
					//Create the message to send
					messageToSend := tba.NewMessage(message.Chat.ID, messageBody)
					//Create the button to delete the message
					rm := tba.NewInlineKeyboardMarkup(
						tba.NewInlineKeyboardRow(
							tba.NewInlineKeyboardButtonData(ctx.Database.GetBotStringValueOrDefaultNoError("deleteMessageText", message.From.LanguageCode), "delme-")))
					//Add the button to the message
					messageToSend.ReplyMarkup = rm
					//Add the message to reply to
					messageToSend.ReplyToMessageID = message.MessageID
					//Send the message
					ctx.Bot.Send(messageToSend)
				} else { //Message is in a group - write just the group lists
					//Add to the message body the chat title
					messageBody += message.Chat.Title + "\n"
					//Get the lists in the group
					lists, _ := ctx.Database.GetLists(message.Chat.ID)
					//For each list
					for i, lst := range lists {
						//If it's the latest list prepend the final char
						if i == len(lists)-1 {
							messageBody += "  ╚ "
						} else {
							//Otherwhise add the standard one
							messageBody += "  ╠ "
						}
						//Add to the message body the list name
						messageBody += lst.Name + "\n"
					} //end of list loop

					//Create the message to send
					messageToSend := tba.NewMessage(message.Chat.ID, messageBody)
					//Create the keyboard
					rm := tba.NewInlineKeyboardMarkup(
						tba.NewInlineKeyboardRow(
							tba.NewInlineKeyboardButtonData(ctx.Database.GetBotStringValueOrDefaultNoError("deleteMessageText", message.From.LanguageCode), "delme-")))
					//assign the keyboard to the message
					messageToSend.ReplyMarkup = rm
					//Set the message to reply to as the command message
					messageToSend.ReplyToMessageID = message.MessageID
					//Send the message
					ctx.Bot.Send(messageToSend)
				} //End user is in group
			} //end user exists
			break

			//Newlist is created to create a new list in a group
		case "/newlist":

			//If the commands are not 2, the syntax is wrong
			if len(args) != 2 {
				//Reply to the message with the syntax error string
				replyDbMessageWithCloseButton(ctx, "newlistSyntaxError")
				//And stop elaborating the message
				return
			}
			//Get if the list name is valid via regex expression of the 2nd argument (/newlist listname <- 2nd argument)
			listNameIsValid, _ := regexp.MatchString("^[a-z\\-_]{1,30}$", args[1])
			//If the list name is invalid
			if !listNameIsValid {
				//Return the message syntax message
				replyDbMessageWithCloseButton(ctx, "newlistSyntaxError")
				//And stop elaborating the message
				return
			}

			//If the message is in a group
			if messageInGroup {
				//If the message comes from a bot admin, a group admin or a user that has the permission to create a new list
				if userIsBotAdmin || userIsGroupAdmin || utils.HasPermission(userPermission, consts.UserPermissionCanCreateList) {
					//Try adding the new list
					err = ctx.Database.AddList(database.List{Name: args[1], GroupID: message.Chat.ID})
					//If there was no error
					if err == nil {
						//Say that the list was created successfully
						replyDbMessageWithCloseButton(ctx, "listCreatedSuccessfully")
					} //TODO: add else condition; if query failed, probably the list name is occupied
				} else { //If the user is not authorized
					//Say that he's not authorized
					replyDbMessageWithCloseButton(ctx, "notAuthorized")
				}

			} else {
				//Need to implement a way of choosing the group to add the list to from the groups the
				//	user is in and has permission to add lists to
				//Not urgent
				//
				//TODO: implement group choosing where is admin
				replyDbMessageWithCloseButton(ctx, "notImplemented")
			}

			break

		case "/deletelist":
			//If the commands are not 2, the syntax is wrong
			if len(args) != 2 {
				//Reply to the message with the syntax error string
				replyDbMessageWithCloseButton(ctx, "deletelistSyntaxError")
				//And stop elaborating the message
				return
			}
			//Reply to the message with the syntax error string
			listNameIsValid, _ := regexp.MatchString("^[a-z\\-_]{1,30}$", args[1])
			if !listNameIsValid {
				replyDbMessageWithCloseButton(ctx, "deletelistSyntaxError")
				//And stop elaborating the message
				return
			}

			//If the message is in a group
			if messageInGroup {
				//If the message comes from a bot admin, a group admin or a user that has the permission to delete a list
				if userIsBotAdmin || userIsGroupAdmin || utils.HasPermission(userPermission, consts.UserPermissionCanRemoveList) {
					//Try deleting the new list
					err = ctx.Database.DeleteListByName(message.Chat.ID, args[1])
					//If there was no error
					if err == nil {
						//Say that the list was deleted successfully
						replyDbMessageWithCloseButton(ctx, "listDeletedSuccessfully")
					} //TODO: add else condition; if query failed, probably the list does not exists
					//TODO: also need to create the string
				} else { //If the user is not authorized, show the message of authorization
					replyDbMessageWithCloseButton(ctx, "notAuthorized")
				}

			} else {
				//Need to implement a way of choosing the group to remove the list from, from the groups the
				//	user is in and has permission to remove lists from
				//Not urgent
				//
				//TODO: implement group choosing where is admin
				replyDbMessageWithCloseButton(ctx, "notImplemented")
			}

			break

		case "/del", "/deleteMessage":
			//If the message refers to another message and the user is bot admin or group admin
			if message.ReplyToMessage != nil && (userIsBotAdmin || userIsGroupAdmin) {
				//Delete the message the user is replying to
				ctx.Bot.DeleteMessage(tba.NewDeleteMessage(message.Chat.ID, message.ReplyToMessage.MessageID))
			}
			break

		case "/registrazione", "/registra", "/registrati", "/registrami", "/signup":
			//If the user is not in database
			if !userExists {
				//We want registration to happen in private, not in public
				if messageInGroup {
					replyDbMessageWithCloseButton(ctx, "onPrivateChatCommand")
					return
				}
				//If the message is in private, try to add the user to the database
				err = ctx.Database.AddUser(database.User{ID: int64(message.From.ID), Nickname: message.From.UserName, Status: consts.UserStatusActive, Locale: message.From.LanguageCode})
				//if there was an error
				if err != nil {
					//Warn about a general error
					replyDbMessageWithCloseButton(ctx, "generalError")
				} else {
					//otherwise say that the user was added successfully
					replyDbMessageWithCloseButton(ctx, "userAddedSuccessfully")
				}
			} else { //user is already registered
				//say that the user is already registered
				replyDbMessageWithCloseButton(ctx, "userAlreadyRegistred")
			}
			break

		//Add the user to a list
		case "/iscrivi", "/iscrivimi", "/join", "/iscrizione", "/entra", "/sottoscrivi", "/subscribe":
			//if the user exists in the database
			if userExists {
				//If the message is in a group, we already know the group to subscribe the user to
				if messageInGroup {
					//Get the available lists in the group
					//The 3rd parameter is how many list we want to get, the 4th is the "offset" (how many lists to skip)
					//TODO: Check if comment is clear
					//The 3rd "magic number" is the maximum amount of rows +1
					//  this because we want to know if there are more list to add, so we just ask for an additional list
					//  and if it exists, we don't show it, but instead we show the pagination button, to see the other lists
					lists, _ := ctx.Database.GetAvailableLists(message.Chat.ID, message.From.ID, consts.MaximumInlineKeyboardRows+1, 0)
					//If there's no list left, reply with another message
					if len(lists) == 0 {
						replyDbMessageWithCloseButton(ctx, "noListsLeft")
						return
					}
					//We then create the inline keyboard
					rows := make([][]tba.InlineKeyboardButton, 0)
					//We then declare a bool that will state if the pagination was added in the loop or not
					//  in case we will need to add it later
					paginationPresent := false
					//Then, we iterate the lists
					for i, lst := range lists {
						//For each iteration, we check if the iteration number exceed the maximum row number
						//i+2 because i starts from 0, so we need to add 1 to have the number of the current button,
						// and +1 to check if we are exceeding the maximum buttons number
						if i+2 > consts.MaximumInlineKeyboardRows {
							//If we are, we add as final row the pagination, to delete the message or show the next page
							rows = append(rows, []tba.InlineKeyboardButton{

								tba.NewInlineKeyboardButtonData(ctx.Database.GetBotStringValueOrDefaultNoError("closeMessageText", ctx.Update.Message.From.LanguageCode), "delme-"),
								tba.NewInlineKeyboardButtonData("➡️", "jo-"+strconv.Itoa(consts.MaximumInlineKeyboardRows-1))})
							//Then we set the bool to true to say that we added the pagination
							paginationPresent = true
							//And interrupt the loop
							break
						}
						//if the list number is not exceeding the maximum button number, we add the list name
						//KeyboardButtonData is just a way to pass a string to the bot itself
						//  The login behind the button data is shown in the callbackQueryhandler.go file
						rows = append(rows, []tba.InlineKeyboardButton{tba.NewInlineKeyboardButtonData(lst.Name, "sub-"+strconv.Itoa(int(lst.ID)))})
					}
					//If the pagination was not added in the loop, we add it here, without adding the button to see the next page
					if !paginationPresent {
						rows = append(rows, []tba.InlineKeyboardButton{
							tba.NewInlineKeyboardButtonData(ctx.Database.GetBotStringValueOrDefaultNoError("closeMessageText", ctx.Update.Message.From.LanguageCode), "delme-"),
							tba.NewInlineKeyboardButtonData("‌‌ ", "ignore")})
					}
					//Then we send the message
					replyMessageDBWithInlineKeyboard(ctx, "availableLists", tba.InlineKeyboardMarkup{InlineKeyboard: rows})
					return
				} //fi messageInGroup

				//The message is not in a group
				replyDbMessageWithCloseButton(ctx, "notImplemented")

			} else { //User is not in DB
				replyDbMessageWithCloseButton(ctx, "userNotRegistred")
			}
			break

		case "/unsubscribe", "/disicrivi", "/disicriviti":
			//If the user exists
			if userExists {
				//We want registration to happen in private, not in public
				//And the message is in a group (so we already know the group)
				if messageInGroup {
					//We get the lists the user is subscribed to
					lists, err := ctx.Database.GetUserGroupListsWithLimits(int64(message.From.ID), message.Chat.ID, consts.MaximumInlineKeyboardRows+1, 0)
					//If an error verified, report it
					if err != nil {
						replyDbMessageWithCloseButton(ctx, "generalError")
						return
					}
					//If the user is not subscribed to any list, send the message and stop here
					if len(lists) == 0 {
						replyDbMessageWithCloseButton(ctx, "noSubscription")
						return
					}
					//From here the procedure is the same as with the join command, so I'm not copying and pasting the comments
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
				} //fi message in group

				replyDbMessageWithCloseButton(ctx, "notImplemented")

			} else { //Uer not registered
				replyDbMessageWithCloseButton(ctx, "userNotRegistred")
			}
			break

		case "/segnalibro", "/salva":

			break

		//The listinterval is the interval between lists
		case "/listinterval":
			//If the message is in a group and the user is a botadmin or a groupadmin
			if messageInGroup && (userIsGroupAdmin || userIsBotAdmin) {
				//If the args are 2
				if len(args) == 2 {
					//Try to parse the 2nd argument as a time duration
					_, err := time.ParseDuration(args[1])
					if err != nil {
						//There was an error, send the syntax command
						replyDbMessageWithCloseButton(ctx, "listintervalSyntaxError")
						//And return
						return
					}
					//If there was no issue, update the setting in the database
					err = ctx.Database.SetSettingValue("listInterval", args[1], int(message.Chat.ID))
					if err == nil {
						replyDbMessageWithCloseButton(ctx, "listintervalSuccess")
						//Send success command
					} else {
						//Something went wrong but the duration parsed successfully, check the logs and send error message
						replyDbMessageWithCloseButton(ctx, "generalError")
					}
				} else {
					//The args are not 2, send the syntax command
					replyDbMessageWithCloseButton(ctx, "listintervalSyntaxError")
				}
			}
			break

		//The userinterval is the interval between invocations from a user
		case "/userinterval":
			//If the message is in a group and the user is a botadmin or a groupadmin
			if messageInGroup && (userIsGroupAdmin || userIsBotAdmin) {
				//If the args are 2
				if len(args) == 2 {
					//Try to parse the 2nd argument as a time duration
					_, err := time.ParseDuration(args[1])
					if err != nil {
						//There was an error, send the syntax command
						replyDbMessageWithCloseButton(ctx, "userintervalSyntaxError")
						//And return
						return
					}
					//If there was no issue, update the setting in the database
					err = ctx.Database.SetSettingValue("userInterval", args[1], int(message.Chat.ID))
					if err == nil {
						replyDbMessageWithCloseButton(ctx, "userintervalSuccess")
						//Send success command
					} else {
						//Something went wrong but the duration parsed successfully, check the logs and send error message
						replyDbMessageWithCloseButton(ctx, "generalError")
					}
				} else {
					//The args are not 2, send the syntax command
					replyDbMessageWithCloseButton(ctx, "userintervalSyntaxError")
				}
			}
			break

		//restart is used to reload the telegram admins within a group
		case "/reloadpermissions", "/ricarica", "/riavvia", "/restart":
			reloadChatAdmins(ctx)
			break

		default:
			//Check if it exists in DB

		}
	} else { //The message is not a command

		//If the message is in a group we can check for thins like lists invocations etc
		if messageInGroup {
			//Get the user interval if present
			userIntervalString, userIntervalError := ctx.Database.GetSettingValue("userInterval", int(message.Chat.ID))
			//If it's nil the setting exists
			if userIntervalError == nil {
				//Convert the string to a timespan
				userInterval, _ := time.ParseDuration(userIntervalString)
				//Get the last time the user invoked a list
				lastInvocation, _ := ctx.Database.GetLastListInvocation(message.From.ID, message.Chat.ID)
				//If the time is greater than the message, the user shouldn't be able to call a list and should be ignored
				if userInterval.Seconds() > 0 && lastInvocation.Add(userInterval).Unix() > message.Time().Unix() {
					return
				}
			}

			//Check if the user can invoke a list

			//Check if list was invoked
			//To do so we have a set of prefixes
			listPrefixes := []string{"@", "#", "!", "."}
			//We add the the possibleList every word that has one of the prefixes
			//TODO: Check with reges the lists
			possibleLists := make([]string, 0)
			//Then we iterate the prefixes, and for each one we see if there are possible lists
			for _, prefix := range listPrefixes {
				//If the message contains the prefix
				if strings.Contains(message.Text, prefix) {
					//we split all the words of the message (by the space), removing eventual commas and semicolons
					words := strings.Split(strings.Replace(strings.Replace(message.Text, ",", "", -1), ";", "", -1), " ")
					//And fore each word
					for _, word := range words {
						//If the word is longer than 1 char
						if len(word) > 1 {
							//And hase the prefix
							if word[0] == prefix[0] {
								//We add it to the list without the prefix
								possibleLists = append(possibleLists, strings.ToLower(word[1:]))
							} //fi prefix check
						} //fi world len
					} //end words loop
				} //fi message contains prefix
			} //end prefix loop

			//If there are not, there's no need to continue
			if len(possibleLists) < 1 {
				return
			}

			//Then we get all the lists in the group
			groupLists, err := ctx.Database.GetLists(message.Chat.ID)
			//And we check for errors
			if err != nil {
				replyDbMessageWithCloseButton(ctx, "generalError")
				return
			}
			//Then we create a slice that will contain all the invoked list
			lists := make([]database.List, 0)
			//We iterate throught all the possible lists
			for _, plist := range possibleLists {
				//And iterate throught all the lists present in the group
				for _, glist := range groupLists {
					//If the list name is the same as the possible list, there'sa match
					if plist == glist.Name {
						//So we append the current list to the lists to "call"
						lists = append(lists, glist)
						//Increment the number of lists contacted by the user if possible
						if userIntervalError == nil {
							ctx.Database.IncrementListsInvokedCount(message.From.ID, message.Chat.ID)
							//and updated the last time the user contacted a list
							ctx.Database.UpdateLastInvocation(message.From.ID, message.Chat.ID, message.Time())
						}
						//And we interrupt the iteration
						break
					} //fi check for list match
				} //end loop of the lists in the groyp
			} //end loop of possible lists in the message

			//if there was no match return
			if len(lists) < 1 {
				return
			}
			//We create a slice of the users who were contacted
			contactedUsers := make([]int64, 0)

			var listInterval time.Duration
			//We get the minimum interval a list should be called
			intervalString, intervalError := ctx.Database.GetSettingValue("listInterval", int(message.Chat.ID))
			if intervalError == nil {
				listInterval, intervalError = time.ParseDuration(intervalString)
			}
			//For each list
			for _, list := range lists {

				//If the list has a interval calling limit (wee see if the error is nil or not)
				if intervalError == nil {
					//If the minimum interval is not passed yet and the list interval is valid (greather than 0)
					//To check so, we add the minimum list timespan to the latest list invokation time
					//  and convert the number to an integer we can compare
					//If the integer is greater than the time of the message, the list cannot be called
					//  so we continue the loop

					if listInterval.Seconds() > 0 && list.LatestInvocation.Add(listInterval).Unix() > message.Time().Unix() {
						//If we go there, it means that the last time we called the list
						//  summed to the list interval is greater than the message time,
						//  so the list shouldn't be calle
						continue
					}
				}
				//Update the list invokation time
				ctx.Database.UpdateListLastInvokation(list.ID, message.Time())
				//We get the list of the subscribers
				subs, _ := ctx.Database.GetSubscribedUsers(list.ID)
				//For each subscriber
				for _, sub := range subs {
					//We set a flag to see if we already called to user (maybe he was in another list)
					found := false
					//And we iterate through the contacted users
					for _, cUse := range contactedUsers {
						//If the user was contacted
						if sub.UserID == cUse {
							//We set the flag to true
							found = true
							//And end the loop
							break
						} //fi check for contacted user
					} //end loop of contacted users

					//If the user was not contacted
					if !found {
						//Then we get info of the user to contact via the ID
						user, _ := ctx.Database.GetUser(int(sub.UserID))
						//Then we get the message to send to the user from the database, and replace the keywords
						//  such as {{categoryName}} with the real category name
						messageToSend := tba.NewMessage(sub.UserID, strings.Replace(strings.Replace(ctx.Database.GetBotStringValueOrDefaultNoError("tagNotification", user.Locale), "{{categoryName}}", list.Name, -1), "{{groupName}}", message.Chat.Title, -1))
						//If the message is in a supergroup, it's possible to get a link to the message
						if message.Chat.IsSuperGroup() {
							//If the group has a username
							if message.Chat.UserName != "" {
								//We generate the links, always by taking from the database the strings
								ikm1 := tba.NewInlineKeyboardButtonURL(ctx.Database.GetBotStringValueOrDefaultNoError("tagNotificationGroupLink", user.Locale), "t.me/"+message.Chat.UserName)
								ikm2 := tba.NewInlineKeyboardButtonURL(ctx.Database.GetBotStringValueOrDefaultNoError("tagNotificationMessageLink", user.Locale), "t.me/"+message.Chat.UserName+"/"+strconv.Itoa(message.MessageID))
								ikm3 := tba.NewInlineKeyboardButtonData(ctx.Database.GetBotStringValueOrDefaultNoError("tagNotificationTag", user.Locale), "tag-"+strconv.FormatInt(message.Chat.ID, 10)+"-"+strconv.Itoa(message.MessageID))
								ikl := []tba.InlineKeyboardButton{ikm1, ikm2, ikm3}
								ikm := tba.NewInlineKeyboardMarkup(ikl)
								//And add to the message the buttons
								messageToSend.ReplyMarkup = ikm
							} else { //The chat hasn't a nickname
								//We add the button to be tagged from the bot at the mesage
								ikm3 := tba.NewInlineKeyboardButtonData(ctx.Database.GetBotStringValueOrDefaultNoError("tagNotificationTag", user.Locale), "tag-"+strconv.FormatInt(message.Chat.ID, 10)+"-"+strconv.Itoa(message.MessageID))
								ikl := []tba.InlineKeyboardButton{ikm3}
								ikm := tba.NewInlineKeyboardMarkup(ikl)
								//And add to the message the buttons
								messageToSend.ReplyMarkup = ikm
							}
						} else { //The message is not a supergroup
							ikm3 := tba.NewInlineKeyboardButtonData(ctx.Database.GetBotStringValueOrDefaultNoError("tagNotificationTag", user.Locale), "tag-"+strconv.FormatInt(message.Chat.ID, 10)+"-"+strconv.Itoa(message.MessageID))
							ikl := []tba.InlineKeyboardButton{ikm3}
							ikm := tba.NewInlineKeyboardMarkup(ikl)
							messageToSend.ReplyMarkup = ikm
						}
						//We then send the message to the user
						ctx.Bot.Send(messageToSend)
						//And add the ID of the user to the slice of the contacted users
						contactedUsers = append(contactedUsers, sub.UserID)
					} //fi user found
				} //end subscribed users loop
			} //end lists loop
			//Notify the user that the lists were called successfully
			if len(contactedUsers) > 0 {
				replyDbMessage(ctx, "listNotificationSuccessMessage")
			}
		} //fi messageIngroup
	} //fi message command
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

func replyDbMessage(ctx *Context, keyString string) {
	messageBody := ctx.Database.GetBotStringValueOrDefaultNoError(keyString, ctx.Update.Message.From.LanguageCode)
	messageToSend := tba.NewMessage(ctx.Update.Message.Chat.ID, messageBody)
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
