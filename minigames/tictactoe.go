package minigames

import (
	"fmt"
	"math"
	"math/rand"

	"bot2/utils"

	"github.com/bwmarrin/discordgo"
)

var (
	TicTacToeGames = make(map[string]chan bool)
	boards         = make(map[string][][]string)

	session       *discordgo.Session = nil
	channelID     string
	GridPlaces    = make([]string, 9)
	boldedSymbols = make(map[string]string)
	randomPlays   = false
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
	boards[user.ID] = [][]string{{EMPTY, EMPTY, EMPTY}, {EMPTY, EMPTY, EMPTY}, {EMPTY, EMPTY, EMPTY}}
	session = s
	channelID = chID

	GridPlaces = []string{ONE, TWO, THREE, FOUR, FIVE, SIX, SEVEN, EIGHT, NINE}
	boldedSymbols[X] = X_BOLD
	boldedSymbols[O] = O_BOLD

	userName := user.Username

	var winner string
	turnCount := 0
	for turnCount < 9 {
		if turnCount%2 == 0 {
			s.ChannelMessageSend(channelID, fmt.Sprintf("It's %s's turn\n"+getBoard(user.ID), userName))
			playerTurn()
			<-reacted
		} else {
			computerTurn(user.ID)
		}

		fmt.Println(getBoard(user.ID))
		winner = checkWin(boards[user.ID])
		if winner != NO_WIN {
			fmt.Println("Winner is: " + winner)
			break
		}
		turnCount++
	}

	if winner == X {
		session.ChannelMessageSend(channelID, fmt.Sprintf("Congratulations %s! You beat the computer 🥳", userName))
	} else if winner == O {
		session.ChannelMessageSend(channelID, "You lost 😢")
	} else if winner == DRAW {
		session.ChannelMessageSend(channelID, "No winner 😐")
	}

	session.ChannelMessageSend(channelID, getBoard(user.ID))
}

func playerTurn() {
	msg, _ := session.ChannelMessageSend(channelID, "Choose a spot:\n"+
		"1 | 2 | 3\n---------\n4 | 5 | 6\n---------\n7 | 8 | 9\n")

	for _, emoji := range GridPlaces {
		session.MessageReactionAdd(channelID, msg.Reference().MessageID, emoji)
	}
}

func HandlePlayerTurn(emoji *discordgo.Emoji, user *discordgo.User) {
	fmt.Printf("%s chose %s\n", user.Username, emoji.Name)

	idx := utils.IndexStr(GridPlaces, emoji.MessageFormat())
	rowIdx, colIdx := getBoardPosFromIdx(idx)

	if boards[user.ID][rowIdx][colIdx] == EMPTY {
		boards[user.ID][rowIdx][colIdx] = X
		session.ChannelMessageSend(channelID, getBoard(user.ID))
		TicTacToeGames[user.ID] <- true
	} else {
		session.ChannelMessageSend(channelID, "This spot is already taken. React again to choose another one")
	}

}

func computerTurn(userID string) {
	availRows, availCols := findAvailSpaces(boards[userID])

	var (
		scoreMultiplier int
		rowIdx          int
		colIdx          int
	)

	if randomPlays {
		availIdx := rand.Intn(len(availRows))
		rowIdx = availRows[availIdx]
		colIdx = availCols[availIdx]
	} else {
		scoreMultiplier = len(availRows)
		_, rowIdx, colIdx = miniMaxMove(utils.Copy2DSliceStr(boards[userID]), O, scoreMultiplier)
	}

	if rowIdx > -1 && colIdx > -1 {
		boards[userID][rowIdx][colIdx] = O
		fmt.Println("Bot turn, bot chose " + GridPlaces[getIdxFromBoardPos(rowIdx, colIdx)])
	}
}

