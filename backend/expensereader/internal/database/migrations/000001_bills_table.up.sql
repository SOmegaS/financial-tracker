CREATE TABLE bills
(
    id       SERIAL PRIMARY KEY,
    category TEXT             NOT NULL,
    user_id  UUID             NOT NULL,
    amount   double precision NOT NULL,
    tmstmp   timestamp        NOT NULL,
    name     TEXT             NOT NULL
);
