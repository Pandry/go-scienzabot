package main

import (
	"flag"
	"log"
	"os"
	"strconv"
	"time"

	"scienzabot/database"
	"scienzabot/embtypes"

	tba "github.com/go-telegram-bot-api/telegram-bot-api"
)

var verboseMode *bool

func main() {

	//The APIToken variable contains the token, taken by default
	//from the environment variable, if not explicited via start argument
	var APIToken string
	var dbPath string

	//Debug is an argument passed directly to the Telegram Bot API
	//Enabilng it, will cause the library to print to STDOUT incoming updates and replies
	verboseMode = flag.Bool("v", false, "Verbose, tells how the messages get elaborated as they \"flow\".")

	//Debug is an argument passed directly to the Telegram Bot API
	//Enabilng it, will cause the library to print to STDOUT incoming updates and replies
	vvPtr := flag.Bool("vv", false, "Sets the Telegram API in debug mode and the bot in verbose mode- outputs the data it sends and receives.")

	//The database path is the path of the SQLite3 database, that the bot will use as a base
	flag.StringVar(&dbPath, "database", "database.sqlite3", "The default databse path")

	//the API token is assigned by default by reading the "TELEGRAM_TOKEN" environment variable
	// if not explicited as an agument
	flag.StringVar(&APIToken, "token", os.Getenv("TELEGRAM_TOKEN"), "The token of the bot. By default the value is taken from the TELEGRAM_TOKEN environment variable.")

	//Parses the flags to read the values
	flag.Parse()

	if *vvPtr == true {
		*verboseMode = true
	}

	//Checking if the bot was submitted
	if APIToken == "" {
		log.Panic("Cannot start the bot!\nThe token env variable (TELEGRAM_TOKEN) is not setted \n")
	}

	//Instatiating the bot with the TELEGRAM_TOEKN environment variable value
	//It needs to be setted.
	//bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_TOEKN"))
	oldBot, err := tba.NewBotAPI(APIToken)
	if err != nil {
		//If there's an error, display it
		log.Panic("Cannot start the bot!\nProbably there's an error with the token!\nError details: ", err, "\n")
	}

	//Now, we want to convert the standard bot API instance, to our one.
	//The reson is because we can assign a receiver to our instance
	//Doing so the code looks cleaner
	//such procedure is called embedding, and is similiar to the
	// OOP's concept of ereditariety
	//
	//https://golang.org/doc/effective_go.html#embedding
	bot := &embtypes.Tgbotapi{BotAPI: oldBot}

	//Initializing the database connection and running the statrtup queries
	db, err := database.InitDatabaseConnection(dbPath)
	if err != nil {
		log.Panic("Impossible to establish a connection with the database !\nError details: ", err, "\n")
	}
	log.Println("Database connection open on SQLite3 database", dbPath)

	//And defer the closing of the database
	defer db.Close()

	//initializing the context.
	//See the context comment for more info
	messageQueue := make(chan *tba.MessageConfig, 10000)
	ctx := Context{bot, nil, db, messageQueue}

	go func(c Context) {
		//Limit is 30 messages per second
		//TODO: Implement limits against same user
		for {
			select {
			case m := <-c.SendQueue:
				log.Printf("Sending message from queue to list. Sending to " + strconv.FormatInt(m.ChatID, 10))
				ctx.Bot.Send(m)
				break
			}
			time.Sleep(time.Second / 20)
		}
	}(ctx)

	//Assinging the debug variable to the lib.
	//This will tell if the library will print the updates it receives and send
	bot.Debug = *vvPtr
	//Notifying the successful connection to the telegram servers
	log.Printf("Bot authorized on account @%s", bot.Self.UserName)
	log.Println("Setting highly verbose mode to", *vvPtr)

	//creates a long possling stream with 0 as starting offset
	u := tba.NewUpdate(0)

	//Sets the timeout for the long polling
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)
	for update := range updates {
		ctx.Update = &embtypes.Tgupdate{Update: &update}
		ctx.route()
	}
}
