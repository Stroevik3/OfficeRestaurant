CREATE TABLE IF NOT EXISTS public.offices (
    id uuid not null primary key,
    "name" varchar not null,
    addres varchar not null,
    createdAt timestamp not null
);