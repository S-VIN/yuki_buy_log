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

