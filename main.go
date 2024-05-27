package main

import (
	"fmt"
	"os"

	"go-bot/handlers"

	"github.com/Valgard/godotenv"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

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
				// Configura a imagem
				imageFileBytes := tgbotapi.FilePath("img/i610.jpg") // Substitua pelo caminho da sua imagem

				// Configura os botões
				btn1 := tgbotapi.NewInlineKeyboardButtonData("Portugues 🇧🇷", "button1")
				btn2 := tgbotapi.NewInlineKeyboardButtonData("English 🇬🇧", "button2")
				btn3 := tgbotapi.NewInlineKeyboardButtonData("Espanhol 🇪🇸", "button3")

				inlineKeys := tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(btn1),
					tgbotapi.NewInlineKeyboardRow(btn2),
					tgbotapi.NewInlineKeyboardRow(btn3),
				)

				// Envia a imagem com os botões
				photoMsg := tgbotapi.NewPhoto(update.Message.Chat.ID, imageFileBytes)
				photoMsg.Caption = "Bem-vindo! Escolha uma opção abaixo."
				photoMsg.ReplyMarkup = inlineKeys

				bot.Send(photoMsg)
			case "/lock":
				inlineBtn := tgbotapi.NewInlineKeyboardButtonData("Pagar com PIX", "pay")
				inlineKeys := tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(inlineBtn),
				)

				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Bem-vindo! Clique no botão para pagar com PIX.")
				msg.ReplyMarkup = inlineKeys

				bot.Send(msg)
			}
		} else if update.CallbackQuery != nil {
			switch update.CallbackQuery.Data {
			case "pay":
				qrCodeImage, _, err := handlers.CreatePixPayment()
				if err != nil {
					fmt.Println(err.Error())
					bot.Send(tgbotapi.NewCallback(update.CallbackQuery.ID, "Erro ao gerar QR code de pagamento."))
					continue
				}

				// Enviar a imagem do QR code
				photoFileBytes := tgbotapi.FileBytes{
					Name:  "qrcode.png",
					Bytes: qrCodeImage,
				}
				photoMsg := tgbotapi.NewPhoto(update.CallbackQuery.Message.Chat.ID, photoFileBytes)
				photoMsg.Caption = "Use o QR code a seguir para realizar o pagamento via PIX."

				bot.Send(photoMsg)
			case "button1":
				bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Você clicou no botão 1"))
			case "button2":
				bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Você clicou no botão 2"))
			case "button3":
				bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Você clicou no botão 3"))
			}
		}
	}
}
