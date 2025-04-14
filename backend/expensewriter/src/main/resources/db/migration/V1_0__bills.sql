-- ===========================================
-- Table: category
-- ===========================================

-- ===========================================
-- Table: bills
-- ===========================================
CREATE TABLE bills
(
    id       SERIAL PRIMARY KEY,
    category TEXT             NOT NULL,
    user_id  UUID             NOT NULL,
    amount   double precision NOT NULL,
    tmstmp   timestamp        NOT NULL,
    name     TEXT             NOT NULL
);

-- ===========================================
-- Index on user_id, category_id, tmstmp
-- ===========================================
CREATE INDEX idx_bills_user_category_time
    ON bills (user_id, category, tmstmp);
