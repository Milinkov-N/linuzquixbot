package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/joho/godotenv"

	"github.com/Milinkov-N/linuzquixbot/pkg/cache"
	"github.com/Milinkov-N/linuzquixbot/pkg/postgres"
	"github.com/Milinkov-N/linuzquixbot/pkg/quiz"
)

type UserData struct {
	selected_quiz string
}

var QUIZZES []*quiz.Quiz
var USERS = cache.NewCacheWithAutoCleanup[int64](1*time.Minute, 5*time.Minute)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	err := godotenv.Load()
	if err != nil {
		log.Println("failed to load .env file.", err.Error())
		return
	}

	TG_API_TOKEN := os.Getenv("TG_API_TOKEN")

	db, err := postgres.New("postgres", "postgres")
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	repo := quiz.NewRepo(db)
	quiz, err := repo.GetQuiz(1)
	if err != nil {
		log.Panic(err.Error())
	}

	err = repo.FetchQestions(quiz)
	if err != nil {
		log.Panic(err.Error())
	}

	for i := range quiz.Questions {
		err := repo.FetchAnswers(&quiz.Questions[i])
		if err != nil {
			log.Panic(err.Error())
		}
	}

	fmt.Println(quiz)

	QUIZZES = append(QUIZZES, quiz)

	opts := []bot.Option{
		bot.WithDefaultHandler(defHandler),
		bot.WithCallbackQueryDataHandler("test:", bot.MatchTypePrefix, quizHandler),
		bot.WithCallbackQueryDataHandler("answer:", bot.MatchTypePrefix, answerHandler),
	}

	b, err := bot.New(TG_API_TOKEN, opts...)
	if err != nil {
		panic(err)
	}

	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, startHandler)

	b.Start(ctx)
}

func startHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	kb := &models.InlineKeyboardMarkup{}

	for _, quiz := range QUIZZES {
		kb.InlineKeyboard = append(kb.InlineKeyboard, []models.InlineKeyboardButton{
			{Text: quiz.Name, CallbackData: quiz.CallbackData},
		})
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

	// FIX: QUIZZES should be a map to callback_data
	// so  this hardcode can be eliminated
	switch update.CallbackQuery.Data {
	case "test:linux_main":
		USERS.Set(update.CallbackQuery.Message.Message.Chat.ID, UserData{
			selected_quiz: "test:linux_main",
		})

		kb := &models.InlineKeyboardMarkup{}

		for _, answer := range QUIZZES[0].Questions[0].Answers {
			kb.InlineKeyboard = append(kb.InlineKeyboard, []models.InlineKeyboardButton{
				{Text: answer.Text, CallbackData: fmt.Sprintf("answer:%d", answer.Id)},
			})
		}

		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
			Text:        "You selected the button: " + update.CallbackQuery.Data,
			ReplyMarkup: kb,
		})

	case "test:txt_editors":
		USERS.Set(update.CallbackQuery.Message.Message.Chat.ID, UserData{
			selected_quiz: "test:txt_editors",
		})

		kb := &models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{
				{
					{Text: "Ответ 1", CallbackData: "answer:1"},
					{Text: "Ответ 2", CallbackData: "answer:2"},
					{Text: "Ответ 3", CallbackData: "answer:3"},
				},
			},
		}

		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
			Text:        "You selected the button: " + update.CallbackQuery.Data,
			ReplyMarkup: kb,
		})
	}

}

func answerHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})

	data, exists := USERS.Get(update.CallbackQuery.Message.Message.Chat.ID)

	if !exists {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.CallbackQuery.Message.Message.Chat.ID,
			Text:   "Похоже время на ответ истекло. Введи команду /start и попробуй пройти тест заново",
		})

		return
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.CallbackQuery.Message.Message.Chat.ID,
		Text: fmt.Sprintf("Your answer for quiz %s was %s",
			data.(UserData).selected_quiz,
			update.CallbackQuery.Data),
	})
}
