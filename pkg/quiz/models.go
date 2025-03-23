package quiz

type Quiz struct {
	Id           uint64
	Name         string
	CallbackData string
	Questions    []Question
}

type Question struct {
	Id      uint64
	Text    string
	Answers []Answer
}

type Answer struct {
	Id      uint64 `db:"id"`
	Text    string `db:"text"`
	IsRight bool   `db:"is_right"`
}

type QuizRaw struct {
	Id           uint64 `db:"id"`
	Name         string `db:"name"`
	CallbackData string `db:"callback_data"`
}

func (qr *QuizRaw) ToQuiz() *Quiz {
	return &Quiz{
		Id:           qr.Id,
		Name:         qr.Name,
		CallbackData: qr.CallbackData,
		Questions:    []Question{},
	}
}

type QuestionRaw struct {
	Id   uint64 `db:"id"`
	Text string `db:"text"`
}

func (qr *QuestionRaw) ToQuestion() *Question {
	return &Question{
		Id:      qr.Id,
		Text:    qr.Text,
		Answers: []Answer{},
	}
}
