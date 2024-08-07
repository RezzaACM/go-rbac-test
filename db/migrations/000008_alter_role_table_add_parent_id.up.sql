BEGIN;

-- Step 1: Alter the table to add the parent_id column
ALTER TABLE
    roles
ADD
    COLUMN parent_id INT;

-- Step 2: Add the foreign key constraint to reference role_id
ALTER TABLE
    roles
ADD
    CONSTRAINT fk_roles_parent FOREIGN KEY (parent_id) REFERENCES roles(id);

COMMIT;