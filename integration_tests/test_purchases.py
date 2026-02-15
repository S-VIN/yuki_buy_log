# Пользователь может создать покупку
def test_create_purchase(req):
    user = req.get_new_user()

    product = {
        'name': 'Milk',
        'volume': '1L',
        'brand': 'DairyBest',
    }
    r = req.post('products', json=product, user=user)
    product_id = r.json()['id']

    purchase = {
        'product_id': product_id,
        'quantity': 2,
        'price': 150,
        'date': '2024-01-15T00:00:00Z',
        'store': 'Store 7',
        'receipt_id': 1,
    }
    r = req.post('purchases', json=purchase, user=user)
    assert r.status_code == 200
    assert 'id' in r.json()


# Пользователь видит свои покупки после создания
def test_get_own_purchases(req):
    user = req.get_new_user()

    product = {
        'name': 'Bread',
        'volume': '1pc',
        'brand': 'Bakery',
    }
    r = req.post('products', json=product, user=user)
    product_id = r.json()['id']

    purchase = {
        'product_id': product_id,
        'quantity': 1,
        'price': 50,
        'date': '2024-01-20T00:00:00Z',
        'store': 'LocalStore',
        'receipt_id': 100,
    }
    r = req.post('purchases', json=purchase, user=user)
    purchase_id = r.json()['id']

    r = req.get('purchases', user=user)
    assert r.status_code == 200
    purchases = r.json()['purchases']
    assert any(p['id'] == purchase_id for p in purchases)


# Создание покупки с невалидными данными (quantity=0) должно вернуть 400
def test_create_purchase_invalid(req):
    user = req.get_new_user()

    product = {
        'name': 'Coffee',
        'volume': '250g',
        'brand': 'CoffeeCo',
    }
    r = req.post('products', json=product, user=user)
    product_id = r.json()['id']

    purchase = {
        'product_id': product_id,
        'quantity': 0,
        'price': 100,
        'date': '2024-01-01',
        'store': 'Store',
        'receipt_id': 1,
    }
    r = req.post('purchases', json=purchase, user=user)
    assert r.status_code == 400


# Пользователь может удалить свою покупку
def test_delete_purchase_success(req):
    user = req.get_new_user()

    product = {
        'name': 'Sugar',
        'volume': '1kg',
        'brand': 'SweetCo',
    }
    r = req.post('products', json=product, user=user)
    product_id = r.json()['id']

    purchase = {
        'product_id': product_id,
        'quantity': 1,
        'price': 100,
        'date': '2024-01-25T00:00:00Z',
        'store': 'Supermarket',
        'receipt_id': 500,
    }
    r = req.post('purchases', json=purchase, user=user)
    assert r.status_code == 200
    purchase_id = r.json()['id']

    r = req.delete('purchases', json={'id': purchase_id}, user=user)
    assert r.status_code == 204

    r = req.get('purchases', user=user)
    assert r.status_code == 200
    purchases = r.json()['purchases']
    assert not any(p['id'] == purchase_id for p in purchases)


# Удаление несуществующей покупки должно вернуть 404
def test_delete_purchase_not_found(req):
    user = req.get_new_user()
    r = req.delete('purchases', json={'id': 999999}, user=user)
    assert r.status_code == 404


# Удаление покупки без авторизации должно вернуть 401
def test_delete_purchase_unauthorized(req):
    r = req.delete('purchases', json={'id': 1})
    assert r.status_code == 401


# Пользователь не может удалить чужую покупку
def test_delete_purchase_different_user(req):
    user1 = req.get_new_user()
    user2 = req.get_new_user()

    product = {
        'name': 'Salt',
        'volume': '500g',
        'brand': 'SeaSalt',
    }
    r = req.post('products', json=product, user=user1)
    product_id = r.json()['id']

    purchase = {
        'product_id': product_id,
        'quantity': 1,
        'price': 50,
        'date': '2024-01-26T00:00:00Z',
        'store': 'Store',
        'receipt_id': 600,
    }
    r = req.post('purchases', json=purchase, user=user1)
    assert r.status_code == 200
    purchase_id = r.json()['id']

    # user2 пытается удалить покупку user1
    r = req.delete('purchases', json={'id': purchase_id}, user=user2)
    assert r.status_code == 404

    # Покупка user1 всё ещё существует
    r = req.get('purchases', user=user1)
    assert r.status_code == 200
    purchases = r.json()['purchases']
    assert any(p['id'] == purchase_id for p in purchases)


# Удаление покупки с невалидным JSON должно вернуть 400
def test_delete_purchase_invalid_json(req):
    user = req.get_new_user()
    r = req.delete('purchases', data='invalid json', user=user)
    assert r.status_code == 400


# Удаление покупки без id или с id=0 должно вернуть 400
def test_delete_purchase_missing_id(req):
    user = req.get_new_user()

    r = req.delete('purchases', json={}, user=user)
    assert r.status_code == 400

    r = req.delete('purchases', json={'id': 0}, user=user)
    assert r.status_code == 400


# Пользователь может последовательно удалить несколько покупок
def test_delete_multiple_purchases(req):
    user = req.get_new_user()

    product = {
        'name': 'Flour',
        'volume': '2kg',
        'brand': 'BakerCo',
    }
    r = req.post('products', json=product, user=user)
    product_id = r.json()['id']

    purchase_ids = []
    stores = ['StoreOne', 'StoreTwo', 'StoreThree']
    for i in range(3):
        purchase = {
            'product_id': product_id,
            'quantity': i + 1,
            'price': 200 + (i * 50),
            'date': f'2024-01-{27 + i:02d}T00:00:00Z',
            'store': stores[i],
            'receipt_id': 700 + i,
        }
        r = req.post('purchases', json=purchase, user=user)
        assert r.status_code == 200
        purchase_ids.append(r.json()['id'])

    for purchase_id in purchase_ids:
        r = req.delete('purchases', json={'id': purchase_id}, user=user)
        assert r.status_code == 204

    r = req.get('purchases', user=user)
    assert r.status_code == 200
    purchases = r.json()['purchases']
    for purchase_id in purchase_ids:
        assert not any(p['id'] == purchase_id for p in purchases)