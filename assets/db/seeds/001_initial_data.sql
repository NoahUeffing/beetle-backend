-- +goose Up
INSERT INTO tasks (title, description, status, created_at, updated_at)
VALUES 
    ('Example Task 1', 'This is an example task', 'pending', NOW(), NOW()),
    ('Example Task 2', 'Another example task', 'in_progress', NOW(), NOW());

-- +goose Down
DELETE FROM tasks WHERE title IN ('Example Task 1', 'Example Task 2'); 