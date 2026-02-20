-- +goose Up
create table users
(
    id uuid default gen_random_uuid() not null,
    username varchar(200) not null,
    email varchar(200) not null,
    password_hash text not null,
    avatar varchar(500) not null,
    skin varchar(100),
    cloak varchar(100),
    registered_at timestamp with time zone not null default current_timestamp,
    is_active boolean not null default true
);
ALTER TABLE users ADD PRIMARY KEY (id);

CREATE UNIQUE INDEX users_email_idx ON users(email);
CREATE UNIQUE INDEX users_username_idx ON users(username);

-- +goose Down
drop table if exists users;