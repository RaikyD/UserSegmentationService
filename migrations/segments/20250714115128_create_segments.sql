-- +goose Up
-- +goose StatementBegin
CREATE TABLE segments (
  id           UUID        PRIMARY KEY,
  segment_name TEXT        NOT NULL UNIQUE,
  type         TEXT        NOT NULL
      CHECK (type IN ('static','dynamic','dynamic_rule')),
  config       JSONB       NOT NULL,
  description  TEXT,
  is_active    BOOLEAN     NOT NULL DEFAULT true,
  created_on   TIMESTAMPTZ NOT NULL DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS segments;
-- +goose StatementEnd
