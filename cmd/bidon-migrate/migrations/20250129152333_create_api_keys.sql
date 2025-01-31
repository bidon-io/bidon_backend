-- +goose Up
-- +goose StatementBegin
CREATE TABLE public.api_keys
(
    id               uuid PRIMARY KEY,
    value            varchar   NOT NULL,

    user_id          bigint    NOT NULL,

    last_accessed_at timestamp,
    created_at       timestamp NOT NULL,
    updated_at       timestamp NOT NULL,

    FOREIGN KEY (user_id) REFERENCES public.users
);
CREATE INDEX api_keys_user_id_idx ON public.api_keys (user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX public.api_keys_user_id_idx;
DROP TABLE public.api_keys;
-- +goose StatementEnd
