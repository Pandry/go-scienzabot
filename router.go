package main

// The router.go file is a file that "routes" a message received to the designed
//	received (e.g. A normal text message should be "routed" to the text message "route"
//	to be elaborated correctly)

func (ctx *Context) route() {

	switch {
	//Swtich the message based on its type

	//Return message if there's no update
	case ctx.Update == nil:
		return

	case ctx.Update.Message != nil:
		//General message
		message := ctx.Update.Message

		switch {
		case message.Text != "":
			logMessage(0, "[", message.MessageID, "]", " Got text message from [", message.From.ID, "]")
			textMessageRoute(ctx)
			break

		case message.Sticker != nil:
			//Sticker
			logMessage(0, "[", message.MessageID, "]", " Got sticker message from [", message.From.ID, "]")
			break

		case message.Photo != nil:
			logMessage(0, "[", message.MessageID, "]", " Got photo from [", message.From.ID, "]")
			//Photo
			break

		case message.Voice != nil:
			logMessage(0, "[", message.MessageID, "]", " Got audio note from [", message.From.ID, "]")
			//Voice audio file

			break

		case message.Document != nil:
			//Document
			logMessage(0, "[", message.MessageID, "]", " Got file from [", message.From.ID, "]")

			break

		case message.Location != nil:
			logMessage(0, "[", message.MessageID, "]", " Got location from [", message.From.ID, "]")

			break

		case message.Contact != nil:
			logMessage(0, "[", message.MessageID, "]", " Got contact from [", message.From.ID, "]")

			break

		case message.NewChatMembers != nil:
			//New user(s)
			for _, user := range *ctx.Update.Message.NewChatMembers {
				logMessage(0, "[", message.MessageID, "]", " User joined (", user.UserName, ")[", message.From.ID, "]")
			}
			if ctx.Update.Message.GroupChatCreated {
				logMessage(0, "[", message.MessageID, "]", " Group created")
			}

			//User joined
			break

		case message.VideoNote != nil:
			//Video circolare
			logMessage(0, "[", message.MessageID, "]", " Got video note from [", message.From.ID, "]")

			break

		case message.Video != nil:
			logMessage(0, "[", message.MessageID, "]", " Got video from [", message.From.ID, "]")

			break

		case message.Venue != nil:
			//NO IDEA     <--- non capisco (Bhez)
			logMessage(0, "[", message.MessageID, "]", " Got venue (to understand what this is) from [", message.From.ID, "]")
			break

		case message.LeftChatMember != nil:
			//User removed (could be the bot)
			logMessage(0, "[", message.MessageID, "]", " A user left or was kicked from the group [", message.From.ID, "]")
			break

		case message.PinnedMessage != nil:
			logMessage(0, "[", message.MessageID, "]", " A message was pinned [", message.Chat.ID, "]")

			break

		case message.NewChatPhoto != nil:
			logMessage(0, "[", message.MessageID, "]", " The chat photo was updated [", message.Chat.ID, "]")

			break

		case message.NewChatTitle != "":
			logMessage(0, "[", message.MessageID, "]", " The chat title was updated [", message.Chat.ID, "]")

			break

		case message.MigrateToChatID != 0:
			logMessage(0, "[", message.MessageID, "]", " The chat migrated to [", message.MigrateToChatID, "]")

			break
		}

	case ctx.Update.EditedMessage != nil:
		//Edited text message
		logMessage(0, "[", ctx.Update.EditedMessage.MessageID, "]", " Got edit message from [", ctx.Update.EditedMessage.From.ID, "]")
		editedTxtMessageRoute(ctx)
		break

	case ctx.Update.CallbackQuery != nil:
		//Callback query
		logMessage(0, "[", ctx.Update.CallbackQuery.ID, "]", " Got callback query from [", ctx.Update.CallbackQuery.From.ID, "]")
		callbackQueryRoute(ctx)
		break

	case ctx.Update.InlineQuery != nil:
		//Inline query
		logMessage(0, "[", ctx.Update.InlineQuery.ID, "]", " Got inline query from [", ctx.Update.InlineQuery.From.ID, "]")

		break

	case ctx.Update.ChannelPost != nil:
		//Channel post
		logMessage(0, "[", ctx.Update.ChannelPost.MessageID, "]", " Got channel post from [", ctx.Update.ChannelPost.MessageID, "]")

		break

	case ctx.Update.EditedChannelPost != nil:
		//Edited channel post
		logMessage(0, "[", ctx.Update.EditedChannelPost.MessageID, "]", " Got edit of channel post from [", ctx.Update.EditedChannelPost.MessageID, "]")
		break

	case ctx.Update.PreCheckoutQuery != nil:
		//Pre checkoput query - useless for now
		logMessage(0, "[", ctx.Update.PreCheckoutQuery.ID, "]", " Got pre-checkout query from [", ctx.Update.PreCheckoutQuery.From.ID, "]")

		break

	case ctx.Update.ShippingQuery != nil:
		//Pre shipping query
		logMessage(0, "[", ctx.Update.ShippingQuery.ID, "]", " Got shipping query from [", ctx.Update.ShippingQuery.From.ID, "]")

		break

	case ctx.Update.ChosenInlineResult != nil:
		//Chosen inline result -> Chosen inline element?
		logMessage(0, "[", ctx.Update.ChosenInlineResult.InlineMessageID, "]", " Got chosen inline result from [", ctx.Update.ChosenInlineResult.From.ID, "]")
		break

	} //Message type swtich
}
