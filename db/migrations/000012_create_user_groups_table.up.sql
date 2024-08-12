BEGIN;

CREATE TABLE IF NOT EXISTS user_groups (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    group_id INT NOT NULL,
    assigned_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (group_id) REFERENCES groups(id),
    UNIQUE(user_id, group_id)
);

CREATE TRIGGER update_updated_at BEFORE
UPDATE
    ON user_groups FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

COMMIT;