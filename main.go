package main

import (
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"

	"discord_bots/bot2/config"
	"discord_bots/bot2/minigames"
	"discord_bots/bot2/utils"

	"github.com/bwmarrin/discordgo"
)

var (
	BotToken string
	BOT_PREFIX = "~"
	searchParams map[string]string
	pinPatterns = make([]string, 0, 20)
	tictactoeGames = make(map[string]chan bool)
)

func main() {
	err := config.ReadConfig()
	utils.HandleError(err)

	session, err := discordgo.New("Bot " + config.BotToken)
	utils.HandleError(err)
	fmt.Println("new session")

	// variable initialization
	searchParams = make(map[string]string)
	searchParams["from"] = "" // TODO autocomplete
	searchParams["has"] = "" // TODO for has and in: popup with options
	searchParams["in"] = ""
	
	// TODO during param. somehow make calendar pop up

	// handlers

	session.AddHandler(func (s *discordgo.Session, r *discordgo.Ready) {
		fmt.Println("Bot is ready")
	})
	session.AddHandler(messageHandler)
	session.AddHandler(reactionHandler)

	session.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAllWithoutPrivileged)

	err = session.Open()
	utils.HandleError(err)
	defer session.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<- stop
	fmt.Println("Bot2 is going to sleep")
}

func messageHandler(s *discordgo.Session, msg *discordgo.MessageCreate) {
	if len(msg.Content) > 0 && string(msg.Content[0]) == BOT_PREFIX {
		if (strings.Contains(msg.Content, BOT_PREFIX + "thread")) {
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
				if (!utils.ContainsInt(utils.ValidArchiveDurations, archiveDuration)) {
					s.ChannelMessageSendReply(msg.ChannelID, "Archive duration must be one of (1440, 60, 4320, 10080).", msg.Reference())
					return
				}
			} else {
				archiveDuration = 60 * 24
			}
			_, err := s.MessageThreadStart(msg.ChannelID, msg.ID, threadArgs[1], archiveDuration)
			utils.HandleError(err)
		} else if (strings.Contains(msg.Content, BOT_PREFIX + "regex_search")) {
			searchArgs := strings.Split(msg.Content, " ")
			var pattern string
			for _, arg := range searchArgs {
				if (arg == BOT_PREFIX + "regexStr_search") {
					continue
				} else if strings.HasPrefix(arg, "from:") {
					searchParams["from"] = arg[5:len(arg)-1]
				} else if strings.HasPrefix(arg, "has:") {
					searchParams["from"] = arg[5:len(arg)-1]
				} else if strings.HasPrefix(arg, "in:") {
					searchParams["from"] = arg[5:len(arg)-1]
				} else {
					pattern = arg
				}
			}

			matched, err := regexp.MatchString(pattern, "...")
			fmt.Println(matched)
			fmt.Println(err)
		} else if (strings.Contains(msg.Content, BOT_PREFIX + "regex_pin")) {
			pinPatArgs := strings.Split(msg.Content, " ")
			if len(pinPatArgs) < 3 {
				s.ChannelMessageSendReply(msg.ChannelID, "Command must be in the form ~regex_pin <add|remove> <pattern>", msg.Reference())
				return
			}
			if pinPatArgs[1] == "add" {
				pinPatterns = append(pinPatterns, pinPatArgs[2])
				fmt.Printf("Just added %s, now have %d pin patterns\n", pinPatArgs[2], len(pinPatterns))
			} else if pinPatArgs[1] == "remove" {
				err := utils.RemoveStr(&pinPatterns, pinPatArgs[2])
				if err != nil {
					s.ChannelMessageSend(msg.ChannelID, "This pattern was never added")
				} else {
					fmt.Printf("Just removed %s, now have %d pin patterns\n", pinPatArgs[2], len(pinPatterns))
				}
			} else {
				s.ChannelMessageSendReply(msg.ChannelID, "invalid argument! must be add or remove", msg.Reference())
				return
			}

		} else if (strings.Contains(msg.Content, BOT_PREFIX + "tictactoe")) {
			s.ChannelMessageSend(msg.ChannelID, "Welcome to TicTacToe! You are X and the computer is O")
			reacted := make(chan bool)
			tictactoeGames[msg.Author.ID] = reacted
			go minigames.PlayTicTacToe(s, msg.ChannelID, msg.Author, reacted)
		} else {
			s.ChannelMessageSendReply(msg.ChannelID, "not a valid command", msg.Reference())
			// TODO: show help page
		} 

	} else if checkMsgMatchPinPatterns(msg.Content) {
		s.ChannelMessagePin(msg.ChannelID, msg.ID)
	}
}

func reactionHandler(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	if (utils.ContainsStr(utils.GetMapKeys(tictactoeGames), r.Member.User.ID) && 
		utils.ContainsStr(minigames.GridPlaces, r.Emoji.MessageFormat())) {
		minigames.HandlePlayerTurn(&r.Emoji, r.Member.User)
		tictactoeGames[r.Member.User.ID] <- true
	}
}

func checkMsgMatchPinPatterns(msg string) bool {
	for _, pattern := range pinPatterns {
		if pattern == "" { continue }
		if matched, _ := regexp.MatchString(pattern, msg); matched{
			fmt.Printf("Message \"%s\" matched %s\n", msg, pattern)
			return true
		}
	}
	return false
}