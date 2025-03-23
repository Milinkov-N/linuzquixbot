package quiz

import (
	"log"
	"os"

	"github.com/jmoiron/sqlx"
)

type QuizRepo struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) QuizRepo {
	return QuizRepo{db: db}
}

func (r *QuizRepo) GetQuiz(id uint64) (*Quiz, error) {
	sql, err := os.ReadFile("sql/queries/00-select_all_quiz_data_by_id.sql")
	if err != nil {
		// bot cannot function properly if it can't query the db, hence the panic()
		panic(err.Error())
	}

	stmt, err := r.db.Preparex(string(sql))
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	quizData := []QuizRawData{}
	err = stmt.Select(&quizData, 1)
	if err != nil {
		return nil, err
	}

	for _, v := range quizData {
		log.Println(v)
	}

	return &Quiz{
		Id:   quizData[0].QuizId,
		Name: quizData[0].QuizName,
		Questions: []Question{
			{
				Id:   quizData[0].QuestionId,
				Text: quizData[0].QuestionText,
				Answers: []Answer{
					{
						Id:      quizData[0].AnswerId,
						Text:    quizData[0].AnswerText,
						IsRight: quizData[0].AnswerIsRight,
					},
				},
			},
		},
	}, nil
}
