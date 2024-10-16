-- +goose Up
-- +goose StatementBegin
CREATE TABLE chat
(
    id         BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE chat_user
(
    id      BIGSERIAL PRIMARY KEY,
    chat_id BIGINT REFERENCES chat (id) ON DELETE CASCADE NOT NULL,
    user_id BIGINT                                        NOT NULL,
    UNIQUE (chat_id, user_id)
);

CREATE TABLE message
(
    id         BIGSERIAL PRIMARY KEY,
    chat_id    BIGINT REFERENCES chat (id) ON DELETE CASCADE NOT NULL,
    user_id    BIGINT                                        NOT NULL,
    text       VARCHAR(10000),
    created_at TIMESTAMP                                     NOT NULL DEFAULT NOW()
);

-- Индексы
CREATE INDEX idx_message_chat_id ON message (chat_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_message_chat_id;

DROP TABLE IF EXISTS message;
DROP TABLE IF EXISTS chat_user;
DROP TABLE IF EXISTS chat;
-- +goose StatementEnd
