package main

import (
	"os"

	irc "github.com/akrennmair/joe-irc-adapter"
	"github.com/go-joe/joe"
	"github.com/pkg/errors"
)

type ExampleBot struct {
	*joe.Bot
}

func main() {
	cfg := irc.Config{
		Address: os.Getenv("IRC_SERVER"),
		Nick:    os.Getenv("IRC_NICK"),
		Name:    "User " + os.Getenv("IRC_NICK"),
		Channel: os.Getenv("IRC_CHANNEL"),
	}

	b := &ExampleBot{
		Bot: joe.New("my-irc-bot",
			irc.Adapter(cfg),
		),
	}

	b.Respond("remember (.+) is (.+)", b.Remember)
	b.Respond("what is (.+)", b.WhatIs)

	if err := b.Run(); err != nil {
		b.Logger.Fatal(err.Error())
	}
}

func (b *ExampleBot) Remember(msg joe.Message) error {
	key, value := msg.Matches[0], msg.Matches[1]
	msg.Respond("OK, I'll remember %s is %s", key, value)
	return b.Store.Set(key, value)
}

func (b *ExampleBot) WhatIs(msg joe.Message) error {
	key := msg.Matches[0]
	var value string
	ok, err := b.Store.Get(key, &value)
	if err != nil {
		return errors.Wrapf(err, "failed to retrieve key %q from brain", key)
	}

	if ok {
		msg.Respond("%s is %s", key, value)
	} else {
		msg.Respond("I do not remember %q", key)
	}

	return nil
}
