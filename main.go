package main

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	tb "gopkg.in/tucnak/telebot.v2"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	b, err := tb.NewBot(tb.Settings{
		Token:  os.Getenv("TOKEN"),
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	var commands = make([]tb.Command, 0)

	for name, command := range getCommandList() {
		log.Printf("Command /%s registered", name)
		var f, description = command(b)
		commands = append(commands, tb.Command{Text: name, Description: description})
		b.Handle("/"+name, f)
	}

	err = b.SetCommands(commands)
	if err != nil {
		log.Println("WARNING: Cant update list of commands:")
		log.Println(err)
	}

	b.Start()
}
