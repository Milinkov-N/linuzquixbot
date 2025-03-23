SELECT
    qz.id quiz_id, qz.name quiz_name,
    q.id question_id, q.text question_text,
    a.id answer_id, a.text answer_text, a.is_right answer_is_right
FROM quiz qz
INNER JOIN question q ON qz.id = q.quiz_id
INNER JOIN answer a ON q.id = a.question_id
WHERE qz.id = $1;