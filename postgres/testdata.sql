-- Users (password hash = 'password' for all)
INSERT INTO users (login, password_hash) VALUES
('alice', '$2a$10$mjbFCHMv1O8tXrIxTTju.euFex/plavfT875Rjsz5RWxjOunAG4QO'),
('bob',   '$2a$10$mjbFCHMv1O8tXrIxTTju.euFex/plavfT875Rjsz5RWxjOunAG4QO'),
('charlie','$2a$10$mjbFCHMv1O8tXrIxTTju.euFex/plavfT875Rjsz5RWxjOunAG4QO');
-- alice=1, bob=2, charlie=3

-- Products
INSERT INTO products (name, volume, brand, default_tags, user_id) VALUES
('Tea',         '500ml', 'BrandA', ARRAY['healthy','drink'],         1),
('Coffee',      '250g',  'BrandB', ARRAY['energy','drink'],          1),
('Juice',       '1L',    'BrandC', ARRAY['healthy','fruit'],         1),
('Milk',        '1L',    'BrandD', ARRAY['healthy','dairy'],         2),
('Protein Bar', '60g',   'BrandE', ARRAY['sport','snack'],           2),
('Water',       '0.5L',  'BrandF', ARRAY['drink'],                   3);
-- tea=1, coffee=2, juice=3, milk=4, protein bar=5, water=6

-- Groups: alice(1) + bob(2) in group 1, charlie(3) alone in group 2
INSERT INTO group_members (group_id, user_id, member_number) VALUES
(1, 1, 1),
(1, 2, 2),
(2, 3, 1);

-- Purchases (receipt_id 1 and 2 are alice's receipts, 3 is bob's)
INSERT INTO purchases (product_id, quantity, price, date, store, tags, receipt_id, user_id) VALUES
(1, 3, 150,  '2025-01-05', 'Magnit',   ARRAY['healthy','drink','morning'],  1, 1),
(2, 1, 220,  '2025-01-05', 'Magnit',   ARRAY['energy','drink','morning'],   1, 1),
(3, 2, 180,  '2025-01-12', 'Lenta',    ARRAY['healthy','fruit'],            2, 1),
(4, 2, 130,  '2025-01-08', 'Perekrestok', ARRAY['healthy','dairy'],         3, 2),
(5, 4, 400,  '2025-01-15', 'SportMaster', ARRAY['sport','snack'],           3, 2),
(1, 1, 155,  '2025-02-01', 'Magnit',   ARRAY['drink'],                     NULL, 1),
(6, 6, 90,   '2025-02-03', 'Vkusvill', ARRAY['drink'],                     NULL, 3),
(2, 2, 440,  '2025-02-10', 'Lenta',    ARRAY['energy','drink'],            NULL, 2),
(4, 1, 135,  '2025-02-14', 'Perekrestok', ARRAY['healthy','dairy'],        NULL, 1),
(3, 3, 270,  '2025-03-01', 'Vkusvill', ARRAY['healthy','fruit'],           NULL, 3);

-- Invites: charlie invites alice to join his group
INSERT INTO invites (from_user_id, to_user_id) VALUES
(3, 1);