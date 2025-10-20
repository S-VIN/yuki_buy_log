INSERT INTO users (login, password_hash) VALUES
('alice','$2a$10$mjbFCHMv1O8tXrIxTTju.euFex/plavfT875Rjsz5RWxjOunAG4QO'),
('bob','$2a$10$mjbFCHMv1O8tXrIxTTju.euFex/plavfT875Rjsz5RWxjOunAG4QO');

INSERT INTO products (name, volume, brand, default_tags) VALUES
('Tea','500ml','Brand1','healthy,drink'),
('Coffee','250g','Brand2','energy,drink');

INSERT INTO purchases (product_id, quantity, price, date, store, tags, receipt_id, user_id) VALUES
(1, 2, 100, '2023-03-01', 'Store', '{"healthy","drink","morning"}', 1, 1),
(2, 1, 200, '2023-03-02', 'Store', '{"energy","drink","work"}', 1, 2);
