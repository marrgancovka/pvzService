CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    email TEXT UNIQUE,
    role TEXT NOT NULL DEFAULT 'employee',
    password TEXT NOT NULL
);