-- Seed initial users
-- Password for all test users: "password123" (hashed with bcrypt)

INSERT INTO users (id, email, name, password, is_active, created_at, updated_at) VALUES
(
    uuid_generate_v4(),
    'admin@example.com',
    'Admin User',
    '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy',
    true,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
),
(
    uuid_generate_v4(),
    'user@example.com',
    'Regular User',
    '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy',
    true,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
);
