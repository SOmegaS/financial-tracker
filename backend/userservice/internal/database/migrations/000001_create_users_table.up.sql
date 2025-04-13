create table users
(
    id       uuid PRIMARY KEY,
    username text not null unique,
    password text not null
);
