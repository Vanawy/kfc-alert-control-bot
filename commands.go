package main

import (
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strings"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

type Command struct {
	Func      func(*tb.Bot, *tb.Message)
	Desc      string
	AdminOnly bool
}

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
	var cmd = exec.Command("pm2", "status")
	out, err := cmd.Output()
	if err != nil {
		log.Println(err)
		b.Send(m.Sender, fmt.Sprintf("Error: %v", err))
	}

	rows := strings.Split(string(out), "\n")

	var title []string
	var processes [][]string

	for i, row := range rows {
		fields := strings.Split(row, "│")
		if len(fields) <= 1 {
			continue
		}
		var processInfo []string
		addProcessInfo := false
		for _, field := range fields {
			if i == 1 {
				title = append(title, strings.TrimSpace(field))
			} else {
				processInfo = append(processInfo, strings.TrimSpace(field))
				addProcessInfo = true
			}
		}
		if addProcessInfo {
			processes = append(processes, processInfo)
		}
	}

	pattern, _ := regexp.Compile("^(id|name|version|uptime|cpu|mem|status)$")

	var result string

	for i, f := range title {
		match := pattern.MatchString(f)
		if match {
			pInfo := ""
			for _, p := range processes {
				pInfo += p[i] + " "
			}
			result = result + fmt.Sprintf("%s | %s\n", f, pInfo)
		}
	}

	b.Send(m.Sender, result)
}

func logs(b *tb.Bot, m *tb.Message) {
	var cmd = exec.Command("tail", "/home/vanawy/.pm2/logs/test-out.log")
	out, err := cmd.Output()
	if err != nil {
		log.Println(err)
		b.Send(m.Sender, fmt.Sprintf("Error: %v", err))
	}

	b.Send(m.Sender, string(out))
}
