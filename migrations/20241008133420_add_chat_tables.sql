-- +goose Up
-- +goose StatementBegin
CREATE TABLE chat
(
    id         SERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE chat_user
(
    id      SERIAL PRIMARY KEY,
    chat_id INT REFERENCES chat (id) ON DELETE CASCADE NOT NULL,
    user_id INT                                        NOT NULL,
    CONSTRAINT uq_chat_user UNIQUE (chat_id, user_id)
);

CREATE TABLE message
(
    id         SERIAL PRIMARY KEY,
    chat_id    INT REFERENCES chat (id) ON DELETE CASCADE NOT NULL,
    user_id    INT                                        NOT NULL,
    text       TEXT,
    created_at TIMESTAMP                                  NOT NULL DEFAULT NOW()
);

-- Индексы
CREATE INDEX idx_chat_user_chat_id_user_id ON chat_user (chat_id, user_id);
CREATE INDEX idx_message_chat_id ON message (chat_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_chat_user_chat_id_user_id;
DROP INDEX IF EXISTS idx_message_chat_id;

DROP TABLE IF EXISTS message;
DROP TABLE IF EXISTS chat_user;
DROP TABLE IF EXISTS chat;
-- +goose StatementEnd
