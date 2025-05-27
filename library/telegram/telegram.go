package telegram

import (
	"fmt"
	"log"
	"net/http"

	"github.com/fritz-immanuel/erajaya-be-tech-test/configs"
)

func Send(message string) error {
	config, errConfig := configs.GetConfiguration()
	if errConfig != nil {
		log.Fatalln("failed to get configuration: ", errConfig)
	}

	if config.TeleBotToken == "" || config.TeleGroupID == "" {
		return fmt.Errorf("Telegram Bot Token or Group ID is empty")
	}

	BOT_TOKEN := config.TeleBotToken // ex: 1234567890:ABC-DEF1234ghIkl-zyx57W2v1u123ew11
	GROUP_ID := config.TeleGroupID   // ex: -123456789

	// to use 'enter' (line feed) replace '\n' into '%0A'
	message = fmt.Sprintf("%s%%0A%%0A%s", config.ServerName, message)

	url := "https://api.telegram.org/bot" + BOT_TOKEN + "/sendMessage?chat_id=" + GROUP_ID + "&text=" + message

	_, err := http.Get(url)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	// fmt.Println("===>", resp)

	return err
}
