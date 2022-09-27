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
const DRAW string = "GAME_DRAW"
const NO_WIN string = "NO_WIN"

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

const X string = "X"
const O string = "O"
const X_BOLD string = "**X**"
const O_BOLD string = "**O**"

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
			computerTurn(user)
		}
		
		winner = checkWin(boards[user.ID])
		if (winner != NO_WIN) {
			break
		}
		turnCount++
	}
	
	if (winner == X) {
		session.ChannelMessageSend(channelID, fmt.Sprintf("Congratulations %s! You beat the computer 🥳", userName))
	} else if (winner == O) {
		session.ChannelMessageSend(channelID, "You lost 😢")
	} else if (winner == DRAW) {
		session.ChannelMessageSend(channelID, "No winner 😐")
	}

	session.ChannelMessageSend(channelID, getBoard(user.ID))
}

func playerTurn() {
	msg, _ := session.ChannelMessageSend(channelID, "Choose a spot:\n" + 
		"1 | 2 | 3\n---------\n4 | 5 | 6\n---------\n7 | 8 | 9\n")

	for _, emoji := range GridPlaces {
		session.MessageReactionAdd(channelID, msg.Reference().MessageID, emoji)
	}
}

func HandlePlayerTurn(emoji *discordgo.Emoji, user *discordgo.User) {
	fmt.Printf("%s chose %s\n", user.Username, emoji.Name)

	idx := utils.IndexStr(GridPlaces, emoji.MessageFormat())
	rowIdx, colIdx := getBoardPosFromIdx(idx)
	boards[user.ID][rowIdx][colIdx] = "X"

	session.ChannelMessageSend(channelID, getBoard(user.ID))
}

func computerTurn(user *discordgo.User) {
	fmt.Println("computer turn")
	availRows, _ := findAvailSpaces(boards[user.ID])
	scoreMultiplier := len(availRows)
	_, rowIdx, colIdx := miniMaxRecursive(boards[user.ID], O, scoreMultiplier)
	if (rowIdx > -1 && colIdx > -1) {
		boards[user.ID][rowIdx][colIdx] = "O"
	}
}

func miniMaxRecursive(board [][]string, playerSymb string, scoreMultiplier int) (bestScore int, bestScoreRow int, bestScoreCol int) {
	winner := checkWin(board)
	var terminalScore int
	if (winner != NO_WIN) {
		if (winner == X) {
			terminalScore = -1
		} else if (winner == O) {
			terminalScore = 1
		} else if (winner == DRAW) {
			terminalScore = 0
		}
		return terminalScore * scoreMultiplier, -1, -1
	}

	bestScoreRow = -1
	bestScoreCol = -1
	if (playerSymb == X) {
		bestScore = math.MaxInt
	} else if (playerSymb == O) {
		bestScore = math.MinInt
	}
	
	availRows, availCols := findAvailSpaces(board)
	for i := range availRows {
		newBoard := utils.Copy2DSliceStr(board)
		newBoard[availRows[i]][availCols[i]] = playerSymb
		var nextPlayerSymb string
		if (playerSymb == X) {
			nextPlayerSymb = O
		} else if (playerSymb == O) {
			nextPlayerSymb = X
		}
		currScore, _, _ := miniMaxRecursive(newBoard, nextPlayerSymb, (scoreMultiplier - 1))
		if (playerSymb == X && currScore < bestScore) {
			bestScore = currScore
			bestScoreRow = availRows[i]
			bestScoreCol = availCols[i]
		} else if (playerSymb == O && currScore > bestScore) {
			bestScore = currScore
			bestScoreRow = availRows[i]
			bestScoreCol = availCols[i]
		}
	}
	return 
}

func findAvailSpaces(board [][]string) (rows []int, cols []int){
	for i, row := range board {
		for j, item := range row {
			if (item == EMPTY) {
				rows = append(rows, i)
				cols = append(cols, j)
			}
		}
	}
	return
}

func getBoardPosFromIdx(idx int) (row int, col int) {
	row = int(math.Floor(float64(idx) / 3))
	col = idx - row * 3
	return
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

func checkWin(board [][]string) (winner string) {
	rowSame := false
	colSame := false
	diagPosSame := false
	diagNegSame := false
	for i, row := range board {
		for j := range row {
			if (j < len(row) - 1) {
				rowSame = row[j] != EMPTY && row[j] == row[j + 1]
			}
			if (i < len(board) - 1) {
				colSame = board[i][j] != EMPTY && board[i][j] == board[i + 1][j]
				if (colSame && i == len(row) - 2) {
					if (board[i][j] == "X") {
						winner = X
						for k := 0; k < len(board); k ++ {
							board[k][j] = X_BOLD
						}
						break
					} else if (board[i][j] == "O") {
						winner = O
						for k := 0; k < len(board); k++ {
							board[k][j] = O_BOLD
						}
						break
					}
				}
			}
			if (i == j && i < len(board) - 1) {
				diagPosSame = board[i][j] != EMPTY && board[i][j] == board[i + 1][j + 1]
			}
			if (i == len(board) - j - 1 && j > 0) {
				diagNegSame = board[i][j] != EMPTY && board[i][j] == board[i + 1][j - 1]
			}
		}
		if (rowSame) {
			if (board[i][0] == "X") {
				winner = X
				for j := 0; j < len(row); j++ {
					board[i][j] = X_BOLD
				}
				break
			} else if (board[i][0] == "O") {
				winner = O
				for j := 0; j < len(row); j++ {
					board[i][j] = O_BOLD
				}
				break
			}
		}
	}

	if (diagPosSame) {
		if (board[0][0] == "X") {
			winner = X
			for i := 0; i < len(board); i ++ {
				board[i][i] = X_BOLD
			}
		} else if (board[0][0] == "O") {
			winner = O
			for i := 0; i < len(board); i ++ {
				board[i][i] = O_BOLD
			}
		}
	}
	if (diagNegSame) {
		if (board[0][2] == "X") {
			winner = X
			for i := 0; i < len(board); i++ {
				board[i][len(board) - i - 1] = X_BOLD
			}
		} else if (board[0][2] == "O") {
			winner = O
			for i := 0; i < len(board); i++ {
				board[i][len(board) - i - 1] = O_BOLD
			}
		}
	}

	boardFull := true
	for i := 0; i < len(board); i++ {
		for j := 0; j < len(board[i]); j++ {
			boardFull = boardFull && board[i][j] != EMPTY
		}
	}
	if (boardFull) { 
		winner = DRAW 
	} else {
		winner = NO_WIN
	}
	
	done = true
	return
}