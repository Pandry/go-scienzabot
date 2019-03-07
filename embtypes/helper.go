package embtypes

import (
	"scienzabot/consts"

	tba "github.com/go-telegram-bot-api/telegram-bot-api"
)

// The helpers.go file contains helper functions that extends the default library
//TODO: needs to be tested, written offline

//SendLongMessage is used to send a message longer than the telegram limit (splitting it in more messages)
func (t *Tgbotapi) SendLongMessage(chatID int64, message string, ReplyToMessageID int, replyMarkup interface{}, parseMode string) {
	//TODO: make messageSplitters a constant
	messageSplitters := []string{"\n\n\n", "\n\n", "\n", " "}
	sentChars := 0
	//Cycle all the message string
	for len(message) > sentChars {
		messageToSend := ""
		//set the base index to the offset
		if consts.MaximumMessageLength+consts.MaximumMessageLengthMargin < len(message[sentChars:]) {
			splitterOffset := -1
			splitterLen := 0
			//The left chars are still too much and need to be splitted
			//For each "message splitter"
			for _, messageSplitter := range messageSplitters {
				//So we set i to the max len the message can have and start seaching backward for
				// a char we can split the message at.
				//We seach until we find the char we are searching for or we search too in deep and
				//  go beyond the margins
				for i := consts.MaximumMessageLength + consts.MaximumMessageLengthMargin + sentChars; splitterOffset == -1 && i > sentChars+consts.MaximumMessageLength-consts.MaximumMessageLengthMargin; i-- {
					//Then we check for each "margin char" if it corresponds to be a valid
					//  point where we can split the message
					if message[sentChars+i-len(messageSplitter):sentChars+i] == messageSplitter {
						//We found where the message can be splitted
						splitterOffset = sentChars + i - len(messageSplitter)
						splitterLen = len(messageSplitter)
						//If splitter offset was found, we can exit the loop
						break
					}
				}

			}
			//If splitter offset is still -1, a "hard split" needs to be done
			if splitterOffset == -1 {
				messageToSend = message[sentChars:consts.MaximumMessageLengthMargin]
			} else {
				//A splitter was found
				messageToSend = message[sentChars : sentChars+splitterOffset+splitterLen]
			}

		} else {
			//The chars left to send are less than the max message len and does not need to be splitted
			messageToSend = message[sentChars:]
		}
		//Finally we update the sent chars for the next iteration
		sentChars += len(messageToSend)
		//send the message
		msg := tba.NewMessage(chatID, messageToSend)
		msg.ReplyToMessageID = ReplyToMessageID
		msg.ReplyMarkup = replyMarkup
		msg.ParseMode = parseMode
		t.Send(msg)
	}
}

/*
func (u *Tgupdate) ReplyToUpdate(message string) {

}

*/
