CREATE TABLE IF NOT EXISTS receptions(
    id UUID PRIMARY KEY,
    date_time TIMESTAMPTZ NOT NULL DEFAULT now(),
    pvz_id UUID REFERENCES pvz(id) ON DELETE SET NULL,
    status TEXT NOT NULL DEFAULT 'in_progress'
);
