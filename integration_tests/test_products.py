# Пользователь может создать продукт
def test_create_product(req):
    user = req.get_new_user()
    product = {
        'name': 'Tea',
        'volume': '500ml',
        'brand': 'Brand1',
        'default_tags': ['drink', 'tea'],
    }
    r = req.post('products', json=product, user=user)
    assert r.status_code == 200
    assert 'id' in r.json()


# Пользователь видит свои продукты после создания
def test_get_own_products(req):
    user = req.get_new_user()
    product = {
        'name': 'Apple',
        'volume': '1kg',
        'brand': 'FreshFarm2024',
        'default_tags': ['fruit'],
    }
    r = req.post('products', json=product, user=user)
    product_id = r.json()['id']

    r = req.get('products', user=user)
    assert r.status_code == 200
    products = r.json()['products']
    assert any(p['id'] == product_id for p in products)


# Создание продукта с невалидными данными (спецсимволы в brand) должно вернуть 400
def test_create_product_invalid(req):
    user = req.get_new_user()
    product = {
        'name': 'Tea',
        'volume': '500ml',
        'brand': 'Brand!',
        'default_tags': ['drink'],
    }
    r = req.post('products', json=product, user=user)
    assert r.status_code == 400


# Пользователь может обновить свой продукт
def test_update_product_success(req):
    user = req.get_new_user()
    product = {
        'name': 'OriginalProduct',
        'volume': '500ml',
        'brand': 'OriginalBrand',
        'default_tags': ['original'],
    }
    r = req.post('products', json=product, user=user)
    assert r.status_code == 200
    product_id = r.json()['id']

    updated_product = {
        'id': product_id,
        'name': 'UpdatedProduct',
        'volume': '1L',
        'brand': 'UpdatedBrand',
        'default_tags': ['updated', 'new'],
    }
    r = req.put('products', json=updated_product, user=user)
    assert r.status_code == 200

    response_product = r.json()
    assert response_product['id'] == product_id
    assert response_product['name'] == 'UpdatedProduct'
    assert response_product['volume'] == '1L'
    assert response_product['brand'] == 'UpdatedBrand'
    assert 'updated' in response_product['default_tags']
    assert 'new' in response_product['default_tags']

    r = req.get('products', user=user)
    assert r.status_code == 200
    products = r.json()['products']
    updated = next((p for p in products if p['id'] == product_id), None)
    assert updated is not None
    assert updated['name'] == 'UpdatedProduct'
    assert updated['volume'] == '1L'


# Обновление несуществующего продукта должно вернуть 404
def test_update_product_not_found(req):
    user = req.get_new_user()
    updated_product = {
        'id': 999999,
        'name': 'UpdatedProduct',
        'volume': '1L',
        'brand': 'UpdatedBrand',
        'default_tags': ['test'],
    }
    r = req.put('products', json=updated_product, user=user)
    assert r.status_code == 404


# Обновление продукта без авторизации должно вернуть 401
def test_update_product_unauthorized(req):
    updated_product = {
        'id': 1,
        'name': 'UpdatedProduct',
        'volume': '1L',
        'brand': 'UpdatedBrand',
    }
    r = req.put('products', json=updated_product)
    assert r.status_code == 401


# Пользователь не может обновить чужой продукт
def test_update_product_different_user(req):
    user1 = req.get_new_user()
    user2 = req.get_new_user()

    product = {
        'name': 'UserOneProduct',
        'volume': '500ml',
        'brand': 'BrandOne',
    }
    r = req.post('products', json=product, user=user1)
    assert r.status_code == 200
    product_id = r.json()['id']

    updated_product = {
        'id': product_id,
        'name': 'HackedProduct',
        'volume': '1L',
        'brand': 'HackedBrand',
    }
    r = req.put('products', json=updated_product, user=user2)
    assert r.status_code == 404

    r = req.get('products', user=user1)
    assert r.status_code == 200
    products = r.json()['products']
    original = next((p for p in products if p['id'] == product_id), None)
    assert original is not None
    assert original['name'] == 'UserOneProduct'


# Обновление продукта с невалидным JSON должно вернуть 400
def test_update_product_invalid_json(req):
    user = req.get_new_user()
    r = req.put('products', data='invalid json', user=user)
    assert r.status_code == 400


