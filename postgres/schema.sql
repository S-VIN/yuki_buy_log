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
    category VARCHAR(30) NOT NULL,
    description VARCHAR(150),
    creation_date DATE NOT NULL
);

CREATE TABLE purchases (
    id SERIAL PRIMARY KEY,
    product_id INTEGER NOT NULL REFERENCES products(id),
    quantity INTEGER NOT NULL CHECK (quantity >= 1 AND quantity <= 100000),
    price INTEGER NOT NULL CHECK (price >= 1 AND price <= 100000000),
    date DATE NOT NULL,
    store VARCHAR(30) NOT NULL,
    receipt_id INTEGER,
    user_id INTEGER REFERENCES users(id)
);

