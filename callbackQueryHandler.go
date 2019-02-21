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
		user    database.User
		isAdmin bool
		err     error
	)

	if userExists := ctx.Database.UserExists(message.From.ID); userExists {
		user, err = ctx.Database.GetUser(message.From.ID)
		isAdmin = utils.HasPermission(int(user.Permissions), consts.UserPermissionAdmin)
	}

	switch {
	case message.Data != "" && strings.Contains(message.Data, "-"):
		switch args := strings.Split(message.Data, "-"); args[0] {
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
