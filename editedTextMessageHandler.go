package main

func editedTxtMessageRoute(ctx *Context) {
	message := ctx.Update.EditedMessage

	if ctx.Database.UserExists(message.From.ID) {
		ctx.Database.UpdateUserLastSeen(message.From.ID, message.Time())
	}
}
