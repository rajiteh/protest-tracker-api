package commands

import (
	"fmt"

	"github.com/sirupsen/logrus"
	tb "gopkg.in/tucnak/telebot.v3"
)

var log = logrus.New()

type CommandHandler struct {
	bot              *tb.Bot
	ResponseHandlers []tb.HandlerFunc
}

type SecuredCommandHandler struct {
	CommandHandler
	AllowedUsers []int64
}

type NotAllowedError struct {
	Reason string
}

func (e *NotAllowedError) Error() string {
	return e.Reason
}

func (sch *SecuredCommandHandler) VerifiedHandler(handler tb.HandlerFunc) func(c tb.Context) error {
	return func(c tb.Context) error {
		senderId := c.Sender().ID
		if ok := sliceContainsInt64(sch.AllowedUsers, senderId); ok {
			return handler(c)
		}

		return &NotAllowedError{
			Reason: fmt.Sprintf("User ID %d not allowed to perform this action", senderId),
		}
	}
}

func sliceContainsInt64(s []int64, e int64) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