# Обновление продукта без id или с id=0 должно вернуть 400
def test_update_product_missing_id(req):
    user = req.get_new_user()

    product = {
        'name': 'TestProduct',
        'volume': '500ml',
        'brand': 'TestBrand',
    }
    r = req.post('products', json=product, user=user)
    assert r.status_code == 200

    updated_product = {
        'name': 'UpdatedProduct',
        'volume': '1L',
        'brand': 'UpdatedBrand',
    }
    r = req.put('products', json=updated_product, user=user)
    assert r.status_code == 400

    updated_product['id'] = 0
    r = req.put('products', json=updated_product, user=user)
    assert r.status_code == 400


# Обновление продукта с цифрами в имени допустимо, а со спецсимволами - нет
def test_update_product_invalid_validation(req):
    user = req.get_new_user()

    product = {
        'name': 'ValidProduct',
        'volume': '500ml',
        'brand': 'ValidBrand',
    }
    r = req.post('products', json=product, user=user)
    assert r.status_code == 200
    product_id = r.json()['id']

    # Цифры в имени и бренде допустимы
    updated_product = {
        'id': product_id,
        'name': 'Valid123',
        'volume': '1L',
        'brand': 'Brand2024',
    }
    r = req.put('products', json=updated_product, user=user)
    assert r.status_code == 200, f'Expected 200 for name with digits, got {r.status_code}'

    # Спецсимволы в бренде недопустимы
    updated_product = {
        'id': product_id,
        'name': 'ValidProduct',
        'volume': '1L',
        'brand': 'Invalid@Brand',
    }
    r = req.put('products', json=updated_product, user=user)
    assert r.status_code == 400

    # Спецсимволы в имени недопустимы
    updated_product = {
        'id': product_id,
        'name': 'Product#Name',
        'volume': '1L',
        'brand': 'ValidBrand',
    }
    r = req.put('products', json=updated_product, user=user)
    assert r.status_code == 400


# Пользователь может обновить продукт с пустыми тегами
def test_update_product_empty_tags(req):
    user = req.get_new_user()

    product = {
        'name': 'ProductWithTags',
        'volume': '500ml',
        'brand': 'BrandName',
        'default_tags': ['tag1', 'tag2'],
    }
    r = req.post('products', json=product, user=user)
    assert r.status_code == 200
    product_id = r.json()['id']

    updated_product = {
        'id': product_id,
        'name': 'ProductNoTags',
        'volume': '1L',
        'brand': 'BrandName',
        'default_tags': [],
    }
    r = req.put('products', json=updated_product, user=user)
    assert r.status_code == 200

    response_product = r.json()
    assert len(response_product['default_tags']) == 0


# Пользователь может обновить один продукт несколько раз
def test_update_product_multiple_times(req):
    user = req.get_new_user()

    product = {
        'name': 'FirstVersion',
        'volume': '500ml',
        'brand': 'BrandOne',
    }
    r = req.post('products', json=product, user=user)
    assert r.status_code == 200
    product_id = r.json()['id']

    # Первое обновление
    updated_product = {
        'id': product_id,
        'name': 'SecondVersion',
        'volume': '1L',
        'brand': 'BrandTwo',
    }
    r = req.put('products', json=updated_product, user=user)
    assert r.status_code == 200
    assert r.json()['name'] == 'SecondVersion'

    # Второе обновление
    updated_product = {
        'id': product_id,
        'name': 'ThirdVersion',
        'volume': '2L',
        'brand': 'BrandThree',
    }
    r = req.put('products', json=updated_product, user=user)
    assert r.status_code == 200
    assert r.json()['name'] == 'ThirdVersion'
    assert r.json()['volume'] == '2L'

    # Проверяем финальное состояние
    r = req.get('products', user=user)
    products = r.json()['products']
    final = next((p for p in products if p['id'] == product_id), None)
    assert final is not None
    assert final['name'] == 'ThirdVersion'
    assert final['volume'] == '2L'
    assert final['brand'] == 'BrandThree'


