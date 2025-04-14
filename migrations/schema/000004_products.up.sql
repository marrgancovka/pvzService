CREATE TABLE IF NOT EXISTS products (
    id UUID PRIMARY KEY,
    date_time TIMESTAMPTZ NOT NULL DEFAULT now(),
    type TEXT NOT NULL,
    reception_id UUID REFERENCES receptions(id) ON DELETE SET NULL
)
