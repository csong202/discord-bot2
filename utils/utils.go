package utils

import (
	"bot2/minigamestypes"
	"errors"
	"fmt"
)

var ValidArchiveDurations []int = []int{1440, 60, 4320, 10080}

func HandleError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
}

func ContainsInt(slice []int, value int) bool {
	for _, elem := range slice {
		if elem == value {
			return true
		}
	}
	return false
}

func ContainsStr(slice []string, value string) bool {
	for _, elem := range slice {
		if elem == value {
			return true
		}
	}
	return false
}

func RemoveStr(slice *[]string, value string) error {
	if !ContainsStr(*slice, value) {
		return errors.New("slice does not contain the value")
	}

	n := len(*slice) - 1
	newSlice := make([]string, n)

	idx := 0
	if n > 0 {
		for i, elem := range *slice {
			if elem == value {
				newSlice[i] = ""
				idx = i
				break
			}
		}
	}

	for i := idx; i < n; i++ {
		newSlice[i] = (*slice)[i+1]
	}

	*slice = make([]string, n)
	for i := 0; i < n; i++ {
		(*slice)[i] = newSlice[i]
	}

	return nil
}

func IndexStr(slice []string, value string) int {
	for i, elem := range slice {
		if elem == value {
			return i
		}
	}
	return -1
}

func GetMapKeys(myMap map[string]*minigamestypes.TicTacToeGameMeta) (keys []string) {
	keys = make([]string, len(myMap))
	i := 0
	for key := range myMap {
		keys[i] = key
		i++
	}
	return
}

func Copy2DSliceStr(matrix [][]string) [][]string {
	matrixCopy := make([][]string, len(matrix))
	for i, row := range matrix {
		matrixCopy[i] = make([]string, len(matrix[i]))
		copy(matrixCopy[i], row)
	}
	return matrixCopy
}
