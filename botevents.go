package tgbotapi

// BotInitializeEvent called once before any other events so the bot can initialize
// any state configurations before it starts running.
type BotInitializeEvent func(bot *BotAPI) bool

// BotDisposeEvent called once before the bot is shutdown so it can clean up any
// resources and save any state.
type BotDisposeEvent func(bot *BotAPI)

// BotUpdateEvent called when the bot receives an Update object, and wants to process the raw Update.
type BotUpdateEvent func(bot *BotAPI, msg *Update) bool

// BotMessageEvent called when the bot receives a Message update.
type BotMessageEvent func(bot *BotAPI, msg *Message) bool

// BotCommandEvent called when the bot receives a Command.
type BotCommandEvent func(bot *BotAPI, cmd string, msg *Message) bool

// BotInlineQueryEvent called when the bot receives an InlineQuery update.
type BotInlineQueryEvent func(bot *BotAPI, query *InlineQuery) bool

// BotChosenInlineResultEvent called when the bot receives a ChosenInlineResult update.
type BotChosenInlineResultEvent func(bot *BotAPI, result *ChosenInlineResult) bool

// BotCallbackQueryEvent called when the bot receives a CallbackQuery update.
type BotCallbackQueryEvent func(bot *BotAPI, query *CallbackQuery) bool

// BotShippingQueryEvent called when the bot receives a ShippQuery update.
type BotShippingQueryEvent func(bot *BotAPI, query *ShippingQuery) bool

// BotPreCheckoutQueryEvent called when the bot received a PreCheckoutQuery update.
type BotPreCheckoutQueryEvent func(bot *BotAPI, query *PreCheckoutQuery) bool

// BotEventHandlers contains function pointers for handling the different events.
type BotEventHandlers struct {
	OnInitialize         BotInitializeEvent
	OnDispose            BotDisposeEvent
	OnUpdate             BotUpdateEvent
	OnCommand            BotCommandEvent
	OnMessage            BotMessageEvent
	OnEditedMessage      BotMessageEvent
	OnChannelPost        BotMessageEvent
	OnEditedChannelPost  BotMessageEvent
	OnInlineQuery        BotInlineQueryEvent
	OnChosenInlineResult BotChosenInlineResultEvent
	OnCallbackQuery      BotCallbackQueryEvent
	OnShippingQuery      BotShippingQueryEvent
	OnPreCheckoutQuery   BotPreCheckoutQueryEvent
}

// RunBot is the main loop that runs the bot and dispatches messages to registered event
// handlers. All handlers are optional.
func RunBot(bot *BotAPI, handler BotEventHandlers) {
	u := NewUpdate(0)
	u.Timeout = 60
	updates, _ := bot.GetUpdatesChan(u)

	defer func() {
		log.Infof("Shutting down %s", bot.Self.UserName)
		// Must have initialize function in order to call dispose function.
		if handler.OnInitialize != nil && handler.OnDispose != nil {
			handler.OnDispose(bot)
		}
	}()

	if handler.OnInitialize != nil && !handler.OnInitialize(bot) {
		return
	}

	var keepgoing = true
	for update := range updates {
		// Support generic handling of the raw update message.
		if handler.OnUpdate != nil && !handler.OnUpdate(bot, &update) {
			break
		}

		// Call specific event handlers
		switch {
		case update.Message != nil && update.Message.IsCommand() && handler.OnCommand != nil:
			keepgoing = handler.OnCommand(bot, update.Message.Command(), update.Message)
		case update.Message != nil && handler.OnMessage != nil:
			keepgoing = handler.OnMessage(bot, update.Message)
		case update.EditedMessage != nil && handler.OnEditedMessage != nil:
			keepgoing = handler.OnEditedMessage(bot, update.EditedMessage)
		case update.ChannelPost != nil && handler.OnChannelPost != nil:
			keepgoing = handler.OnChannelPost(bot, update.ChannelPost)
		case update.EditedChannelPost != nil && handler.OnEditedChannelPost != nil:
			keepgoing = handler.OnEditedChannelPost(bot, update.EditedChannelPost)
		case update.InlineQuery != nil && handler.OnInlineQuery != nil:
			keepgoing = handler.OnInlineQuery(bot, update.InlineQuery)
		case update.ChosenInlineResult != nil && handler.OnChosenInlineResult != nil:
			keepgoing = handler.OnChosenInlineResult(bot, update.ChosenInlineResult)
		case update.CallbackQuery != nil && handler.OnCallbackQuery != nil:
			keepgoing = handler.OnCallbackQuery(bot, update.CallbackQuery)
		case update.ShippingQuery != nil && handler.OnShippingQuery != nil:
			keepgoing = handler.OnShippingQuery(bot, update.ShippingQuery)
		case update.PreCheckoutQuery != nil && handler.OnPreCheckoutQuery != nil:
			keepgoing = handler.OnPreCheckoutQuery(bot, update.PreCheckoutQuery)
		default:
			if bot.Debug {
				log.Infof("Unhandled Bot Event: %v", update)
			}
			keepgoing = true
		}

		if !keepgoing {
			break
		}
	}
}
