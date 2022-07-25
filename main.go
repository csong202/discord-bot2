package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"

	"discord_bots/bot2/config"
	"discord_bots/bot2/utils"

	// "./config"
	// "./utils"

	"github.com/bwmarrin/discordgo"
)

var (
	BotToken string
	CMD_SYMB = "~"
)

func main() {
	err := config.ReadConfig()
	utils.HandleError(err)

	session, err := discordgo.New("Bot " + config.BotToken)
	utils.HandleError(err)
	fmt.Println("new session")

	session.AddHandler(func (s *discordgo.Session, r *discordgo.Ready) {
		fmt.Println("Bot is ready")
	})

	session.AddHandler(func (s *discordgo.Session, msg *discordgo.MessageCreate) {
		if strings.Contains(msg.Content, CMD_SYMB) {
			if (strings.Contains(msg.Content, CMD_SYMB + "thread")) {
				threadArgs := strings.Split(msg.Content, " ")
				var archiveDuration int
				if (len(threadArgs) < 2) {
					s.ChannelMessageSendReply(msg.ChannelID, "You must include the thread name!", msg.Reference())
					return
				}
				if (len(threadArgs) > 2) {
					var err error
					archiveDuration, err = strconv.Atoi(threadArgs[2])
					if (err != nil) {
						s.ChannelMessageSendReply(msg.ChannelID, "Archive duration must be an integer!", msg.Reference())
						return
					}
					if (!utils.Contains(utils.ValidArchiveDurations, archiveDuration)) {
						s.ChannelMessageSendReply(msg.ChannelID, "Archive duration must be one of (1440, 60, 4320, 10080).", msg.Reference())
						return
					}
				} else {
					archiveDuration = 60 * 24
				}
				_, err := s.MessageThreadStart(msg.ChannelID, msg.ID, threadArgs[1], archiveDuration)
				if err != nil {
					return
				}
			} else {
				s.ChannelMessageSendReply(msg.ChannelID, "not a valid command", msg.Reference())
				// TODO: show help page
			}
		}
	})

	session.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAllWithoutPrivileged)

	err = session.Open()
	utils.HandleError(err)
	defer session.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<- stop
	fmt.Println("Bot2 says bye-bye")
}