package quiz

type QuizRawData struct {
	QuizId        uint64 `db:"quiz_id"`
	QuizName      string `db:"quiz_name"`
	QuestionId    uint64 `db:"question_id"`
	QuestionText  string `db:"question_text"`
	AnswerId      uint64 `db:"answer_id"`
	AnswerText    string `db:"answer_text"`
	AnswerIsRight bool   `db:"answer_is_right"`
}

type Quiz struct {
	Id        uint64
	Name      string
	Questions []Question
}

type Question struct {
	Id      uint64
	Text    string
	Answers []Answer
}

type Answer struct {
	Id      uint64
	Text    string
	IsRight bool
}
