package service

import (
	"fmt"
	"gems_go_back/pkg/model"
	"gems_go_back/pkg/repository"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"os"
	"strconv"
	"strings"
)

type WithdrawService struct {
	repo repository.Withdraw
}

var errTG error
var bot *tgbotapi.BotAPI

var TOKEN string
var channelID string

func NewWithdrawService(repo repository.Withdraw) *WithdrawService {
	TOKEN = os.Getenv("TELEGRAM_TOKEN")
	channelID = os.Getenv("TELEGRAM_CHANNEL")
	bot, errTG = tgbotapi.NewBotAPI(TOKEN)
	fmt.Println(TOKEN)
	if errTG != nil {
		fmt.Println(errTG)
	}
	return &WithdrawService{repo: repo}
}

func (s *WithdrawService) TelegramBot() {
	go s.HandleUpdatesTelegram(bot)
}

func (s *WithdrawService) CreateWithdraw(currentWithdraw model.Withdraw) error {
	createdWithdraw, err := s.repo.CreateWithdraw(currentWithdraw)
	if err != nil {
		return err
	}

	callbackData := fmt.Sprintf("perform_task_%d", createdWithdraw.ID)

	button := tgbotapi.NewInlineKeyboardButtonData("Выполнить", callbackData)
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(button),
	)

	text := fmt.Sprintf(
		"📋 Новый заказ №%d\n\n"+
			"👤 Пользователь:\n"+
			"├ ID: %s\n"+
			"└ Username: %s\n\n"+
			"🛒 Заказ:\n"+
			"└ Гемы: %d\n",
		createdWithdraw.ID,
		createdWithdraw.UserId,
		createdWithdraw.Username,
		createdWithdraw.Amount)
	msg := tgbotapi.NewMessageToChannel(channelID, text)
	msg.ReplyMarkup = keyboard

	_, err = bot.Send(msg)
	return err
}

func (s *WithdrawService) HandleUpdatesTelegram(bot *tgbotapi.BotAPI) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Panic(err)
	}

	for update := range updates {
		if update.CallbackQuery != nil {
			callback := update.CallbackQuery
			data := callback.Data

			if strings.HasPrefix(data, "perform_task_") {
				orderID := strings.TrimPrefix(data, "perform_task_")
				s.HandlePerformTask(callback, orderID)
			} else if strings.HasPrefix(data, "finish_task_") {
				orderID := strings.TrimPrefix(data, "finish_task_")
				orderIDToInt, _ := strconv.Atoi(orderID)
				currentWithdraw, _ := s.repo.GetWithdraw(orderIDToInt)
				s.HandleFinishTask(callback, currentWithdraw)
			}
		}
	}
}

func (s *WithdrawService) HandlePerformTask(callback *tgbotapi.CallbackQuery, orderID string) {
	user := callback.From

	oldText := callback.Message.Text
	text := "✅ Заказ " + oldText[27:] + fmt.Sprintf("\n\n🤝 @%s взял в работу", user.UserName)

	editMsg := tgbotapi.NewEditMessageText(
		callback.Message.Chat.ID,
		callback.Message.MessageID,
		text,
	)

	_, err := bot.Send(editMsg)
	if err != nil {
		log.Println("Error editing message:", err)
	}

	// Отправляем личное сообщение пользователю с деталями заказа и кнопкой "Завершить"
	callbackDataFinish := fmt.Sprintf("finish_task_%s", orderID)
	finishButton := tgbotapi.NewInlineKeyboardButtonData("Завершить", callbackDataFinish)
	finishKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(finishButton),
	)

	orderIDToInt, err := strconv.Atoi(orderID)
	if err != nil {
		fmt.Println(err)
	}
	currentWithdraw, err := s.repo.GetWithdraw(orderIDToInt)
	taskDetailsToUser := fmt.Sprintf(
		"📋 Новый заказ №%d\n\n"+
			"👤 Пользователь:\n"+
			"├ ID: %s\n"+
			"├ Username: %s\n"+
			"├ Email: %s\n"+
			"└ Code: %d\n\n"+
			"📋 Заказ:\n"+
			"└ Гемы: %d\n",
		currentWithdraw.ID,
		currentWithdraw.UserId,
		currentWithdraw.Username,
		currentWithdraw.AccountEmail,
		currentWithdraw.Code,
		currentWithdraw.Amount)
	privateMsg := tgbotapi.NewMessage(int64(user.ID), taskDetailsToUser)
	privateMsg.ReplyMarkup = finishKeyboard
	_, err = bot.Send(privateMsg)
	if err != nil {
		log.Println("Error sending private message:", err)
	}
	response := tgbotapi.NewCallback(callback.ID, "Вы приняли задачу!")
	bot.AnswerCallbackQuery(response)
}

func (s *WithdrawService) HandleFinishTask(callback *tgbotapi.CallbackQuery, currentWithdraw model.Withdraw) {
	s.repo.CompleteWithdraw(currentWithdraw.ID)

	responseText := fmt.Sprintf(
		"📋 Заказ №%d\n\n"+
			"👤 Пользователь:\n"+
			"├ ID: %s\n"+
			"├ Username: %s\n"+
			"├ Email: %s\n"+
			"└ Code: %d\n\n"+
			"📋 Заказ:\n"+
			"└ Гемы: %d\n\n"+
			"Выполнен ✅✅✅",
		currentWithdraw.ID,
		currentWithdraw.UserId,
		currentWithdraw.Username,
		currentWithdraw.AccountEmail,
		currentWithdraw.Code,
		currentWithdraw.Amount)

	editMsg := tgbotapi.NewEditMessageText(
		callback.Message.Chat.ID,
		callback.Message.MessageID,
		responseText,
	)

	_, err := bot.Send(editMsg)
	if err != nil {
		log.Println("Error editing message:", err)
	}

	response := tgbotapi.NewCallback(callback.ID, "Задание завершено!")
	bot.AnswerCallbackQuery(response)
}
