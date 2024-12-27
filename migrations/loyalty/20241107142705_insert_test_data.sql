-- +goose Up
-- +goose StatementBegin
INSERT INTO public.loyalty (username, reservation_count, status, discount)
VALUES ('Test Max', 25, 'GOLD', 10);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
