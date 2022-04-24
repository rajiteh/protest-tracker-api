package bot

import (
	"context"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/reliefeffortslk/protest-tracker-api/pkg/bot/commands"
	"go.uber.org/multierr"
	tb "gopkg.in/tucnak/telebot.v3"
	"gopkg.in/tucnak/telebot.v3/middleware"
)

var log = logrus.New()

var bot *tb.Bot

type BotService struct {
}

func (bs *BotService) Serve(_ context.Context) error {
	log.Info("Creating BotService.")
	createBot()
	return nil
}

// struct ConverstaionSt
func createBot() {
	var err error
	pref := tb.Settings{
		Token:  os.Getenv("TELEGRAM_APITOKEN"),
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err = tb.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	// Load bot Admins
	botAdmins := []int64{}
	for _, adminStr := range strings.Split(os.Getenv("TELEGRAM_BOT_ADMINS"), ",") {
		if adminStr == "" {
			continue
		}
		adminIdInt, err := strconv.ParseInt(adminStr, 10, 64)
		if err != nil {
			log.Fatalf("malformed admin id in env: %s (%v)", adminStr, err)
			return
		}
		botAdmins = append(botAdmins, adminIdInt)
	}

	bot.Use(middleware.Logger())

	var botCommands []commands.CommandHandler
	botCommands = append(botCommands,
		commands.NewSubscribe(bot).CommandHandler,
		commands.NewIngest(bot, botAdmins).CommandHandler,
		commands.NewNear(bot).CommandHandler,
	)

	setupGlobalHandlers(botCommands)

	bot.Start()

}

func setupGlobalHandlers(botCommands []commands.CommandHandler) {
	// Setup global handlers for location and text
	resolveHandlers := func(ctx tb.Context, handlerResolver func(ch commands.CommandHandler) []tb.HandlerFunc) error {
		var errors []error
		for _, botCommand := range botCommands {
			handlers := handlerResolver(botCommand)
			for _, handler := range handlers {
				if err := handler(ctx); err != nil {
					errors = append(errors, err)
				}
			}
		}
		return multierr.Combine(errors...)
	}

	commonEvents := []string{
		tb.OnLocation,
		tb.OnText,
	}

	for _, event := range commonEvents {
		bot.Handle(event, func(ctx tb.Context) error {
			return resolveHandlers(ctx, func(ch commands.CommandHandler) []tb.HandlerFunc {
				return ch.ResponseHandlers
			})
		})

	}
}
