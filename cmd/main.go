package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/joho/godotenv"
)

var selected_quiz = "__none__"

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	err := godotenv.Load()
	if err != nil {
		println("failed to load .env file.", err.Error())
		return
	}

	TG_API_TOKEN := os.Getenv("TG_API_TOKEN")

	opts := []bot.Option{
		bot.WithDefaultHandler(defHandler),
		bot.WithCallbackQueryDataHandler("test", bot.MatchTypePrefix, quizHandler),
		bot.WithCallbackQueryDataHandler("answer", bot.MatchTypePrefix, answerHandler),
	}

	b, err := bot.New(TG_API_TOKEN, opts...)
	if err != nil {
		panic(err)
	}

	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, startHandler)

	b.Start(ctx)
}

func startHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	kb := &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "Общий тест по Linux", CallbackData: "test_main"},
				{Text: "Текстовые редакторы Linux", CallbackData: "test_txt_editors"},
			},
		},
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        "Привет! Выбери один из тестов, которые ты желаешь пройти:",
		ReplyMarkup: kb,
	})
}

func defHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Напиши /start чтобы открыть меню доступных тестов",
	})
}

func quizHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})

	switch update.CallbackQuery.Data {
	case "test_main":
		selected_quiz = "test_main"

		kb := &models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{
				{
					{Text: "Ответ 1", CallbackData: "answer_1"},
					{Text: "Ответ 2", CallbackData: "answer_2"},
					{Text: "Ответ 3", CallbackData: "answer_3"},
					{Text: "Ответ 4", CallbackData: "answer_4"},
				},
			},
		}

		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
			Text:        "You selected the button: " + update.CallbackQuery.Data,
			ReplyMarkup: kb,
		})

	case "test_txt_editors":
		selected_quiz = "test_txt_editors"

		kb := &models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{
				{
					{Text: "Ответ 1", CallbackData: "answer_1"},
					{Text: "Ответ 2", CallbackData: "answer_2"},
					{Text: "Ответ 3", CallbackData: "answer_3"},
				},
			},
		}

		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
			Text:        "You selected the `button`: " + update.CallbackQuery.Data,
			ReplyMarkup: kb,
		})
	}

}

func answerHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.CallbackQuery.Message.Message.Chat.ID,
		Text:   fmt.Sprintf("Your answer for quiz %s was %s", selected_quiz, update.CallbackQuery.Data),
	})
}
