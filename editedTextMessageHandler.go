package main

// The editedTextMessageHandler.go file contains the code that gets executed when a
//	message is edited

func editedTxtMessageRoute(ctx *Context) {
	message := ctx.Update.EditedMessage

	if ctx.Database.UserExists(message.From.ID) {
		ctx.Database.UpdateUserLastSeen(message.From.ID, message.Time())
	}
}
