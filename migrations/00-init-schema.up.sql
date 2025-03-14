CREATE TABLE quiz (
    id BIGSERIAL PRIMARY KEY,
    "name" VARCHAR(64) 
);

CREATE TABLE question (
    id BIGSERIAL PRIMARY KEY,
    "text" VARCHAR(255),
    quiz_id BIGSERIAL REFERENCES quiz(id)
);

CREATE TABLE answer (
    id BIGSERIAL PRIMARY KEY,
    "text" VARCHAR(32),
    is_right BOOLEAN,
    question_id BIGSERIAL REFERENCES question(id)
);

INSERT INTO quiz ("name") VALUES ('Базовый тест Linux');
INSERT INTO question ("text", quiz_id) VALUES ('Какая команда используется для вывода содержимого папки на экран', 1);
INSERT INTO answer ("text", question_id, is_right) VALUES ('ls', 1, false);
INSERT INTO answer ("text", question_id, is_right) VALUES ('cat', 1, true);
INSERT INTO answer ("text", question_id, is_right) VALUES ('awk', 1, false);
INSERT INTO answer ("text", question_id, is_right) VALUES ('pwd', 1, false);
