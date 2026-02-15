# Новый пользователь не должен быть в группе
def test_initial_group_empty(req):
    user = req.get_new_user()
    r = req.get('group', user=user)
    assert r.status_code == 200
    assert r.json()['members'] == []


# Пользователь может отправить инвайт другому пользователю
def test_send_invite(req):
    user1 = req.get_new_user()
    user2 = req.get_new_user()

    r = req.post('invite', json={'login': user2.login}, user=user1)
    assert r.status_code == 200
    assert 'invite_id' in r.json()


# Пользователь видит входящие инвайты
def test_receive_invite(req):
    user1 = req.get_new_user()
    user2 = req.get_new_user()

    r = req.post('invite', json={'login': user2.login}, user=user1)
    assert r.status_code == 200

    r = req.get('invite', user=user2)
    assert r.status_code == 200
    invites = r.json()['invites']
    assert any(inv['from_login'] == user1.login for inv in invites)


# Взаимные инвайты создают группу
def test_mutual_invite_creates_group(req):
    user1 = req.get_new_user()
    user2 = req.get_new_user()

    r = req.post('invite', json={'login': user2.login}, user=user1)
    assert r.status_code == 200

    r = req.post('invite', json={'login': user1.login}, user=user2)
    assert r.status_code == 200
    data = r.json()
    assert 'mutual_invite' in data or 'group' in data.get('message', '').lower()

    r = req.get('group', user=user1)
    assert r.status_code == 200
    members = r.json()['members']
    assert len(members) == 2
    member_logins = [m['login'] for m in members]
    assert user1.login in member_logins
    assert user2.login in member_logins


# Инвайт несуществующему пользователю должен вернуть 404
def test_invite_to_nonexistent_user(req):
    user = req.get_new_user()
    r = req.post('invite', json={'login': 'nonexistent_user_xyz'}, user=user)
    assert r.status_code == 404


# До создания группы пользователи не видят продукты друг друга
def test_products_isolation_before_group(req):
    user1 = req.get_new_user()
    user2 = req.get_new_user()

    product1 = {
        'name': 'ProductOne',
        'volume': '500ml',
        'brand': 'BrandOne',
        'default_tags': ['tag1'],
    }
    r = req.post('products', json=product1, user=user1)
    assert r.status_code == 200

    product2 = {
        'name': 'ProductTwo',
        'volume': '1L',
        'brand': 'BrandTwo',
        'default_tags': ['tag2'],
    }
    r = req.post('products', json=product2, user=user2)
    assert r.status_code == 200

    r = req.get('products', user=user1)
    assert r.status_code == 200
    products = r.json()['products']
    assert len(products) == 1
    assert products[0]['name'] == 'ProductOne'


# После создания группы пользователи видят продукты друг друга
def test_products_shared_in_group(req):
    user1 = req.get_new_user()
    user2 = req.get_new_user()

    product1 = {
        'name': 'SharedProductOne',
        'volume': '500ml',
        'brand': 'BrandOne',
    }
    r = req.post('products', json=product1, user=user1)
    product1_id = r.json()['id']

    product2 = {
        'name': 'SharedProductTwo',
        'volume': '1L',
        'brand': 'BrandTwo',
    }
    r = req.post('products', json=product2, user=user2)
    product2_id = r.json()['id']

    # Создаём группу через взаимные инвайты
    req.post('invite', json={'login': user2.login}, user=user1)
    req.post('invite', json={'login': user1.login}, user=user2)

    r = req.get('products', user=user1)
    assert r.status_code == 200
    products = r.json()['products']
    assert len(products) == 2
    product_ids = [p['id'] for p in products]
    assert product1_id in product_ids
    assert product2_id in product_ids
    assert all('user_id' in p for p in products)


