package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/joho/godotenv"
)

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
		bot.WithCallbackQueryDataHandler("test", bot.MatchTypePrefix, cbHandler),
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

func cbHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.CallbackQuery.Message.Message.Chat.ID,
		Text:   "You selected the button: " + update.CallbackQuery.Data,
	})
}