func miniMaxMove(board [][]string, playerSymb string, scoreMultiplier int) (bestScore int, bestScoreRow int, bestScoreCol int) {
	winner := checkWin(board)
	var terminalScore int
	if winner != NO_WIN {
		if winner == DRAW {
			terminalScore = 0
		} else if winner == O {
			terminalScore = 1
		} else if winner == X {
			terminalScore = -1
		}
		return terminalScore * scoreMultiplier, -1, -1
	}

	bestScoreRow = -1
	bestScoreCol = -1
	if playerSymb == X {
		bestScore = math.MaxInt
	} else if playerSymb == O {
		bestScore = math.MinInt
	}

	availRows, availCols := findAvailSpaces(board)
	for i := range availRows {
		newBoard := utils.Copy2DSliceStr(board)
		newBoard[availRows[i]][availCols[i]] = playerSymb
		var nextPlayerSymb string
		if playerSymb == X {
			nextPlayerSymb = O
		} else if playerSymb == O {
			nextPlayerSymb = X
		}
		currScore, _, _ := miniMaxMove(newBoard, nextPlayerSymb, len(availRows))
		if playerSymb == X && currScore < bestScore {
			bestScore = currScore
			bestScoreRow = availRows[i]
			bestScoreCol = availCols[i]
		} else if playerSymb == O && currScore > bestScore {
			bestScore = currScore
			bestScoreRow = availRows[i]
			bestScoreCol = availCols[i]
		}
	}
	return
}

func findAvailSpaces(board [][]string) (rows []int, cols []int) {
	for i, row := range board {
		for j, item := range row {
			if item == EMPTY {
				rows = append(rows, i)
				cols = append(cols, j)
			}
		}
	}
	return
}

func getBoardPosFromIdx(idx int) (row int, col int) {
	row = int(math.Floor(float64(idx) / 3))
	col = idx - row*3
	return
}

func getIdxFromBoardPos(row int, col int) int {
	return 3*row + col
}

func getBoard(userID string) (output string) {
	output = ""
	board := boards[userID]
	for i, row := range board {
		for j, mark := range row {
			output += mark
			if j < len(row)-1 {
				output += " | "
			}
		}
		if i < len(board)-1 {
			output += "\n---------\n"
		}
	}
	return
}

func checkWin(board [][]string) (winner string) {
	winner = NO_WIN
	rowSame := true
	colSame := true
	diagPosSame := true
	diagNegSame := true

	for i, row := range board {
		rowSame = true
		for j := range row {
			if j < len(row)-1 {
				rowSame = rowSame && row[j] != EMPTY && row[j] == row[j+1]
			}
		}
		if rowSame {
			winner = board[i][0]
			for j := 0; j < len(row); j++ {
				board[i][j] = boldedSymbols[winner]
			}
			break
		}
	}

	for j := range board[0] {
		colSame = true
		for i := range board {
			if i < len(board)-1 {
				colSame = colSame && board[i][j] != EMPTY && board[i][j] == board[i+1][j]
			}
		}
		if colSame {
			winner = board[0][j]
			for i := 0; i < len(board); i++ {
				board[i][j] = boldedSymbols[winner]
			}
			break
		}
	}

	for i := range board {
		if i < len(board)-1 {
			diagPosSame = diagPosSame && board[i][i] != EMPTY && board[i][i] == board[i+1][i+1]
		}
		if i > 0 {
			diagNegSame = diagNegSame && board[i][len(board)-i-1] != EMPTY && board[i][len(board)-i-1] == board[i-1][len(board)-i]
		}
	}

	if diagPosSame {
		winner = board[0][0]
		for i := 0; i < len(board); i++ {
			board[i][i] = boldedSymbols[winner]
		}
	}
	if diagNegSame {
		winner = board[0][len(board)-1]
		for i := 0; i < len(board); i++ {
			board[i][len(board)-i-1] = boldedSymbols[winner]
		}
	}

	if winner == NO_WIN {
		boardFull := true
	outer:
		for i := 0; i < len(board); i++ {
			for j := 0; j < len(board[i]); j++ {
				boardFull = boardFull && board[i][j] != EMPTY
				if !boardFull {
					break outer
				}
			}
		}
		if boardFull {
			winner = DRAW
		}
	}

	// fmt.Printf("rowSame = %t, colSame = %t, diagPosSame = %t, diagNegSame = %t\n",
	// 	rowSame, colSame, diagPosSame, diagNegSame)

	return
}
