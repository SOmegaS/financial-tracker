CREATE TABLE bills (
  id uuid PRIMARY KEY,
  user_id text,
  amount integer,
  category_name text,
  ts timestamp
);
