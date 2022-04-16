//go:generate go-localize -input localizations_src -output localizations

package commands

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"

	tb "gopkg.in/tucnak/telebot.v3"

	"github.com/rusq/tbcomctl/v4"
)

type SubscribeCommand struct {
	CommandHandler
}

func NewSubscribe(bot *tb.Bot) *SubscribeCommand {
	sh := SubscribeCommand{}
	sh.bot = bot

	sh.setupHandlers()

	return &sh
}

func (sh *SubscribeCommand) setupHandlers() {
	log.Info("Setting up handlers for subscribe")

	locationPinInput := tbcomctl.NewInputText("location_pin", "Send a location pin from telegram app:", processLocationInput(sh.bot), tbcomctl.IOptValueResolver(func(m *tb.Message) (string, error) {
		location, err := json.Marshal(m.Location)
		if err != nil {
			return "", err
		}
		return string(location), nil
	}))

	radiusRangeInput := tbcomctl.NewPicklist("radius_pick",
		tbcomctl.NewStaticTVC(
			"Select the distance from this point for this subscription:",
			[]string{
				"1",
				"3",
				"5",
				"10",
			},
			processRadiusInput(sh.bot),
		),
		tbcomctl.PickOptOverwrite(false),
		tbcomctl.PickOptBtnPattern([]uint{4}),
	)

	tvc := tbcomctl.NewStaticTVC("", nil, nil)
	tvc.TextFn = finalizeSubscribeForm(sh.bot)
	finalizerInput := tbcomctl.NewMessage("finalize", tvc)

	form := tbcomctl.NewForm(locationPinInput, radiusRangeInput, finalizerInput).SetOverwrite(false).SetRemoveButtons(true)

	sh.bot.Handle("/subscribe", form.Handler)
	sh.ResponseHandlers = append(sh.ResponseHandlers, form.OnTextMiddleware(func(ctx tb.Context) error {
		return nil
	}))
}

func processLocationInput(b *tb.Bot) func(ctx context.Context, c tb.Context) error {
	return func(ctx context.Context, c tb.Context) error {
		log.Info("Inside process location input")
		loc := c.Message().Location

		if loc == nil {
			return tbcomctl.NewInputError("expected a location pin from telegram, please try again.")
		}
		return nil
	}
}

func processRadiusInput(b *tb.Bot) func(ctx context.Context, c tb.Context) error {
	return func(ctx context.Context, c tb.Context) error {
		log.Info("Inside process radius input")
		radiusStr := c.Data()
		if _, err := strconv.Atoi(radiusStr); err != nil {
			log.Error("Something went wrong trying to get radius", err.Error())
			return tbcomctl.NewInputError("expected a valid radius value from telegram, please try again.")
		}

		return nil
	}
}

func finalizeSubscribeForm(b *tb.Bot) func(ctx context.Context, c tb.Context) (string, error) {
	return func(ctx context.Context, c tb.Context) (string, error) {
		log.Info("Finalizing form...")
		if ctrl, ok := tbcomctl.ControllerFromCtx(ctx); ok {
			form := ctrl.Form()
			data := form.Data(c.Sender())
			log.Println("form values so far: ", data)
		} else {
			return "There was a problem", errors.New("something went wrong trying to process the form")
		}

		return "Success", nil
	}
}
