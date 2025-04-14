create table users
(
    id       uuid PRIMARY KEY,
    username text not null,
    password text not null
);

create index users_user_id_idx on users (id)