# После создания группы пользователи видят покупки друг друга
def test_purchases_shared_in_group(req):
    user1 = req.get_new_user()
    user2 = req.get_new_user()

    product1 = {'name': 'PurchProdOne', 'volume': '1L', 'brand': 'BrandA'}
    r = req.post('products', json=product1, user=user1)
    product1_id = r.json()['id']

    product2 = {'name': 'PurchProdTwo', 'volume': '1L', 'brand': 'BrandB'}
    r = req.post('products', json=product2, user=user2)
    product2_id = r.json()['id']

    # Создаём группу
    req.post('invite', json={'login': user2.login}, user=user1)
    req.post('invite', json={'login': user1.login}, user=user2)

    purchase1 = {
        'product_id': product1_id,
        'quantity': 2,
        'price': 1000,
        'date': '2024-01-01T00:00:00Z',
        'store': 'StoreOne',
        'tags': ['buy2024'],
        'receipt_id': 100,
    }
    r = req.post('purchases', json=purchase1, user=user1)
    purchase1_id = r.json()['id']

    purchase2 = {
        'product_id': product2_id,
        'quantity': 3,
        'price': 2000,
        'date': '2024-01-02T00:00:00Z',
        'store': 'Store2',
        'tags': ['sale'],
        'receipt_id': 200,
    }
    r = req.post('purchases', json=purchase2, user=user2)
    purchase2_id = r.json()['id']

    r = req.get('purchases', user=user1)
    assert r.status_code == 200
    purchases = r.json()['purchases']
    assert len(purchases) == 2
    purchase_ids = [p['id'] for p in purchases]
    assert purchase1_id in purchase_ids
    assert purchase2_id in purchase_ids
    assert all('user_id' in p for p in purchases)


# Пользователь вне группы не видит данные группы
def test_isolation_for_non_group_users(req):
    user1 = req.get_new_user()
    user2 = req.get_new_user()
    user3 = req.get_new_user()

    # user1 и user2 создают группу
    req.post('invite', json={'login': user2.login}, user=user1)
    req.post('invite', json={'login': user1.login}, user=user2)

    product = {'name': 'GroupProduct', 'volume': '1L', 'brand': 'BrandX'}
    req.post('products', json=product, user=user1)

    # user3 не в группе и не должен видеть продукты
    r = req.get('products', user=user3)
    assert r.status_code == 200
    products = r.json()['products']
    assert len(products) == 0


# Пользователь может выйти из группы
def test_leave_group(req):
    user1 = req.get_new_user()
    user2 = req.get_new_user()

    req.post('invite', json={'login': user2.login}, user=user1)
    req.post('invite', json={'login': user1.login}, user=user2)

    r = req.delete('group', user=user2)
    assert r.status_code == 200

    r = req.get('group', user=user2)
    assert r.status_code == 200
    assert r.json()['members'] == []


# Группа автоматически удаляется, когда остаётся 1 участник
def test_group_auto_deletion(req):
    user1 = req.get_new_user()
    user2 = req.get_new_user()

    req.post('invite', json={'login': user2.login}, user=user1)
    req.post('invite', json={'login': user1.login}, user=user2)

    r = req.delete('group', user=user2)
    assert r.status_code == 200

    r = req.get('group', user=user1)
    assert r.status_code == 200
    assert r.json()['members'] == []


# После выхода из группы пользователи не видят данные друг друга
def test_data_isolation_after_leaving_group(req):
    user1 = req.get_new_user()
    user2 = req.get_new_user()

    product1 = {'name': 'ProductOne', 'volume': '1L', 'brand': 'BrandA'}
    r = req.post('products', json=product1, user=user1)
    product1_id = r.json()['id']

    product2 = {'name': 'ProductTwo', 'volume': '1L', 'brand': 'BrandB'}
    r = req.post('products', json=product2, user=user2)

    # Создаём группу
    req.post('invite', json={'login': user2.login}, user=user1)
    req.post('invite', json={'login': user1.login}, user=user2)

    r = req.get('products', user=user1)
    assert len(r.json()['products']) == 2

    # user2 выходит из группы
    req.delete('group', user=user2)

    r = req.get('products', user=user1)
    assert r.status_code == 200
    products = r.json()['products']
    assert len(products) == 1
    assert products[0]['id'] == product1_id


# Выход из группы, когда не состоишь в группе, должен вернуть 400
def test_leave_group_when_not_in_group(req):
    user = req.get_new_user()
    r = req.delete('group', user=user)
    assert r.status_code == 400


