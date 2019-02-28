package embtypes

import (
	tba "github.com/go-telegram-bot-api/telegram-bot-api"
)

// The embtypes (short for embedded types) is the package used to extend the library types.
// This because we can make a receiver using an embedded type, and without editing the
//	library we are using

//https://golang.org/doc/effective_go.html#embedding

//Tgbotapi is the abstractoin of botapi to do things like add func receivees
//TODO write real description
type Tgbotapi struct {
	*tba.BotAPI
}

//Tgupdate is the abstractoin of botapi to do things like add func receivees
//TODO write real description
type Tgupdate struct {
	*tba.Update
}
