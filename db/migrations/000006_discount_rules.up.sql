BEGIN;

CREATE TABLE IF NOT EXISTS discount_rules(
    id BIGSERIAL PRIMARY KEY,
    discount_id BIGSERIAL,
    type VARCHAR (100) NOT NULL,
    value DECIMAL (25, 2) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (discount_id) REFERENCES discounts(id)
);

COMMIT;