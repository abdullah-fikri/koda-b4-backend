-- Active: 1762629892012@@172.17.0.3@5432@postgres

CREATE TABLE users(
    id SERIAL PRIMARY KEY,
    email varchar(100) UNIQUE NOT NULL,
    password TEXT,
    role TEXT DEFAULT 'user',
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

CREATE TABLE profile (
    id SERIAL PRIMARY KEY,
    users_id BIGINT,
    username varchar(100),
    phone VARCHAR(20),
    address VARCHAR(100),
   FOREIGN KEY (users_id) REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);