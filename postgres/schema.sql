DROP TABLE IF EXISTS invites CASCADE;
DROP TABLE IF EXISTS groups CASCADE;
DROP TABLE IF EXISTS purchases CASCADE;
DROP TABLE IF EXISTS products CASCADE;
DROP TABLE IF EXISTS users CASCADE;

-- Create tables
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    login VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(60) NOT NULL
);

CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(30) NOT NULL,
    volume VARCHAR(10) NOT NULL,
    brand VARCHAR(30) NOT NULL,
    default_tags VARCHAR(30) NOT NULL,
    user_id INTEGER REFERENCES users(id)
);

CREATE TABLE purchases (
    id SERIAL PRIMARY KEY,
    product_id INTEGER NOT NULL REFERENCES products(id),
    quantity INTEGER NOT NULL,
    price INTEGER NOT NULL,
    date DATE NOT NULL,
    store VARCHAR(30) NOT NULL,
    tags TEXT[],
    receipt_id INTEGER,
    user_id INTEGER REFERENCES users(id)
);

CREATE TABLE groups (
    id SERIAL,
    user_id INTEGER NOT NULL REFERENCES users(id),
    PRIMARY KEY (id, user_id),
    UNIQUE (user_id)
);

CREATE TABLE invites (
    id SERIAL PRIMARY KEY,
    from_user_id INTEGER NOT NULL REFERENCES users(id),
    to_user_id INTEGER NOT NULL REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (from_user_id, to_user_id),
    CHECK (from_user_id != to_user_id)
);

