# Golang bindings for the Telegram Bot API

This package is a fork to add minimal new features, and to bring support for the
latest Telegram Bot API. The original project was not being updated to the latest
Telegram Bot API versions.

The scope of this project remains close to the original project, but adds
a simple command dispatching model that makes it easy to get your bot up
and running quickly without having to implement the same boilerplate code
each time.

Use `github.com/semog/telegram-bot-api` for the latest version.

Join [the original development group](https://t.me/go_telegram_bot_api) if
you want to ask questions or discuss development. Remember that this is a branch
from the original development version.

## Example

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
	"regexp"

	tg "github.com/semog/telegram-bot-api"
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
	OnMessage:    mybotOnMessage,
}

var actionCmd,
	quitCmd *regexp.Regexp
var basicCmd = regexp.MustCompile(`(?i)^/`)

// Initialize global data, and the bot commands with optional botname attached.
func mybotOnInitialize(bot *tg.BotAPI) bool {
	botname := bot.Self.UserName
	actionCmd = regexp.MustCompile(fmt.Sprintf(`(?i)^/?action(@%s)?`, botname))
	quitCmd = regexp.MustCompile(fmt.Sprintf(`(?i)^/?quit(@%s)?`, botname))
	return true
}

func mybotOnDispose(bot *tg.BotAPI) {
	// Do any cleanup of external resources.
}

// mybotOnMessage is the main handler that receives text messages
// from the users. Parse the text to look for bot commands.
func mybotOnMessage(bot *tg.BotAPI, msg *tg.Message) bool {
	log.Printf("MsgFrom: User %s %s (%s): %s",
		msg.From.FirstName, msg.From.LastName, msg.From.UserName, msg.Text)
	switch {
	case !basicCmd.MatchString(msg.Text) && len(msg.Text) > 0:
		// This is a non-command message sent to the bot.
		doTextReply(bot, msg)
	case actionCmd.MatchString(msg.Text):
		doAction(bot, msg)
	case quitCmd.MatchString(msg.Text):
		doQuit(bot, msg)
		return false
	}
	return true
}

func doTextReply(bot *tg.BotAPI, msg *tg.Message) {
	log.Printf("[%s] %s", msg.From.UserName, msg.Text)
	replymsg := tgbotapi.NewMessage(msg.Chat.ID, msg.Text)
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
	"gopkg.in/telegram-bot-api.v4"
	"log"
	"net/http"
)

func main() {
	bot, err := tgbotapi.NewBotAPI("MyAwesomeBotToken")
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	_, err = bot.SetWebhook(tgbotapi.NewWebhookWithCert("https://www.google.com:8443/"+bot.Token, "cert.pem"))
	if err != nil {
		log.Fatal(err)
	}

	updates := bot.ListenForWebhook("/" + bot.Token)
	go http.ListenAndServeTLS("0.0.0.0:8443", "cert.pem", "key.pem", nil)

	for update := range updates {
		log.Printf("%+v\n", update)
	}
}
```

If you need, you may generate a self signed certficate, as this requires
HTTPS / TLS. The above example tells Telegram that this is your
certificate and that it should be trusted, even though it is not
properly signed.

    openssl req -x509 -newkey rsa:2048 -keyout key.pem -out cert.pem -days 3560 -subj "//O=Org\CN=Test" -nodes

Now that [Let's Encrypt](https://letsencrypt.org) has entered public beta,
you may wish to generate your free TLS certificate there.
