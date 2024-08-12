BEGIN;

CREATE TABLE IF NOT EXISTS group_roles (
    id SERIAL PRIMARY KEY,
    group_id INT NOT NULL,
    role_id INT NOT NULL,
    assigned_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (group_id) REFERENCES groups(id),
    FOREIGN KEY (role_id) REFERENCES roles(id),
    UNIQUE(group_id, role_id)
);


CREATE TRIGGER update_updated_at BEFORE
UPDATE
    ON group_roles FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

COMMIT;