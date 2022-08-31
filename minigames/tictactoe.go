package minigames

import (
	"fmt"
	"math"

	"discord_bots/bot2/utils"

	"github.com/bwmarrin/discordgo"
)

var (
	boards = make(map[string][][]string)
	
	session *discordgo.Session = nil
	channelID string
	GridPlaces = make([]string, 9)
	gridPlacesUnicode = make([]string, 9)
	done = false;
)

const COMPUTER string = "Computer"
const DRAW string = "Draw"

const EMPTY string = "   "
const ONE string = "1️⃣"
const TWO string = "2️⃣"
const THREE string = "3️⃣"
const FOUR string = "4️⃣"
const FIVE string = "5️⃣"
const SIX string = "6️⃣"
const SEVEN string = "7️⃣"
const EIGHT string = "8️⃣"
const NINE string = "9️⃣"

const X string = "❌"
const O string = "⭕"

func PlayTicTacToe(s *discordgo.Session, chID string, user *discordgo.User, reacted chan bool) {
	boards[user.ID] = [][]string{[]string {EMPTY, EMPTY, EMPTY}, []string {EMPTY, EMPTY, EMPTY}, []string {EMPTY, EMPTY, EMPTY}}
	session = s
	channelID = chID

	GridPlaces = []string{ONE, TWO, THREE, FOUR, FIVE, SIX, SEVEN, EIGHT, NINE}

	userName := user.Username

	var winner string
	turnCount := 0
	for turnCount < 9 {
		if (turnCount % 2 == 0) {
			s.ChannelMessageSend(channelID, fmt.Sprintf("It's %s's turn\n" + getBoard(user.ID), userName))
			playerTurn()
			<- reacted
		} else {
			computerTurn()
		}
		
		winner = checkWin(user)
		if (winner == userName || winner == COMPUTER || winner == DRAW) {
			break
		}
		turnCount++
	}
	
	if (winner == userName) {
		session.ChannelMessageSend(channelID, fmt.Sprintf("Congratulations %s! You beat the computer 🥳", userName))
	} else if (winner == COMPUTER) {
		session.ChannelMessageSend(channelID, "You lost 😢")
	} else if (winner == DRAW) {
		session.ChannelMessageSend(channelID, "No winner 😐")
	}

	session.ChannelMessageSend(channelID, getBoard(user.ID))
}

func playerTurn() {
	msg, _ := session.ChannelMessageSend(channelID, "Choose a spot:\n" + 
		"1 | 2 | 3\n---------\n4 | 5 | 6\n---------\n7 | 8 | 9")

	for _, emoji := range GridPlaces {
		session.MessageReactionAdd(channelID, msg.Reference().MessageID, emoji)
	}
}

func computerTurn() {
	//https://www.freecodecamp.org/news/how-to-make-your-tic-tac-toe-game-unbeatable-by-using-the-minimax-algorithm-9d690bad4b37/
	fmt.Println("computer turn")
}

func HandlePlayerTurn(emoji *discordgo.Emoji, user *discordgo.User) {
	fmt.Printf("%s chose %s\n", user.Username, emoji.Name)

	idx := utils.IndexStr(GridPlaces, emoji.MessageFormat())
	rowIdx := int(math.Floor(float64(idx) / 3))
	boards[user.ID][rowIdx][idx - rowIdx * 3] = "X"

	session.ChannelMessageSend(channelID, getBoard(user.ID))
}

func getBoard(userID string) (output string) {
	output = ""
	board := boards[userID]
	for i, row := range board {
		for j, mark := range row {
			output += mark
			if (j < len(row) - 1) {
				output += " | "
			}
		}
		if (i < len(board) - 1) { output += "\n---------\n" }
	}
	return
}

func checkWin(user *discordgo.User) (winner string) {
	userName := user.Username
	userID := user.ID
	board := boards[userID]

	rowSame := false
	colSame := false
	diagPosSame := false
	diagNegSame := false
	for i, row := range board {
		for j := range row {
			if (j < len(row) - 1) {
				rowSame = row[j] != EMPTY && row[j] == row[j + 1]
			}
			if (i < len(row) - 1) {
				colSame = board[i][j] != EMPTY && board[i][j] == board[i + 1][j]
				if (colSame && i == len(row) - 2) {
					if (board[i][j] == "X") {
						winner = userName
						for k := 0; k < len(board); k ++ {
							board[k][j] = X
						}
						break
					} else if (board[i][j] == "O") {
						winner = COMPUTER
						for k := 0; k < len(board); k++ {
							board[k][j] = O
						}
						break
					}
				}
			}
			if (i == j && i < len(board) - 1) {
				diagPosSame = board[i][j] != EMPTY && board[i][j] == board[i + 1][j + 1]
			}
			if (i == len(board) - j - 1 && j > 0) {
				diagPosSame = board[i][j] != EMPTY && board[i][j] == board[i + 1][j - 1]
			}
			// fmt.Printf("i = %d, j = %d, rowSame: %t, colSame: %t, diagPosSame: %t, diagNegSame: %t\n", 
			// 	i, j, rowSame, colSame, diagPosSame, diagNegSame)
		}
		if (rowSame) {
			if (board[i][0] == "X") {
				winner = userName
				for j := 0; j < len(row); j++ {
					board[i][j] = X
				}
				break
			} else if (board[i][0] == "O") {
				winner = COMPUTER
				for j := 0; j < len(row); j++ {
					board[i][j] = O
				}
				break
			}
		}
	}

	if (diagPosSame) {
		if (board[0][0] == "X") {
			winner = userName
			for i := 0; i < len(board); i ++ {
				board[i][i] = X
			}
		} else if (board[0][0] == "O") {
			winner = COMPUTER
			for i := 0; i < len(board); i ++ {
				board[i][i] = O
			}
		}
	}
	if (diagNegSame) {
		if (board[0][2] == "X") {
			winner = userName
			for i := 0; i < len(board); i++ {
				board[i][len(board) - i - 1] = X
			}
		} else if (board[0][2] == "O") {
			winner = COMPUTER
			for i := 0; i < len(board); i++ {
				board[i][len(board) - i - 1] = O
			}
		}
	}

	boardFull := true
	for i := 0; i < len(board); i++ {
		for j := 0; j < len(board[i]); j++ {
			boardFull = boardFull && board[i][j] != EMPTY
		}
	}
	if (boardFull) { winner = DRAW }
	
	done = true
	return
}