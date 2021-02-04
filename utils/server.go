package utils

import (
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// TelegramServer is a Object for Telegram API Server
type TelegramServer struct {
	BotConnetion        *tgbotapi.BotAPI
	updateConfig        tgbotapi.UpdateConfig
	consumerFunctionMap map[string]func(updates *tgbotapi.Update) error
}

// NewTelegramServer is a constructor for TelegramServer
func NewTelegramServer(botAPIToken string, debug bool, timeout int) (*TelegramServer, error) {

	server := new(TelegramServer)

	bot, err := tgbotapi.NewBotAPI(botAPIToken)

	if err != nil {
		return nil, err
	}

	bot.Debug = true

	log.Printf("[Telegram Server] Authorized on account %s", bot.Self.UserName)

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = timeout

	server.BotConnetion = bot
	server.updateConfig = updateConfig
	server.consumerFunctionMap = make(map[string]func(updates *tgbotapi.Update) error)

	return server, nil
}

// RunServer is a function to run Telegram Server
func (server *TelegramServer) RunServer() {

	updateChan, err := server.BotConnetion.GetUpdatesChan(server.updateConfig)

	if err != nil {
		log.Fatalf("[Telegram Server] Error Running Consumer Server : %s\n", err.Error())
	}

	for update := range updateChan {
		go server.consumeUpdate(&update)
	}

}

// RegisterConsumerController is a function to register a function controller
func (server *TelegramServer) RegisterConsumerController(command string, controllerFunction func(updates *tgbotapi.Update) error) {

	if _, ok := server.consumerFunctionMap[command]; ok {
		log.Fatalf("[Telegram Server][RegisterConsumerController] '%s' command duplicated\n", command)
	}

	server.consumerFunctionMap[command] = controllerFunction

}

func (server *TelegramServer) consumeUpdate(updates *tgbotapi.Update) {

	command := strings.Split(updates.Message.Text, " ")[0]

	function, ok := server.consumerFunctionMap[command]

	if !ok {
		msgText := "Command Not Found, Please Try Again"
		server.sendSimpleMessage(updates.Message.Chat.ID, msgText)
		return
	}

	err := function(updates)

	if err != nil {
		msgText := "Function Error, Please Try Again"
		server.sendSimpleMessage(updates.Message.Chat.ID, msgText)
		return
	}

}

func (server *TelegramServer) sendSimpleMessage(chatID int64, messageText string) error {
	msg := tgbotapi.NewMessage(chatID, messageText)
	_, err := server.BotConnetion.Send(msg)
	return err
}
