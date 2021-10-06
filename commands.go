package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

type Command struct {
	Func      func(*tb.Bot, *tb.Message)
	Desc      string
	AdminOnly bool
}

const (
	P_NAME   = 1 + 1
	P_UPTIME = 6 + 1
	P_STATUS = 8 + 1
	P_CPU    = 9 + 1
	P_MEM    = 10 + 1
)

func getCommandList() map[string]Command {

	var commands = make(map[string]Command)
	commands["ping"] = Command{
		Func:      ping,
		Desc:      "Sends 'pong' back",
		AdminOnly: false,
	}
	commands["coupons_date"] = Command{
		Func:      couponsDate,
		Desc:      "Get date of coupons on kfc.by",
		AdminOnly: true,
	}
	commands["status"] = Command{
		Func:      pm2status,
		Desc:      "Run pm2 status and get bot info",
		AdminOnly: true,
	}
	commands["logs"] = Command{
		Func:      logs,
		Desc:      "Get logs",
		AdminOnly: true,
	}
	return commands
}

func ping(b *tb.Bot, m *tb.Message) {
	b.Send(m.Sender, "pong 🏓")
}

func couponsDate(b *tb.Bot, m *tb.Message) {
	b.Send(m.Sender, fmt.Sprintf(
		"%s - %s",
		time.Now().Format("02 Jan 06"),
		time.Now().Add(time.Hour*24*10).Format("02 Jan 06")))
}

func pm2status(b *tb.Bot, m *tb.Message) {
	var pm2status = exec.Command("pm2", "status")
	pipe, _ := pm2status.StdoutPipe()

	var grep = exec.Command("grep", os.Getenv("BOT_PROCESS_NAME"))
	grep.Stdin = pipe

	pm2status.Start()
	out, err := grep.Output()
	if err != nil {
		log.Println(err)
		b.Send(m.Sender, fmt.Sprintf("Error: %v", err))
	}

	data := strings.Split(string(out), "│")
	if len(data) < 10 {
		b.Send(m.Sender, "Something went wrong")
	}

	var result string = fmt.Sprintf("Mem: %s\nCPU: %s\nUptime: %s\nStatus: %s",
		data[P_MEM], data[P_CPU], data[P_UPTIME], data[P_STATUS])

	b.Send(m.Sender, result)
}

func logs(b *tb.Bot, m *tb.Message) {
	var cmd = exec.Command("tail", os.Getenv("LOGS_LOCATION"))
	out, err := cmd.Output()
	if err != nil {
		log.Println(err)
		b.Send(m.Sender, fmt.Sprintf("Error: %v", err))
	}

	b.Send(m.Sender, string(out))
}
