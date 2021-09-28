package main

import (
	"log"
	"os"
	"strconv"

	tb "gopkg.in/tucnak/telebot.v2"
)

func registerCommandHandlers(b *tb.Bot) {
	var adminID, _ = strconv.Atoi(os.Getenv("ADMIN_USER_ID"))
	var commands = make([]tb.Command, 0)

	for name, command := range getCommandList() {
		commands = append(commands, tb.Command{
			Text:        name,
			Description: command.Desc,
		})
		b.Handle("/"+name, buildHandler(b, command, adminID))
		log.Printf("Command /%s registered", name)
	}

	var err = b.SetCommands(commands)
	if err != nil {
		log.Println("WARNING: Cant update list of commands:")
		log.Println(err)
	}
}

func buildHandler(b *tb.Bot, command Command, adminID int) func(*tb.Message) {
	log.Println(command)
	return func(m *tb.Message) {
		log.Println(command, m.Sender, adminID)
		if command.AdminOnly && m.Sender.ID != adminID {
			b.Send(m.Sender, "Access denied")
			return
		}
		command.Func(b, m)
	}
}
