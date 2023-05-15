DROP TABLE IF EXISTS public.users;

CREATE TABLE public.users (
    id INTEGER NOT NULL UNIQUE,
    name TEXT,
    username TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    refresh_token TEXT,
    refresh_token_eat INTEGER,
    role TEXT NOT NULL,
    CONSTRAINT pk_users PRIMARY KEY (id)
)