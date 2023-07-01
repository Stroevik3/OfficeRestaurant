CREATE TABLE IF NOT EXISTS public.users (
    id uuid not null primary key,
    "name" varchar not null,
    officeId uuid not null,
    createdAt timestamp not null,
    FOREIGN KEY (officeId) REFERENCES offices (id)
);