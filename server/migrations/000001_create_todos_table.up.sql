CREATE TABLE IF NOT EXISTS todos (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    todo text NOT NULL,
    description text NOT NULL
);