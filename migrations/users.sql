DROP TABLE IF EXISTS public.users;

CREATE TABLE public.users (
    id SERIAL PRIMARY KEY,
    name TEXT,
    username TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    role TEXT NOT NULL
)