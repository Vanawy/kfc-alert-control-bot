package main

import (
	"fmt"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

func getCommandList() map[string]func(b *tb.Bot) (func(m *tb.Message), string) {
	var commands = make(map[string]func(b *tb.Bot) (func(m *tb.Message), string))
	commands["ping"] = ping
	commands["coupons_date"] = couponsDate
	return commands
}

// func buildCommand(func)

func ping(b *tb.Bot) (func(m *tb.Message), string) {
	return func(m *tb.Message) {
		b.Send(m.Sender, "pong")
	}, "Sends 'pong' back"
}

func couponsDate(b *tb.Bot) (func(m *tb.Message), string) {
	return func(m *tb.Message) {
		b.Send(m.Sender, fmt.Sprintf(
			"%s - %s",
			time.Now().Format("02 Jan 06"),
			time.Now().Add(time.Hour*24*10).Format("02 Jan 06")))
	}, "Get date of coupons on kfc.by"
}
