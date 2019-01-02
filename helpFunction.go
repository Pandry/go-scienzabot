package main

//SendHelpMessage is the method thet replies to the
//help requests sending the help message
func (ctx *Context) SendHelpMessage() {
	//Shall reply to the user sending the help message
	//TODO: put help message in database
	//TODO: Grab help text from database
	//TODO: reply only to user in allowed group, private chat
	/*
		helpMessage, err := ctx.Database.GetSettingValue("HelpMessage",
			ctx.Update.Update.Message.From.LanguageCode)
		if err != nil {
			log.Panic("Something's wrong here!\nCannot get help message from database!", err)
		}

		ctx.Update.ReplyToUpdate(helpMessage)*/

}
