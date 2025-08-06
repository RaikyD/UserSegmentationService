-- +goose Up
-- +goose StatementBegin
CREATE TABLE user_segment_assignment (
    segment_id      UUID         NOT NULL
     REFERENCES segments(id)
         ON DELETE CASCADE,
    user_id         UUID         NOT NULL,
    assignment_type TEXT         NOT NULL
     CHECK (assignment_type IN ('manual','auto')),
    assigned_at     TIMESTAMPTZ  NOT NULL DEFAULT now(),
    PRIMARY KEY (segment_id, user_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_segment_assignment;
-- +goose StatementEnd
