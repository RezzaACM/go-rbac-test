BEGIN;

-- Step 1: Remove the foreign key constraint
ALTER TABLE
    roles DROP CONSTRAINT fk_roles_parent;

-- Step 2: Drop the parent_id column
ALTER TABLE
    roles DROP COLUMN parent_id;

COMMIT;