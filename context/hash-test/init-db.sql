CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(255) NOT NULL
);


INSERT INTO users (username, password_hash, role) VALUES 
('admin', '$2a$10$gNkZAwY7BXT9k5K8Y4FDXOPKnwBz7ThqcaT7AtjQuVuLYtgd8vk5K', 'admin'),
('test', '$2a$10$pWxMjkSPfzfRg7gOv0OuEOWMpD3Miv9ZUApLAqhNzgYfLc13DmnXK', 'user')
ON CONFLICT (username) DO NOTHING;