# Свободный пользователь может присоединиться к существующей группе через взаимные инвайты
def test_can_invite_to_expand_group(req):
    user1 = req.get_new_user()
    user2 = req.get_new_user()
    user3 = req.get_new_user()

    # user1 и user2 создают группу
    req.post('invite', json={'login': user2.login}, user=user1)
    req.post('invite', json={'login': user1.login}, user=user2)

    r = req.get('group', user=user1)
    assert len(r.json()['members']) == 2

    # user3 отправляет инвайт user1 (в группе)
    r = req.post('invite', json={'login': user1.login}, user=user3)
    assert r.status_code == 200
    assert 'invite_id' in r.json()

    # user1 отправляет инвайт user3 — взаимный инвайт расширяет группу
    r = req.post('invite', json={'login': user3.login}, user=user1)
    assert r.status_code == 200
    data = r.json()
    assert 'mutual_invite' in data or 'group' in data.get('message', '').lower()

    # Проверяем, что в группе 3 участника
    r = req.get('group', user=user1)
    members = r.json()['members']
    assert len(members) == 3
    member_logins = [m['login'] for m in members]
    assert user1.login in member_logins
    assert user2.login in member_logins
    assert user3.login in member_logins

    r = req.get('group', user=user3)
    assert len(r.json()['members']) == 3


# Пользователи из разных групп не могут отправлять инвайты друг другу
def test_cannot_invite_users_in_different_groups(req):
    user1 = req.get_new_user()
    user2 = req.get_new_user()
    user3 = req.get_new_user()
    user4 = req.get_new_user()

    # Группа 1: user1 и user2
    req.post('invite', json={'login': user2.login}, user=user1)
    req.post('invite', json={'login': user1.login}, user=user2)

    # Группа 2: user3 и user4
    req.post('invite', json={'login': user4.login}, user=user3)
    req.post('invite', json={'login': user3.login}, user=user4)

    r = req.get('group', user=user1)
    assert len(r.json()['members']) == 2

    r = req.get('group', user=user3)
    assert len(r.json()['members']) == 2

    # user1 (группа 1) пытается пригласить user3 (группа 2)
    r = req.post('invite', json={'login': user3.login}, user=user1)
    assert r.status_code == 400
    assert 'different groups' in r.text.lower()

    # user3 (группа 2) пытается пригласить user1 (группа 1)
    r = req.post('invite', json={'login': user1.login}, user=user3)
    assert r.status_code == 400
    assert 'different groups' in r.text.lower()


# Полный цикл: регистрация, создание данных, группа, шаринг, выход
def test_complete_sharing_flow(req):
    user1 = req.get_new_user()
    user2 = req.get_new_user()

    # Каждый создаёт продукт
    r = req.post('products', json={'name': 'ProductA', 'volume': '1L', 'brand': 'BrandA'}, user=user1)
    assert r.status_code == 200

    r = req.post('products', json={'name': 'ProductB', 'volume': '1L', 'brand': 'BrandB'}, user=user2)
    assert r.status_code == 200

    # Изоляция до группы
    r = req.get('products', user=user1)
    assert len(r.json()['products']) == 1

    # Создаём группу
    req.post('invite', json={'login': user2.login}, user=user1)
    req.post('invite', json={'login': user1.login}, user=user2)

    # Общий доступ
    r = req.get('products', user=user1)
    assert len(r.json()['products']) == 2
    r = req.get('products', user=user2)
    assert len(r.json()['products']) == 2

    # user2 выходит
    req.delete('group', user=user2)

    # Группа удалена
    r = req.get('group', user=user1)
    assert r.json()['members'] == []

    # Изоляция после выхода
    r = req.get('products', user=user1)
    products = r.json()['products']
    assert len(products) == 1


# При создании группы участникам назначаются номера 1 и 2
def test_member_numbers_assigned_on_group_creation(req):
    user1 = req.get_new_user()
    user2 = req.get_new_user()

    req.post('invite', json={'login': user2.login}, user=user1)
    req.post('invite', json={'login': user1.login}, user=user2)

    r = req.get('group', user=user1)
    assert r.status_code == 200
    members = r.json()['members']
    assert len(members) == 2

    member_numbers = sorted([m['member_number'] for m in members])
    assert member_numbers == [1, 2]
    assert len(set(member_numbers)) == 2


