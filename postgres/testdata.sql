INSERT INTO users (login, password_hash) VALUES
('alice','$2a$10$mjbFCHMv1O8tXrIxTTju.euFex/plavfT875Rjsz5RWxjOunAG4QO'),
('bob','$2a$10$mjbFCHMv1O8tXrIxTTju.euFex/plavfT875Rjsz5RWxjOunAG4QO');

INSERT INTO products (name, volume, brand, category, description, creation_date, user_login) VALUES
('Tea','500ml','Brand1','Drink','Green tea','2023-01-01','alice'),
('Coffee','250g','Brand2','Drink','Ground coffee','2023-02-01','bob');

INSERT INTO purchases (product_id, quantity, price, date, store, receipt_id, user_login) VALUES
(1, 2, 100, '2023-03-01', 'Store', 1, 'alice'),
(2, 1, 200, '2023-03-02', 'Store', 1, 'bob');
