import os
import subprocess
import time
import uuid

import pytest
import requests

BASE_URL = "http://localhost:8080"
COMPOSE_FILE = os.path.join(os.path.dirname(__file__), "..", "docker-compose.yml")


def detect_compose():
    for cmd in (["docker", "compose"], ["docker-compose"]):
        try:
            if subprocess.run(cmd + ["version"], capture_output=True).returncode == 0:
                return cmd
        except FileNotFoundError:
            continue
    return None


COMPOSE_CMD = detect_compose()


def run_command(cmd):
    result = subprocess.run(cmd, capture_output=True, text=True)
    if result.returncode != 0:
        raise RuntimeError(
            f"{' '.join(cmd)} failed with code {result.returncode}\n"
            f"stdout:\n{result.stdout}\n"
            f"stderr:\n{result.stderr}"
        )
    return result


@pytest.fixture(scope="session", autouse=True)
def start_system():
    if COMPOSE_CMD is None:
        pytest.skip("docker compose not available")
    try:
        run_command(COMPOSE_CMD + ["-f", COMPOSE_FILE, "up", "--build", "-d"])
    except (FileNotFoundError, RuntimeError) as e:
        pytest.skip(str(e))
    for _ in range(30):
        try:
            requests.get(BASE_URL, timeout=1)
            break
        except Exception:
            time.sleep(1)
    yield
    try:
        run_command(COMPOSE_CMD + ["-f", COMPOSE_FILE, "down", "-v"])
    except Exception:
        pass


def register_and_login():
    login = f"user_{uuid.uuid4().hex[:8]}"
    password = "password"
    r = requests.post(f"{BASE_URL}/register", json={"login": login, "password": password})
    assert r.status_code == 200
    r = requests.post(f"{BASE_URL}/login", json={"login": login, "password": password})
    assert r.status_code == 200
    token = r.json()["token"]
    headers = {"Authorization": f"Bearer {token}"}
    return login, password, headers


def test_full_flow():
    _, _, headers = register_and_login()

    product_payload = {
        "name": "Tea",
        "volume": "500ml",
        "brand": "Brand1",
        "category": "Drink",
        "description": "Green tea",
    }
    r = requests.post(f"{BASE_URL}/products", json=product_payload, headers=headers)
    assert r.status_code == 200
    product_id = r.json()["id"]

    r = requests.get(f"{BASE_URL}/products", headers=headers)
    assert r.status_code == 200
    products = r.json()["products"]
    assert any(p["id"] == product_id for p in products)

    purchase_payload = {
        "product_id": product_id,
        "quantity": 1,
        "price": 100,
        "date": "2024-01-01",
        "store": "Store",
        "receipt_id": 1,
    }
    r = requests.post(f"{BASE_URL}/purchases", json=purchase_payload, headers=headers)
    assert r.status_code == 200
    purchase_id = r.json()["id"]

    r = requests.get(f"{BASE_URL}/purchases", headers=headers)
    assert r.status_code == 200
    purchases = r.json()["purchases"]
    assert any(p["id"] == purchase_id for p in purchases)


def test_login_bad_credentials():
    login, _, _ = register_and_login()
    r = requests.post(f"{BASE_URL}/login", json={"login": login, "password": "wrong"})
    assert r.status_code == 401


def test_access_without_token():
    r = requests.get(f"{BASE_URL}/products")
    assert r.status_code == 401
    r = requests.get(f"{BASE_URL}/purchases")
    assert r.status_code == 401


def test_create_product_invalid():
    _, _, headers = register_and_login()
    payload = {
        "name": "Tea",
        "volume": "500ml",
        "brand": "Brand!",
        "category": "Drink",
        "description": "Green tea",
    }
    r = requests.post(f"{BASE_URL}/products", json=payload, headers=headers)
    assert r.status_code == 400


def test_create_purchase_invalid():
    _, _, headers = register_and_login()
    product = {
        "name": "Tea",
        "volume": "500ml",
        "brand": "Brand1",
        "category": "Drink",
        "description": "Green tea",
    }
    r = requests.post(f"{BASE_URL}/products", json=product, headers=headers)
    product_id = r.json()["id"]
    payload = {
        "product_id": product_id,
        "quantity": 0,
        "price": 100,
        "date": "2024-01-01",
        "store": "Store",
        "receipt_id": 1,
    }
    r = requests.post(f"{BASE_URL}/purchases", json=payload, headers=headers)
    assert r.status_code == 400


def test_method_not_allowed():
    r = requests.get(f"{BASE_URL}/register")
    assert r.status_code == 405
    r = requests.get(f"{BASE_URL}/login")
    assert r.status_code == 405


def test_family_flow():
    login_a, _, headers_a = register_and_login()
    login_b, _, headers_b = register_and_login()

    r = requests.get(f"{BASE_URL}/family/members", headers=headers_a)
    assert r.status_code == 200
    assert r.json()["members"] == []

    r = requests.post(
        f"{BASE_URL}/family/invite", json={"login": login_b}, headers=headers_a
    )
    assert r.status_code == 200

    r = requests.get(f"{BASE_URL}/family/invitations", headers=headers_b)
    assert login_a in r.json()["invitations"]

    r = requests.post(
        f"{BASE_URL}/family/respond",
        json={"login": login_a, "accept": True},
        headers=headers_b,
    )
    assert r.status_code == 200

    r = requests.get(f"{BASE_URL}/family/members", headers=headers_a)
    assert set(r.json()["members"]) == {login_a, login_b}
    r = requests.get(f"{BASE_URL}/family/members", headers=headers_b)
    assert set(r.json()["members"]) == {login_a, login_b}

    product_payload = {
        "name": "Tea",
        "volume": "500ml",
        "brand": "Brand1",
        "category": "Drink",
        "description": "Green tea",
    }
    r = requests.post(f"{BASE_URL}/products", json=product_payload, headers=headers_a)
    product_id = r.json()["id"]
    purchase_payload = {
        "product_id": product_id,
        "quantity": 1,
        "price": 100,
        "date": "2024-01-01",
        "store": "Store",
        "receipt_id": 1,
    }
    r = requests.post(
        f"{BASE_URL}/purchases", json=purchase_payload, headers=headers_a
    )
    assert r.status_code == 200

    r = requests.get(f"{BASE_URL}/purchases", headers=headers_b)
    assert any(p["login"] == login_a for p in r.json()["purchases"])

    r = requests.delete(f"{BASE_URL}/family/leave", headers=headers_b)
    assert r.status_code == 200
    r = requests.get(f"{BASE_URL}/family/members", headers=headers_b)
    assert r.json()["members"] == []
    r = requests.get(f"{BASE_URL}/family/members", headers=headers_a)
    assert r.json()["members"] == []