# Пользователь может обновить продукт с 5 тегами
def test_update_product_with_many_tags(req):
    user = req.get_new_user()

    product = {
        'name': 'ProductTags',
        'volume': '500ml',
        'brand': 'BrandTags',
        'default_tags': ['initial'],
    }
    r = req.post('products', json=product, user=user)
    assert r.status_code == 200, f'Failed to create product: {r.status_code} {r.text}'
    product_id = r.json()['id']

    tags = ['apple', 'banana', 'cherry', 'date', 'elderberry']
    updated_product = {
        'id': product_id,
        'name': 'ProductManyTags',
        'volume': '1L',
        'brand': 'BrandTags',
        'default_tags': tags,
    }
    r = req.put('products', json=updated_product, user=user)
    assert r.status_code == 200, f'Failed to update product: {r.status_code} {r.text}'
    assert len(r.json()['default_tags']) == 5


# Пользователь может обновить продукт с максимальным количеством тегов (10)
def test_update_product_with_max_tags(req):
    user = req.get_new_user()

    product = {
        'name': 'MaxTagProduct',
        'volume': '500ml',
        'brand': 'MaxTagBrand',
        'default_tags': ['start'],
    }
    r = req.post('products', json=product, user=user)
    assert r.status_code == 200
    product_id = r.json()['id']

    tags = ['apple', 'banana', 'cherry', 'date', 'elderberry', 'fig', 'grape', 'honeydew', 'imbe', 'jujube']
    updated_product = {
        'id': product_id,
        'name': 'UpdatedMaxTags',
        'volume': '1L',
        'brand': 'MaxTagBrand',
        'default_tags': tags,
    }
    r = req.put('products', json=updated_product, user=user)
    assert r.status_code == 200, f'Failed to update with 10 tags: {r.status_code} {r.text}'
    assert len(r.json()['default_tags']) == 10


# Обновление продукта с больше чем 10 тегами должно вернуть 400
def test_update_product_too_many_tags(req):
    user = req.get_new_user()

    product = {
        'name': 'ProductTest',
        'volume': '500ml',
        'brand': 'BrandTest',
        'default_tags': ['initial'],
    }
    r = req.post('products', json=product, user=user)
    assert r.status_code == 200
    product_id = r.json()['id']

    tags = ['apple', 'banana', 'cherry', 'date', 'elderberry', 'fig', 'grape', 'honeydew', 'imbe', 'jujube', 'kiwi']
    updated_product = {
        'id': product_id,
        'name': 'ProductTooManyTags',
        'volume': '1L',
        'brand': 'BrandTest',
        'default_tags': tags,
    }
    r = req.put('products', json=updated_product, user=user)
    assert r.status_code == 400


# Продукты, бренды, теги и магазины могут содержать цифры и пробелы
def test_product_with_digits_and_spaces(req):
    user = req.get_new_user()

    product = {
        'name': 'Coca Cola 2024',
        'volume': '500ml',
        'brand': 'Brand123',
        'default_tags': ['tag1', 'sale2024'],
    }
    r = req.post('products', json=product, user=user)
    assert r.status_code == 200, f'Failed to create product with digits: {r.status_code} {r.text}'
    product_id = r.json()['id']

    r = req.get('products', user=user)
    products = r.json()['products']
    created_product = next(p for p in products if p['id'] == product_id)
    assert created_product['name'] == 'Coca Cola 2024'
    assert created_product['brand'] == 'Brand123'
    assert 'tag1' in created_product['default_tags']
    assert 'sale2024' in created_product['default_tags']

    # Покупка с цифрами в магазине и тегах
    purchase = {
        'product_id': product_id,
        'quantity': 3,
        'price': 299,
        'date': '2024-01-15T00:00:00Z',
        'store': 'Store 24',
        'tags': ['discount50', 'promo2024'],
        'receipt_id': 999,
    }
    r = req.post('purchases', json=purchase, user=user)
    assert r.status_code == 200, f'Failed to create purchase with digits: {r.status_code} {r.text}'
    purchase_id = r.json()['id']

    r = req.get('purchases', user=user)
    purchases = r.json()['purchases']
    created_purchase = next(p for p in purchases if p['id'] == purchase_id)
    assert created_purchase['store'] == 'Store 24'
    assert 'discount50' in created_purchase['tags']
    assert 'promo2024' in created_purchase['tags']