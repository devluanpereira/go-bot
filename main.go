package main

import (
	"fmt"
	"os"

	"go-bot/handlers"

	"github.com/Valgard/godotenv"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var userLanguages = make(map[int64]string)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		panic(err)
	}

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_API_TOKEN"))
	if err != nil {
		panic(err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {

			switch update.Message.Text {
			case "/start":
				sendWelcomeMessage(bot, update.Message.Chat.ID)
			case "/lock":
				sendPixButton(bot, update.Message.Chat.ID)
			}
		} else if update.CallbackQuery != nil {

			switch update.CallbackQuery.Data {
			case "button1":
				userLanguages[update.CallbackQuery.Message.Chat.ID] = "portuguese"
				clearInlineKeyboard(bot, update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID)
				sendPlanOptions(bot, update.CallbackQuery.Message.Chat.ID, "Bem-vindo! Escolha um plano abaixo:")
			case "button2":
				userLanguages[update.CallbackQuery.Message.Chat.ID] = "english"
				clearInlineKeyboard(bot, update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID)
				sendPlanOptions(bot, update.CallbackQuery.Message.Chat.ID, "Welcome! Choose a plan below:")
			case "button3":
				userLanguages[update.CallbackQuery.Message.Chat.ID] = "spanish"
				clearInlineKeyboard(bot, update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID)
				sendPlanOptions(bot, update.CallbackQuery.Message.Chat.ID, "¡Bienvenido! Elija un plan a continuación:")
			case "plan1":
				sendPixPayment(bot, update.CallbackQuery.Message.Chat.ID)
			case "plan2":
				sendPixPayment(bot, update.CallbackQuery.Message.Chat.ID)
			case "pay":
				qrCodeImage, _, err := handlers.CreatePixPayment()
				if err != nil {
					fmt.Println(err.Error())
					bot.Send(tgbotapi.NewCallback(update.CallbackQuery.ID, "Erro ao gerar QR code de pagamento."))
					continue
				}

				photoFileBytes := tgbotapi.FileBytes{
					Name:  "qrcode.png",
					Bytes: qrCodeImage,
				}
				photoMsg := tgbotapi.NewPhoto(update.CallbackQuery.Message.Chat.ID, photoFileBytes)
				photoMsg.Caption = "Use o QR code a seguir para realizar o pagamento via PIX."
				bot.Send(photoMsg)
			}
		}
	}
}

func sendWelcomeMessage(bot *tgbotapi.BotAPI, chatID int64) {
	imageFileBytes := tgbotapi.FilePath("img/i610.jpg") // Substitua pelo caminho da sua imagem

	btn1 := tgbotapi.NewInlineKeyboardButtonData("Portugues 🇧🇷", "button1")
	btn2 := tgbotapi.NewInlineKeyboardButtonData("English 🇬🇧", "button2")
	btn3 := tgbotapi.NewInlineKeyboardButtonData("Espanhol 🇪🇸", "button3")

	inlineKeys := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(btn1),
		tgbotapi.NewInlineKeyboardRow(btn2),
		tgbotapi.NewInlineKeyboardRow(btn3),
	)

	photoMsg := tgbotapi.NewPhoto(chatID, imageFileBytes)
	photoMsg.Caption = "Bem-vindo! Escolha uma opção abaixo."
	photoMsg.ReplyMarkup = inlineKeys

	bot.Send(photoMsg)
}

func sendPixButton(bot *tgbotapi.BotAPI, chatID int64) {
	inlineBtn := tgbotapi.NewInlineKeyboardButtonData("Pagar com PIX", "pay")
	inlineKeys := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(inlineBtn),
	)

	msg := tgbotapi.NewMessage(chatID, "Bem-vindo! Clique no botão para pagar com PIX.")
	msg.ReplyMarkup = inlineKeys

	bot.Send(msg)
}

func sendPlanOptions(bot *tgbotapi.BotAPI, chatID int64, message string) {
	btn1 := tgbotapi.NewInlineKeyboardButtonData("Plano 1", "plan1")
	btn2 := tgbotapi.NewInlineKeyboardButtonData("Plano 2", "plan2")

	inlineKeys := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(btn1, btn2),
	)

	msg := tgbotapi.NewMessage(chatID, message)
	msg.ReplyMarkup = inlineKeys

	bot.Send(msg)
}

func sendPixPayment(bot *tgbotapi.BotAPI, chatID int64) {
	inlineBtn := tgbotapi.NewInlineKeyboardButtonData("Pagar com PIX", "pay")
	inlineKeys := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(inlineBtn),
	)

	msg := tgbotapi.NewMessage(chatID, "Clique no botão para pagar com PIX.")
	msg.ReplyMarkup = inlineKeys

	bot.Send(msg)
}

func clearInlineKeyboard(bot *tgbotapi.BotAPI, chatID int64, messageID int) {
	edit := tgbotapi.NewEditMessageText(chatID, messageID, "")
	bot.Send(edit)
}
