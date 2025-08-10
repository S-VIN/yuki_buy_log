INSERT INTO users (name) VALUES ('Alice'), ('Bob');

INSERT INTO products (name, volume, brand, category, description, creation_date) VALUES
('Tea','500ml','Brand1','Drink','Green tea','2023-01-01'),
('Coffee','250g','Brand2','Drink','Ground coffee','2023-02-01');

INSERT INTO purchases (product_id, quantity, price, date, store, receipt_id, user_id) VALUES
(1, 2, 100, '2023-03-01', 'Store', 1, 1),
(2, 1, 200, '2023-03-02', 'Store', 1, 2);
