CREATE INDEX idx_bills_user_category_time
    ON bills (user_id, category, tmstmp);
