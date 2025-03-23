package quiz

import (
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type QuizRepo struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) QuizRepo {
	return QuizRepo{db: db}
}

func (r *QuizRepo) GetQuiz(id uint64) (*Quiz, error) {
	sql, _, err := sq.Select("id, \"name\", callback_data").
		From("quiz").
		Where("id = $1").
		ToSql()

	if err != nil {
		return nil, err
	}

	log.Println(sql)

	quizData := QuizRaw{}
	err = r.db.Get(&quizData, sql, 1)
	if err != nil {
		return nil, err
	}

	log.Println(quizData)

	return quizData.ToQuiz(), nil
}

func (r *QuizRepo) FetchQestions(quiz *Quiz) error {
	sql, _, err := sq.Select("id, \"text\"").
		From("question").
		Where("quiz_id = $1").
		ToSql()

	if err != nil {
		return err
	}

	log.Println(sql)

	questions := []QuestionRaw{}
	stmt, err := r.db.Preparex(sql)
	if err != nil {
		return err
	}

	err = stmt.Select(&questions, quiz.Id)
	if err != nil {
		return err
	}

	log.Println(questions)

	for _, q := range questions {
		quiz.Questions = append(quiz.Questions, *q.ToQuestion())
	}

	return nil
}

func (r *QuizRepo) FetchAnswers(question *Question) error {
	sql, _, err := sq.Select("id, \"text\", is_right").
		From("answer").
		Where("question_id = $1").
		ToSql()

	if err != nil {
		return err
	}

	log.Println(sql)

	answers := []Answer{}
	stmt, err := r.db.Preparex(sql)
	if err != nil {
		return err
	}

	err = stmt.Select(&answers, question.Id)
	if err != nil {
		return err
	}

	log.Println(answers)

	question.Answers = answers

	return nil
}
