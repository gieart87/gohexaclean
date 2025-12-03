-- Seed initial users
-- Password for all test users: "password" (hashed with bcrypt)

INSERT INTO users (email, name, password, created_at, updated_at) VALUES
(
    'admin@example.com',
    'Admin User',
    '$2y$10$XDImN7MTiUaQvqvnbuD09Ok6/kgJ2bpnXQe0It21aXL3ZNil9wX..',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
),
(
    'user@example.com',
    'Regular User',
    '$2y$10$XDImN7MTiUaQvqvnbuD09Ok6/kgJ2bpnXQe0It21aXL3ZNil9wX..',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
);
