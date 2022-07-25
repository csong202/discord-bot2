package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"discord_bots/bot2/utils"
	// "../utils"
)

var (
	BotToken string
	config *configStruct
)

type configStruct struct {
	BotToken string `json : "BotToken"`
}

func ReadConfig() error {
	fmt.Println("Reading config file...")
	file, err := ioutil.ReadFile("./config.json")
	utils.HandleError(err)

	fmt.Println(string(file))
	err = json.Unmarshal(file, &config)
	utils.HandleError(err)

	BotToken = config.BotToken
	return nil
}

