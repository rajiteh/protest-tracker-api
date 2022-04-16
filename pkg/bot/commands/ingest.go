package commands

import (
	"fmt"

	"github.com/reliefeffortslk/protest-tracker-api/pkg/ingestor"
	tb "gopkg.in/tucnak/telebot.v3"
)

type IngestCommand struct {
	SecuredCommandHandler
}

func NewIngest(bot *tb.Bot, allowedUsers []int64) *IngestCommand {
	ic := IngestCommand{}
	ic.bot = bot
	ic.AllowedUsers = allowedUsers

	ic.setupHandlers()

	return &ic
}

func (ic *IngestCommand) setupHandlers() {
	log.Info("Setting up hanlders for ingest")
	ic.bot.Handle("/ingest", ic.VerifiedHandler(func(c tb.Context) error {
		c.Send("Starting ingestion")
		defer c.Send("Finished ingestion")

		if count, err := ingestor.IngestFromAll(); err != nil {
			c.Send(fmt.Sprintf("Ingestion failed: %s", err.Error()))
			return err
		} else {
			c.Send(fmt.Sprintf("Inserted %d rows", count))
		}
		return nil
	}))
}
