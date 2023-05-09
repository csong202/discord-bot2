package minigamestypes

type TicTacToeGameMeta struct {
	ReactChannel chan bool
	LastMsgID    string
	Board        [][]string
	ChID         string
}