# При расширении группы номера участников последовательные
def test_member_numbers_sequential_on_expansion(req):
    user1 = req.get_new_user()
    user2 = req.get_new_user()
    user3 = req.get_new_user()

    req.post('invite', json={'login': user2.login}, user=user1)
    req.post('invite', json={'login': user1.login}, user=user2)

    req.post('invite', json={'login': user3.login}, user=user1)
    req.post('invite', json={'login': user1.login}, user=user3)

    r = req.get('group', user=user1)
    assert r.status_code == 200
    members = r.json()['members']
    assert len(members) == 3

    member_numbers = sorted([m['member_number'] for m in members])
    assert member_numbers == [1, 2, 3]


# После выхода участника номера перенумеровываются
def test_member_numbers_renumbered_after_leave(req):
    user1 = req.get_new_user()
    user2 = req.get_new_user()
    user3 = req.get_new_user()

    # Создаём группу из 3 участников
    req.post('invite', json={'login': user2.login}, user=user1)
    req.post('invite', json={'login': user1.login}, user=user2)
    req.post('invite', json={'login': user3.login}, user=user1)
    req.post('invite', json={'login': user1.login}, user=user3)

    r = req.get('group', user=user1)
    members = r.json()['members']
    assert len(members) == 3

    # Находим участника с номером 2 и выводим его из группы
    member2 = next((m for m in members if m['member_number'] == 2), None)
    assert member2 is not None

    users_map = {user1.login: user1, user2.login: user2, user3.login: user3}
    leaving_user = users_map[member2['login']]
    remaining_users = [u for u in [user1, user2, user3] if u.login != leaving_user.login]

    r = req.delete('group', user=leaving_user)
    assert r.status_code == 200

    # Оставшиеся участники перенумерованы как 1 и 2
    r = req.get('group', user=remaining_users[0])
    assert r.status_code == 200
    remaining_members = r.json()['members']
    assert len(remaining_members) == 2

    member_numbers = sorted([m['member_number'] for m in remaining_members])
    assert member_numbers == [1, 2], f'Expected [1, 2] but got {member_numbers}'


# После выхода участника идентичности оставшихся участников сохраняются
def test_member_numbers_preserved_for_non_leaving_members(req):
    user1 = req.get_new_user()
    user2 = req.get_new_user()
    user3 = req.get_new_user()
    user4 = req.get_new_user()

    # Создаём группу из 4 участников
    req.post('invite', json={'login': user2.login}, user=user1)
    req.post('invite', json={'login': user1.login}, user=user2)
    req.post('invite', json={'login': user3.login}, user=user1)
    req.post('invite', json={'login': user1.login}, user=user3)
    req.post('invite', json={'login': user4.login}, user=user1)
    req.post('invite', json={'login': user1.login}, user=user4)

    r = req.get('group', user=user1)
    members = r.json()['members']
    assert len(members) == 4

    # Находим участника с номером 2 и выводим его
    member2 = next((m for m in members if m['member_number'] == 2), None)
    leaving_login = member2['login']

    users_map = {
        user1.login: user1,
        user2.login: user2,
        user3.login: user3,
        user4.login: user4,
    }
    leaving_user = users_map[leaving_login]
    remaining_users = [u for u in [user1, user2, user3, user4] if u.login != leaving_login]

    req.delete('group', user=leaving_user)

    r = req.get('group', user=remaining_users[0])
    remaining_members = r.json()['members']
    assert len(remaining_members) == 3

    # Все оставшиеся участники на месте
    remaining_logins = {m['login'] for m in remaining_members}
    expected_logins = {user1.login, user2.login, user3.login, user4.login} - {leaving_login}
    assert remaining_logins == expected_logins

    # Номера последовательные
    member_numbers = sorted([m['member_number'] for m in remaining_members])
    assert member_numbers == [1, 2, 3]


# Номера участников всегда между 1 и 5
def test_member_number_in_range(req):
    user1 = req.get_new_user()
    user2 = req.get_new_user()

    req.post('invite', json={'login': user2.login}, user=user1)
    req.post('invite', json={'login': user1.login}, user=user2)

    r = req.get('group', user=user1)
    members = r.json()['members']

    for member in members:
        assert 'member_number' in member
        assert 1 <= member['member_number'] <= 5, f'Member number {member["member_number"]} is out of range'