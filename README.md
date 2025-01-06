# Golang bindings for the Telegram Bot API

This package is a fork to add minimal new features, and to bring support for the
latest Telegram Bot API.

[![Go Reference](https://pkg.go.dev/badge/github.com/go-telegram-bot-api/telegram-bot-api/v5.svg)](https://pkg.go.dev/github.com/go-telegram-bot-api/telegram-bot-api/v5)
[![Test](https://github.com/go-telegram-bot-api/telegram-bot-api/actions/workflows/test.yml/badge.svg)](https://github.com/go-telegram-bot-api/telegram-bot-api/actions/workflows/test.yml)

All methods are fairly self-explanatory, and reading the [godoc](https://pkg.go.dev/github.com/go-telegram-bot-api/telegram-bot-api/v5) page should
explain everything. If something isn't clear, open an issue or submit
a pull request.

There are more tutorials and high-level information on the website, [go-telegram-bot-api.dev](https://go-telegram-bot-api.dev).

The scope of this project remains close to the original project, but adds
a simple command dispatching model that makes it easy to get your bot up
and running quickly without having to implement the same boilerplate code
each time.

Use `github.com/semog/telegram-bot-api` for the latest version.

Join [the original development group](https://t.me/go_telegram_bot_api) if
you want to ask questions or discuss development. Remember that this is a branch
from the original development version.

## Example

First, ensure the library is installed and up to date by running
`go get -u github.com/semog/telegram-bot-api/v5`.

This sample shows a main() function that connects to the bot,
and then starts the command listener loop. This is all that is required in the
main() function. The RunBot() function will handle running your bot and dispatch
messages via the event handlers. It will call the OnInitialize
handler once upon startup. It will call the OnDispose handler once on shutdown.
It will then call the OnMessage handler whenever a new message
is received. The OnMessage handler should parse out any command messages that
the bot has registered.

```go
package main

import (
	"log"

	tg "github.com/semog/telegram-bot-api/v5"
)

func main() {
	bot, err := tg.NewBotAPI("MyAwesomeBotToken")
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Connected to Bot: %s (%s)", bot.Self.FirstName, bot.Self.UserName)
	tg.RunBot(bot, mybothandlers)
}

// ------------ The following can be placed into a separate source file -----------------
// mybothandlers maps the dispatch function handlers.
var mybothandlers = tg.BotEventHandlers{
	OnInitialize: mybotOnInitialize,
	OnDispose:    mybotOnDispose,
	OnCommand:    mybotOnCommand,
}

// Initialize global data, and the bot commands with optional botname attached.
func mybotOnInitialize(bot *tg.BotAPI) bool {
	botname := bot.Self.UserName
	return true
}

func mybotOnDispose(bot *tg.BotAPI) {
	// Do any cleanup of external resources.
}

// mybotOnCommand is the main handler that receives command messages
func mybotOnCommand(bot *tg.BotAPI, cmd string, msg *tg.Message) bool {
	log.Printf("Command From: User %s %s (%s): %s - %s",
		msg.From.FirstName, msg.From.LastName, msg.From.UserName, cmd, msg.Text)
	switch {
	case "action":
		doAction(bot, msg)
	case "quit":
		doQuit(bot, msg)
		return false
	}
	return true
}

func doTextReply(bot *tg.BotAPI, msg *tg.Message) {
	log.Printf("[%s] %s", msg.From.UserName, msg.Text)
	replymsg := tg.NewMessage(msg.Chat.ID, msg.Text)
	replymsg.ReplyToMessageID = msg.MessageID
	bot.Send(replymsg)
}

func doAction(bot *tg.BotAPI, msg *tg.Message) {
	// Do /action command
}

func doQuit(bot *tg.BotAPI, msg *tg.Message) {
	// Do /quit action
}
```

If you need to use webhooks (if you wish to run on Google App Engine),
you may use a slightly different method.

```go
package main

import (
	"log"
	"net/http"
	tg "github.com/semog/telegram-bot-api/v5"
)

func main() {
	bot, err := tg.NewBotAPI("MyAwesomeBotToken")
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	wh, _ := tg.NewWebhookWithCert("https://www.example.com:8443/"+bot.Token, "cert.pem")

	_, err = bot.Request(wh)
	if err != nil {
		log.Fatal(err)
	}

	info, err := bot.GetWebhookInfo()
	if err != nil {
		log.Fatal(err)
	}

	if info.LastErrorDate != 0 {
		log.Printf("Telegram callback failed: %s", info.LastErrorMessage)
	}

	updates := bot.ListenForWebhook("/" + bot.Token)
	go http.ListenAndServeTLS("0.0.0.0:8443", "cert.pem", "key.pem", nil)

	for update := range updates {
		log.Printf("%+v\n", update)
	}
}
```

If you need, you may generate a self-signed certificate, as this requires
HTTPS / TLS. The above example tells Telegram that this is your
certificate and that it should be trusted, even though it is not
properly signed.

    openssl req -x509 -newkey rsa:2048 -keyout key.pem -out cert.pem -days 3560 -subj "//O=Org\CN=Test" -nodes

Now that [Let's Encrypt](https://letsencrypt.org) is available,
you may wish to generate your free TLS certificate there.
