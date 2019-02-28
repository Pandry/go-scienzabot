package main

import (
	"log"
)

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
			textMessageRoute(ctx)
			break

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

	case ctx.Update.EditedMessage != nil:
		//Edited text message
		editedTxtMessageRoute(ctx)
		break

	case ctx.Update.CallbackQuery != nil:
		//Callback query
		callbackQueryRoute(ctx)
		break

	case ctx.Update.InlineQuery != nil:
		//Inline query
		break

	case ctx.Update.ChannelPost != nil:
		//Channel post
		break

	case ctx.Update.EditedChannelPost != nil:
		//Edited channel post
		break

	case ctx.Update.PreCheckoutQuery != nil:
		//Pre checkoput query - useless for now
		break

	case ctx.Update.ShippingQuery != nil:
		//Pre shipping query
		break

	case ctx.Update.ChosenInlineResult != nil:
		//Chosen inline result -> Chosen inline element?
		break

	} //Message type swtich
}
