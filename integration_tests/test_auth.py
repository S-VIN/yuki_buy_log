def test_register_user(db, req):
    r = req.post('register', json={'login': 'new', 'password': 'password'})
    assert r.status_code == 200, f"Registration failed: {r.status_code} {r.text}"

    token = r.json()['token']
    headers = {'Authorization': f'Bearer {token}'}

    db_result = db.execute(f'''SELECT * FROM users WHERE login='new'; ''')
    assert len(db_result) == 1
    assert db_result[0]['login'] == 'new'


def test_server_health(req):
    r = req.get(f'products')
    assert r.status_code == 401, "Server should require authentication"


def test_login_with_valid_credentials(req):
    user = req.register(login='test_login_with_valid_credentials', password='test_login_with_valid_credentials')

    user2 = req.login(login='test_login_with_valid_credentials', password='test_login_with_valid_credentials')
    assert user2.token != ''
    assert len(user2.headers) != 0


# Логин с неправильным паролем должен вернуть 401
def test_login_with_invalid_credentials(req):
    user = req.get_new_user()
    r = req.post('login', json={'login': user.login, 'password': 'wrong_password'})
    assert r.status_code == 401


# Доступ к защищённым эндпоинтам без токена должен вернуть 401
def test_access_without_token(req):
    r = req.get('products')
    assert r.status_code == 401

    r = req.get('purchases')
    assert r.status_code == 401

    r = req.get('group')
    assert r.status_code == 401

    r = req.get('invite')
    assert r.status_code == 401


# GET запрос на register и login должен вернуть 405
def test_method_not_allowed(req):
    r = req.get('register')
    assert r.status_code == 405

    r = req.get('login')
    assert r.status_code == 405
