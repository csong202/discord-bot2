package utils

import "fmt"

var ValidArchiveDurations []int = []int{1440, 60, 4320, 10080}

func HandleError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
}

func Contains(slice []int, value int) bool {
	for _, elem := range slice {
		if (elem == value) { return true}
	}
	return false
}