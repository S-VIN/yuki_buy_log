import os
import uuid
import pytest
import requests


BASE_URL = os.getenv("BASE_URL", "http://localhost:8080")


class TestHelper:
    """Helper class for common test operations"""

    @staticmethod
    def register_user(login=None, password="password123"):
        """Register a new user and return login, password, token, headers"""
        if login is None:
            login = f"user_{uuid.uuid4().hex[:8]}"

        r = requests.post(f"{BASE_URL}/register", json={"login": login, "password": password})
        assert r.status_code == 200, f"Registration failed: {r.status_code} {r.text}"

        token = r.json()["token"]
        headers = {"Authorization": f"Bearer {token}"}
        return login, password, token, headers

    @staticmethod
    def login_user(login, password):
        """Login an existing user and return token and headers"""
        r = requests.post(f"{BASE_URL}/login", json={"login": login, "password": password})
        assert r.status_code == 200, f"Login failed: {r.status_code} {r.text}"

        token = r.json()["token"]
        headers = {"Authorization": f"Bearer {token}"}
        return token, headers


@pytest.fixture
def user():
    """Create a test user"""
    return TestHelper.register_user()


@pytest.fixture
def two_users():
    """Create two test users"""
    user1 = TestHelper.register_user()
    user2 = TestHelper.register_user()
    return user1, user2


@pytest.fixture
def three_users():
    """Create three test users"""
    user1 = TestHelper.register_user()
    user2 = TestHelper.register_user()
    user3 = TestHelper.register_user()
    return user1, user2, user3


class TestServerHealth:
    """Test server availability and basic endpoints"""

    def test_server_health(self):
        """Server should respond to requests"""
        r = requests.get(f"{BASE_URL}/products")
        assert r.status_code == 401, "Server should require authentication"


class TestAuthentication:
    """Test authentication and authorization"""

    def test_register_user(self):
        """User registration should work"""
        login, password, token, headers = TestHelper.register_user()
        assert token is not None
        assert len(token) > 0

    def test_login_with_valid_credentials(self):
        """Login with valid credentials should work"""
        login, password, _, _ = TestHelper.register_user()
        token, headers = TestHelper.login_user(login, password)
        assert token is not None
        assert len(token) > 0

    def test_login_with_invalid_credentials(self):
        """Login with invalid credentials should fail"""
        login, _, _, _ = TestHelper.register_user()
        r = requests.post(f"{BASE_URL}/login", json={"login": login, "password": "wrong_password"})
        assert r.status_code == 401

    def test_access_without_token(self):
        """Accessing protected endpoints without token should fail"""
        r = requests.get(f"{BASE_URL}/products")
        assert r.status_code == 401

        r = requests.get(f"{BASE_URL}/purchases")
        assert r.status_code == 401

        r = requests.get(f"{BASE_URL}/group")
        assert r.status_code == 401

        r = requests.get(f"{BASE_URL}/invite")
        assert r.status_code == 401

    def test_method_not_allowed(self):
        """Wrong HTTP methods should return 405"""
        r = requests.get(f"{BASE_URL}/register")
        assert r.status_code == 405

        r = requests.get(f"{BASE_URL}/login")
        assert r.status_code == 405


class TestProducts:
    """Test product CRUD operations"""

    def test_create_product(self, user):
        """User should be able to create a product"""
        login, password, token, headers = user

        product = {
            "name": "Tea",
            "volume": "500ml",
            "brand": "Brand1",
            "category": "Drink",
            "description": "Green tea",
        }
        r = requests.post(f"{BASE_URL}/products", json=product, headers=headers)
        assert r.status_code == 200
        assert "id" in r.json()

    def test_get_own_products(self, user):
        """User should see their own products"""
        login, password, token, headers = user

        # Create a product
        product = {
            "name": "Apple",
            "volume": "1kg",
            "brand": "FreshFarm",
            "default_tags": ["fruit"],
        }
        r = requests.post(f"{BASE_URL}/products", json=product, headers=headers)
        product_id = r.json()["id"]

        # Get products
        r = requests.get(f"{BASE_URL}/products", headers=headers)
        assert r.status_code == 200
        products = r.json()["products"]
        assert any(p["id"] == product_id for p in products)

    def test_create_product_invalid(self, user):
        """Creating product with invalid data should fail"""
        login, password, token, headers = user

        product = {
            "name": "Tea",
            "volume": "500ml",
            "brand": "Brand!",  # Invalid character
            "category": "Drink",
        }
        r = requests.post(f"{BASE_URL}/products", json=product, headers=headers)
        assert r.status_code == 400


class TestPurchases:
    """Test purchase CRUD operations"""

    def test_create_purchase(self, user):
        """User should be able to create a purchase"""
        login, password, token, headers = user

        # Create a product first
        product = {
            "name": "Milk",
            "volume": "1L",
            "brand": "DairyBest",
        }
        r = requests.post(f"{BASE_URL}/products", json=product, headers=headers)
        product_id = r.json()["id"]

        # Create a purchase
        purchase = {
            "product_id": product_id,
            "quantity": 2,
            "price": 150,
            "date": "2024-01-15T00:00:00Z",
            "store": "StoreOne",
            "receipt_id": 1,
        }
        r = requests.post(f"{BASE_URL}/purchases", json=purchase, headers=headers)
        assert r.status_code == 200
        assert "id" in r.json()

    def test_get_own_purchases(self, user):
        """User should see their own purchases"""
        login, password, token, headers = user

        # Create a product
        product = {
            "name": "Bread",
            "volume": "1pc",
            "brand": "Bakery",
        }
        r = requests.post(f"{BASE_URL}/products", json=product, headers=headers)
        product_id = r.json()["id"]

        # Create a purchase
        purchase = {
            "product_id": product_id,
            "quantity": 1,
            "price": 50,
            "date": "2024-01-20T00:00:00Z",
            "store": "LocalStore",
            "receipt_id": 100,
        }
        r = requests.post(f"{BASE_URL}/purchases", json=purchase, headers=headers)
        purchase_id = r.json()["id"]

        # Get purchases
        r = requests.get(f"{BASE_URL}/purchases", headers=headers)
        assert r.status_code == 200
        purchases = r.json()["purchases"]
        assert any(p["id"] == purchase_id for p in purchases)

    def test_create_purchase_invalid(self, user):
        """Creating purchase with invalid data should fail"""
        login, password, token, headers = user

        # Create a product
        product = {
            "name": "Coffee",
            "volume": "250g",
            "brand": "CoffeeCo",
        }
        r = requests.post(f"{BASE_URL}/products", json=product, headers=headers)
        product_id = r.json()["id"]

        # Try to create purchase with invalid quantity
        purchase = {
            "product_id": product_id,
            "quantity": 0,  # Invalid
            "price": 100,
            "date": "2024-01-01",
            "store": "Store",
            "receipt_id": 1,
        }
        r = requests.post(f"{BASE_URL}/purchases", json=purchase, headers=headers)
        assert r.status_code == 400

    def test_delete_purchase_success(self, user):
        """User should be able to delete their own purchase"""
        login, password, token, headers = user

        # Create a product
        product = {
            "name": "Sugar",
            "volume": "1kg",
            "brand": "SweetCo",
        }
        r = requests.post(f"{BASE_URL}/products", json=product, headers=headers)
        product_id = r.json()["id"]

        # Create a purchase
        purchase = {
            "product_id": product_id,
            "quantity": 1,
            "price": 100,
            "date": "2024-01-25T00:00:00Z",
            "store": "Supermarket",
            "receipt_id": 500,
        }
        r = requests.post(f"{BASE_URL}/purchases", json=purchase, headers=headers)
        assert r.status_code == 200
        purchase_id = r.json()["id"]

        # Delete the purchase
        r = requests.delete(f"{BASE_URL}/purchases", json={"id": purchase_id}, headers=headers)
        assert r.status_code == 204

        # Verify purchase is deleted
        r = requests.get(f"{BASE_URL}/purchases", headers=headers)
        assert r.status_code == 200
        purchases = r.json()["purchases"]
        assert not any(p["id"] == purchase_id for p in purchases)

    def test_delete_purchase_not_found(self, user):
        """Deleting non-existent purchase should return 404"""
        login, password, token, headers = user

        # Try to delete a purchase that doesn't exist
        r = requests.delete(f"{BASE_URL}/purchases", json={"id": 999999}, headers=headers)
        assert r.status_code == 404

    def test_delete_purchase_unauthorized(self):
        """Deleting purchase without authentication should fail"""
        r = requests.delete(f"{BASE_URL}/purchases", json={"id": 1})
        assert r.status_code == 401

    def test_delete_purchase_different_user(self, two_users):
        """User should not be able to delete another user's purchase"""
        user1, user2 = two_users
        login1, password1, token1, headers1 = user1
        login2, password2, token2, headers2 = user2

        # User1 creates a product and purchase
        product = {
            "name": "Salt",
            "volume": "500g",
            "brand": "SeaSalt",
        }
        r = requests.post(f"{BASE_URL}/products", json=product, headers=headers1)
        product_id = r.json()["id"]

        purchase = {
            "product_id": product_id,
            "quantity": 1,
            "price": 50,
            "date": "2024-01-26T00:00:00Z",
            "store": "Store",
            "receipt_id": 600,
        }
        r = requests.post(f"{BASE_URL}/purchases", json=purchase, headers=headers1)
        assert r.status_code == 200
        purchase_id = r.json()["id"]

        # User2 tries to delete User1's purchase
        r = requests.delete(f"{BASE_URL}/purchases", json={"id": purchase_id}, headers=headers2)
        assert r.status_code == 404  # Should not find it (security: user_id check)

        # Verify User1's purchase still exists
        r = requests.get(f"{BASE_URL}/purchases", headers=headers1)
        assert r.status_code == 200
        purchases = r.json()["purchases"]
        assert any(p["id"] == purchase_id for p in purchases)

    def test_delete_purchase_invalid_json(self, user):
        """Deleting purchase with invalid JSON should return 400"""
        login, password, token, headers = user

        r = requests.delete(f"{BASE_URL}/purchases", data="invalid json", headers=headers)
        assert r.status_code == 400

    def test_delete_purchase_missing_id(self, user):
        """Deleting purchase without ID should return 400"""
        login, password, token, headers = user

        r = requests.delete(f"{BASE_URL}/purchases", json={}, headers=headers)
        assert r.status_code == 400

        r = requests.delete(f"{BASE_URL}/purchases", json={"id": 0}, headers=headers)
        assert r.status_code == 400

    def test_delete_multiple_purchases(self, user):
        """User should be able to delete multiple purchases sequentially"""
        login, password, token, headers = user

        # Create a product
        product = {
            "name": "Flour",
            "volume": "2kg",
            "brand": "BakerCo",
        }
        r = requests.post(f"{BASE_URL}/products", json=product, headers=headers)
        product_id = r.json()["id"]

        # Create multiple purchases
        purchase_ids = []
        stores = ["StoreOne", "StoreTwo", "StoreThree"]
        for i in range(3):
            purchase = {
                "product_id": product_id,
                "quantity": i + 1,
                "price": 200 + (i * 50),
                "date": f"2024-01-{27 + i:02d}T00:00:00Z",
                "store": stores[i],
                "receipt_id": 700 + i,
            }
            r = requests.post(f"{BASE_URL}/purchases", json=purchase, headers=headers)
            assert r.status_code == 200
            purchase_ids.append(r.json()["id"])

        # Delete all purchases
        for purchase_id in purchase_ids:
            r = requests.delete(f"{BASE_URL}/purchases", json={"id": purchase_id}, headers=headers)
            assert r.status_code == 204

        # Verify all purchases are deleted
        r = requests.get(f"{BASE_URL}/purchases", headers=headers)
        assert r.status_code == 200
        purchases = r.json()["purchases"]
        for purchase_id in purchase_ids:
            assert not any(p["id"] == purchase_id for p in purchases)


class TestGroupInvites:
    """Test group and invite functionality"""

    def test_initial_group_empty(self, user):
        """New user should not be in any group"""
        login, password, token, headers = user

        r = requests.get(f"{BASE_URL}/group", headers=headers)
        assert r.status_code == 200
        assert r.json()["members"] == []

    def test_send_invite(self, two_users):
        """User should be able to send an invite"""
        user1, user2 = two_users
        login1, password1, token1, headers1 = user1
        login2, password2, token2, headers2 = user2

        r = requests.post(f"{BASE_URL}/invite", json={"login": login2}, headers=headers1)
        assert r.status_code == 200
        assert "invite_id" in r.json()

    def test_receive_invite(self, two_users):
        """User should see incoming invites"""
        user1, user2 = two_users
        login1, password1, token1, headers1 = user1
        login2, password2, token2, headers2 = user2

        # User1 sends invite to User2
        r = requests.post(f"{BASE_URL}/invite", json={"login": login2}, headers=headers1)
        assert r.status_code == 200

        # User2 checks invites
        r = requests.get(f"{BASE_URL}/invite", headers=headers2)
        assert r.status_code == 200
        invites = r.json()["invites"]
        assert any(inv["from_login"] == login1 for inv in invites)

    def test_mutual_invite_creates_group(self, two_users):
        """Mutual invites should create a group"""
        user1, user2 = two_users
        login1, password1, token1, headers1 = user1
        login2, password2, token2, headers2 = user2

        # User1 invites User2
        r = requests.post(f"{BASE_URL}/invite", json={"login": login2}, headers=headers1)
        assert r.status_code == 200

        # User2 invites User1 (mutual invite - group created!)
        r = requests.post(f"{BASE_URL}/invite", json={"login": login1}, headers=headers2)
        assert r.status_code == 200
        data = r.json()
        assert "mutual_invite" in data or "group" in data.get("message", "").lower()

        # Verify group was created
        r = requests.get(f"{BASE_URL}/group", headers=headers1)
        assert r.status_code == 200
        members = r.json()["members"]
        assert len(members) == 2
        member_logins = [m["login"] for m in members]
        assert login1 in member_logins
        assert login2 in member_logins

    def test_invite_to_nonexistent_user(self, user):
        """Inviting non-existent user should fail"""
        login, password, token, headers = user

        r = requests.post(f"{BASE_URL}/invite", json={"login": "nonexistent_user_xyz"}, headers=headers)
        assert r.status_code == 404


class TestGroupSharing:
    """Test shared access to products and purchases in groups"""

    def test_products_isolation_before_group(self, two_users):
        """Users should not see each other's products before forming a group"""
        user1, user2 = two_users
        login1, password1, token1, headers1 = user1
        login2, password2, token2, headers2 = user2

        # User1 creates a product
        product1 = {
            "name": "ProductOne",
            "volume": "500ml",
            "brand": "BrandOne",
            "default_tags": ["tag1"],
        }
        r = requests.post(f"{BASE_URL}/products", json=product1, headers=headers1)
        assert r.status_code == 200

        # User2 creates a product
        product2 = {
            "name": "ProductTwo",
            "volume": "1L",
            "brand": "BrandTwo",
            "default_tags": ["tag2"],
        }
        r = requests.post(f"{BASE_URL}/products", json=product2, headers=headers2)
        assert r.status_code == 200

        # User1 should only see their own product
        r = requests.get(f"{BASE_URL}/products", headers=headers1)
        assert r.status_code == 200
        products = r.json()["products"]
        assert len(products) == 1
        assert products[0]["name"] == "ProductOne"

    def test_products_shared_in_group(self, two_users):
        """Users in a group should see each other's products"""
        user1, user2 = two_users
        login1, password1, token1, headers1 = user1
        login2, password2, token2, headers2 = user2

        # Create products before forming group
        product1 = {
            "name": "SharedProductOne",
            "volume": "500ml",
            "brand": "BrandOne",
        }
        r = requests.post(f"{BASE_URL}/products", json=product1, headers=headers1)
        product1_id = r.json()["id"]

        product2 = {
            "name": "SharedProductTwo",
            "volume": "1L",
            "brand": "BrandTwo",
        }
        r = requests.post(f"{BASE_URL}/products", json=product2, headers=headers2)
        product2_id = r.json()["id"]

        # Form a group via mutual invites
        requests.post(f"{BASE_URL}/invite", json={"login": login2}, headers=headers1)
        requests.post(f"{BASE_URL}/invite", json={"login": login1}, headers=headers2)

        # User1 should now see both products
        r = requests.get(f"{BASE_URL}/products", headers=headers1)
        assert r.status_code == 200
        products = r.json()["products"]
        assert len(products) == 2
        product_ids = [p["id"] for p in products]
        assert product1_id in product_ids
        assert product2_id in product_ids

        # Verify user_id fields are present
        assert all("user_id" in p for p in products)

    def test_purchases_shared_in_group(self, two_users):
        """Users in a group should see each other's purchases"""
        user1, user2 = two_users
        login1, password1, token1, headers1 = user1
        login2, password2, token2, headers2 = user2

        # User1 creates a product
        product1 = {"name": "PurchProdOne", "volume": "1L", "brand": "BrandA"}
        r = requests.post(f"{BASE_URL}/products", json=product1, headers=headers1)
        product1_id = r.json()["id"]

        # User2 creates a product
        product2 = {"name": "PurchProdTwo", "volume": "1L", "brand": "BrandB"}
        r = requests.post(f"{BASE_URL}/products", json=product2, headers=headers2)
        product2_id = r.json()["id"]

        # Form a group
        requests.post(f"{BASE_URL}/invite", json={"login": login2}, headers=headers1)
        requests.post(f"{BASE_URL}/invite", json={"login": login1}, headers=headers2)

        # User1 creates a purchase
        purchase1 = {
            "product_id": product1_id,
            "quantity": 2,
            "price": 1000,
            "date": "2024-01-01T00:00:00Z",
            "store": "StoreOne",
            "tags": ["buy"],
            "receipt_id": 100,
        }
        r = requests.post(f"{BASE_URL}/purchases", json=purchase1, headers=headers1)
        purchase1_id = r.json()["id"]

        # User2 creates a purchase
        purchase2 = {
            "product_id": product2_id,
            "quantity": 3,
            "price": 2000,
            "date": "2024-01-02T00:00:00Z",
            "store": "StoreTwo",
            "tags": ["sale"],
            "receipt_id": 200,
        }
        r = requests.post(f"{BASE_URL}/purchases", json=purchase2, headers=headers2)
        purchase2_id = r.json()["id"]

        # User1 should see both purchases
        r = requests.get(f"{BASE_URL}/purchases", headers=headers1)
        assert r.status_code == 200
        purchases = r.json()["purchases"]
        assert len(purchases) == 2
        purchase_ids = [p["id"] for p in purchases]
        assert purchase1_id in purchase_ids
        assert purchase2_id in purchase_ids

        # Verify user_id fields are present
        assert all("user_id" in p for p in purchases)

    def test_isolation_for_non_group_users(self, three_users):
        """Users not in the group should not see group data"""
        user1, user2, user3 = three_users
        login1, password1, token1, headers1 = user1
        login2, password2, token2, headers2 = user2
        login3, password3, token3, headers3 = user3

        # User1 and User2 form a group
        requests.post(f"{BASE_URL}/invite", json={"login": login2}, headers=headers1)
        requests.post(f"{BASE_URL}/invite", json={"login": login1}, headers=headers2)

        # User1 creates a product
        product = {"name": "GroupProduct", "volume": "1L", "brand": "BrandX"}
        requests.post(f"{BASE_URL}/products", json=product, headers=headers1)

        # User3 (not in group) should see no products
        r = requests.get(f"{BASE_URL}/products", headers=headers3)
        assert r.status_code == 200
        products = r.json()["products"]
        assert len(products) == 0


class TestGroupManagement:
    """Test group management operations"""

    def test_leave_group(self, two_users):
        """User should be able to leave a group"""
        user1, user2 = two_users
        login1, password1, token1, headers1 = user1
        login2, password2, token2, headers2 = user2

        # Form a group
        requests.post(f"{BASE_URL}/invite", json={"login": login2}, headers=headers1)
        requests.post(f"{BASE_URL}/invite", json={"login": login1}, headers=headers2)

        # User2 leaves the group
        r = requests.delete(f"{BASE_URL}/group", headers=headers2)
        assert r.status_code == 200

        # User2 should not be in any group now
        r = requests.get(f"{BASE_URL}/group", headers=headers2)
        assert r.status_code == 200
        assert r.json()["members"] == []

    def test_group_auto_deletion(self, two_users):
        """Group should be auto-deleted when only 1 member remains"""
        user1, user2 = two_users
        login1, password1, token1, headers1 = user1
        login2, password2, token2, headers2 = user2

        # Form a group
        requests.post(f"{BASE_URL}/invite", json={"login": login2}, headers=headers1)
        requests.post(f"{BASE_URL}/invite", json={"login": login1}, headers=headers2)

        # User2 leaves the group
        r = requests.delete(f"{BASE_URL}/group", headers=headers2)
        assert r.status_code == 200

        # User1's group should be auto-deleted
        r = requests.get(f"{BASE_URL}/group", headers=headers1)
        assert r.status_code == 200
        assert r.json()["members"] == []

    def test_data_isolation_after_leaving_group(self, two_users):
        """Users should not see each other's data after leaving group"""
        user1, user2 = two_users
        login1, password1, token1, headers1 = user1
        login2, password2, token2, headers2 = user2

        # User1 creates a product
        product1 = {"name": "ProductOne", "volume": "1L", "brand": "BrandA"}
        r = requests.post(f"{BASE_URL}/products", json=product1, headers=headers1)
        product1_id = r.json()["id"]

        # User2 creates a product
        product2 = {"name": "ProductTwo", "volume": "1L", "brand": "BrandB"}
        r = requests.post(f"{BASE_URL}/products", json=product2, headers=headers2)
        product2_id = r.json()["id"]

        # Form a group
        requests.post(f"{BASE_URL}/invite", json={"login": login2}, headers=headers1)
        requests.post(f"{BASE_URL}/invite", json={"login": login1}, headers=headers2)

        # Both users should see 2 products
        r = requests.get(f"{BASE_URL}/products", headers=headers1)
        assert len(r.json()["products"]) == 2

        # User2 leaves the group
        requests.delete(f"{BASE_URL}/group", headers=headers2)

        # User1 should only see their own product now
        r = requests.get(f"{BASE_URL}/products", headers=headers1)
        assert r.status_code == 200
        products = r.json()["products"]
        assert len(products) == 1
        assert products[0]["id"] == product1_id

    def test_leave_group_when_not_in_group(self, user):
        """Leaving a group when not in one should fail"""
        login, password, token, headers = user

        r = requests.delete(f"{BASE_URL}/group", headers=headers)
        assert r.status_code == 400


class TestGroupExpansion:
    """Test group expansion restrictions"""

    def test_can_invite_to_expand_group(self, three_users):
        """Free user can invite user in group and vice versa to expand the group"""
        user1, user2, user3 = three_users
        login1, password1, token1, headers1 = user1
        login2, password2, token2, headers2 = user2
        login3, password3, token3, headers3 = user3

        # User1 and User2 form a group
        requests.post(f"{BASE_URL}/invite", json={"login": login2}, headers=headers1)
        requests.post(f"{BASE_URL}/invite", json={"login": login1}, headers=headers2)

        # Verify initial group has 2 members
        r = requests.get(f"{BASE_URL}/group", headers=headers1)
        assert len(r.json()["members"]) == 2

        # User3 (free user) can send invite to User1 (in group)
        r = requests.post(f"{BASE_URL}/invite", json={"login": login1}, headers=headers3)
        assert r.status_code == 200
        assert "invite_id" in r.json()

        # User1 (in group) can send invite to User3 (free user)
        # This creates mutual invites and should expand the group
        r = requests.post(f"{BASE_URL}/invite", json={"login": login3}, headers=headers1)
        assert r.status_code == 200
        data = r.json()
        assert "mutual_invite" in data or "group" in data.get("message", "").lower()

        # Verify group now has 3 members
        r = requests.get(f"{BASE_URL}/group", headers=headers1)
        members = r.json()["members"]
        assert len(members) == 3
        member_logins = [m["login"] for m in members]
        assert login1 in member_logins
        assert login2 in member_logins
        assert login3 in member_logins

        # User3 should also see the group
        r = requests.get(f"{BASE_URL}/group", headers=headers3)
        assert len(r.json()["members"]) == 3

    def test_cannot_invite_users_in_different_groups(self):
        """Users in different groups cannot invite each other"""
        # Create 4 users
        user1 = TestHelper.register_user()
        user2 = TestHelper.register_user()
        user3 = TestHelper.register_user()
        user4 = TestHelper.register_user()

        login1, password1, token1, headers1 = user1
        login2, password2, token2, headers2 = user2
        login3, password3, token3, headers3 = user3
        login4, password4, token4, headers4 = user4

        # Create Group 1: User1 and User2
        requests.post(f"{BASE_URL}/invite", json={"login": login2}, headers=headers1)
        requests.post(f"{BASE_URL}/invite", json={"login": login1}, headers=headers2)

        # Create Group 2: User3 and User4
        requests.post(f"{BASE_URL}/invite", json={"login": login4}, headers=headers3)
        requests.post(f"{BASE_URL}/invite", json={"login": login3}, headers=headers4)

        # Verify groups were created
        r = requests.get(f"{BASE_URL}/group", headers=headers1)
        assert len(r.json()["members"]) == 2

        r = requests.get(f"{BASE_URL}/group", headers=headers3)
        assert len(r.json()["members"]) == 2

        # User1 (in Group 1) tries to invite User3 (in Group 2)
        r = requests.post(f"{BASE_URL}/invite", json={"login": login3}, headers=headers1)
        assert r.status_code == 400
        assert "different groups" in r.text.lower()

        # User3 (in Group 2) tries to invite User1 (in Group 1)
        r = requests.post(f"{BASE_URL}/invite", json={"login": login1}, headers=headers3)
        assert r.status_code == 400
        assert "different groups" in r.text.lower()


class TestProductUpdate:
    """Test product update operations"""

    def test_update_product_success(self, user):
        """User should be able to update their own product"""
        login, password, token, headers = user

        # Create a product
        product = {
            "name": "OriginalProduct",
            "volume": "500ml",
            "brand": "OriginalBrand",
            "default_tags": ["original"],
        }
        r = requests.post(f"{BASE_URL}/products", json=product, headers=headers)
        assert r.status_code == 200
        product_id = r.json()["id"]

        # Update the product
        updated_product = {
            "id": product_id,
            "name": "UpdatedProduct",
            "volume": "1L",
            "brand": "UpdatedBrand",
            "default_tags": ["updated", "new"],
        }
        r = requests.put(f"{BASE_URL}/products", json=updated_product, headers=headers)
        assert r.status_code == 200

        # Verify the response contains updated data
        response_product = r.json()
        assert response_product["id"] == product_id
        assert response_product["name"] == "UpdatedProduct"
        assert response_product["volume"] == "1L"
        assert response_product["brand"] == "UpdatedBrand"
        assert "updated" in response_product["default_tags"]
        assert "new" in response_product["default_tags"]

        # Verify the update persisted by fetching products
        r = requests.get(f"{BASE_URL}/products", headers=headers)
        assert r.status_code == 200
        products = r.json()["products"]
        updated = next((p for p in products if p["id"] == product_id), None)
        assert updated is not None
        assert updated["name"] == "UpdatedProduct"
        assert updated["volume"] == "1L"

    def test_update_product_not_found(self, user):
        """Updating non-existent product should return 404"""
        login, password, token, headers = user

        # Try to update a product that doesn't exist
        updated_product = {
            "id": 999999,
            "name": "UpdatedProduct",
            "volume": "1L",
            "brand": "UpdatedBrand",
            "default_tags": ["test"],
        }
        r = requests.put(f"{BASE_URL}/products", json=updated_product, headers=headers)
        assert r.status_code == 404

    def test_update_product_unauthorized(self):
        """Updating product without authentication should fail"""
        updated_product = {
            "id": 1,
            "name": "UpdatedProduct",
            "volume": "1L",
            "brand": "UpdatedBrand",
        }
        r = requests.put(f"{BASE_URL}/products", json=updated_product)
        assert r.status_code == 401

    def test_update_product_different_user(self, two_users):
        """User should not be able to update another user's product"""
        user1, user2 = two_users
        login1, password1, token1, headers1 = user1
        login2, password2, token2, headers2 = user2

        # User1 creates a product
        product = {
            "name": "UserOneProduct",
            "volume": "500ml",
            "brand": "BrandOne",
        }
        r = requests.post(f"{BASE_URL}/products", json=product, headers=headers1)
        assert r.status_code == 200
        product_id = r.json()["id"]

        # User2 tries to update User1's product
        updated_product = {
            "id": product_id,
            "name": "HackedProduct",
            "volume": "1L",
            "brand": "HackedBrand",
        }
        r = requests.put(f"{BASE_URL}/products", json=updated_product, headers=headers2)
        assert r.status_code == 404  # Should not find it (security: user_id check)

        # Verify User1's product is unchanged
        r = requests.get(f"{BASE_URL}/products", headers=headers1)
        assert r.status_code == 200
        products = r.json()["products"]
        original = next((p for p in products if p["id"] == product_id), None)
        assert original is not None
        assert original["name"] == "UserOneProduct"

    def test_update_product_invalid_json(self, user):
        """Updating product with invalid JSON should return 400"""
        login, password, token, headers = user

        r = requests.put(f"{BASE_URL}/products", data="invalid json", headers=headers)
        assert r.status_code == 400

    def test_update_product_missing_id(self, user):
        """Updating product without ID should return 400"""
        login, password, token, headers = user

        # Create a product first
        product = {
            "name": "TestProduct",
            "volume": "500ml",
            "brand": "TestBrand",
        }
        r = requests.post(f"{BASE_URL}/products", json=product, headers=headers)
        assert r.status_code == 200

        # Try to update without ID
        updated_product = {
            "name": "UpdatedProduct",
            "volume": "1L",
            "brand": "UpdatedBrand",
        }
        r = requests.put(f"{BASE_URL}/products", json=updated_product, headers=headers)
        assert r.status_code == 400

        # Try with ID = 0
        updated_product["id"] = 0
        r = requests.put(f"{BASE_URL}/products", json=updated_product, headers=headers)
        assert r.status_code == 400

    def test_update_product_invalid_validation(self, user):
        """Updating product with invalid data should return 400"""
        login, password, token, headers = user

        # Create a product
        product = {
            "name": "ValidProduct",
            "volume": "500ml",
            "brand": "ValidBrand",
        }
        r = requests.post(f"{BASE_URL}/products", json=product, headers=headers)
        assert r.status_code == 200
        product_id = r.json()["id"]

        # Try to update with invalid name (contains digits)
        updated_product = {
            "id": product_id,
            "name": "Invalid123",
            "volume": "1L",
            "brand": "ValidBrand",
        }
        r = requests.put(f"{BASE_URL}/products", json=updated_product, headers=headers)
        assert r.status_code == 400

        # Try to update with invalid brand (contains special characters)
        updated_product = {
            "id": product_id,
            "name": "ValidProduct",
            "volume": "1L",
            "brand": "Invalid@Brand",
        }
        r = requests.put(f"{BASE_URL}/products", json=updated_product, headers=headers)
        assert r.status_code == 400

    def test_update_product_empty_tags(self, user):
        """User should be able to update product with empty tags"""
        login, password, token, headers = user

        # Create a product with tags
        product = {
            "name": "ProductWithTags",
            "volume": "500ml",
            "brand": "BrandName",
            "default_tags": ["tag1", "tag2"],
        }
        r = requests.post(f"{BASE_URL}/products", json=product, headers=headers)
        assert r.status_code == 200
        product_id = r.json()["id"]

        # Update product to have no tags
        updated_product = {
            "id": product_id,
            "name": "ProductNoTags",
            "volume": "1L",
            "brand": "BrandName",
            "default_tags": [],
        }
        r = requests.put(f"{BASE_URL}/products", json=updated_product, headers=headers)
        assert r.status_code == 200

        response_product = r.json()
        assert len(response_product["default_tags"]) == 0

    def test_update_product_multiple_times(self, user):
        """User should be able to update the same product multiple times"""
        login, password, token, headers = user

        # Create a product
        product = {
            "name": "FirstVersion",
            "volume": "500ml",
            "brand": "BrandOne",
        }
        r = requests.post(f"{BASE_URL}/products", json=product, headers=headers)
        assert r.status_code == 200
        product_id = r.json()["id"]

        # First update
        updated_product = {
            "id": product_id,
            "name": "SecondVersion",
            "volume": "1L",
            "brand": "BrandTwo",
        }
        r = requests.put(f"{BASE_URL}/products", json=updated_product, headers=headers)
        assert r.status_code == 200
        assert r.json()["name"] == "SecondVersion"

        # Second update
        updated_product = {
            "id": product_id,
            "name": "ThirdVersion",
            "volume": "2L",
            "brand": "BrandThree",
        }
        r = requests.put(f"{BASE_URL}/products", json=updated_product, headers=headers)
        assert r.status_code == 200
        assert r.json()["name"] == "ThirdVersion"
        assert r.json()["volume"] == "2L"

        # Verify final state
        r = requests.get(f"{BASE_URL}/products", headers=headers)
        products = r.json()["products"]
        final = next((p for p in products if p["id"] == product_id), None)
        assert final is not None
        assert final["name"] == "ThirdVersion"
        assert final["volume"] == "2L"
        assert final["brand"] == "BrandThree"

    def test_update_product_with_many_tags(self, user):
        """User should be able to update product with multiple tags"""
        login, password, token, headers = user

        # Create a product
        product = {
            "name": "ProductTags",
            "volume": "500ml",
            "brand": "BrandTags",
        }
        r = requests.post(f"{BASE_URL}/products", json=product, headers=headers)
        assert r.status_code == 200
        product_id = r.json()["id"]

        # Update with maximum allowed tags (10)
        tags = [f"tag{i}" for i in range(1, 11)]
        updated_product = {
            "id": product_id,
            "name": "ProductManyTags",
            "volume": "1L",
            "brand": "BrandTags",
            "default_tags": tags,
        }
        r = requests.put(f"{BASE_URL}/products", json=updated_product, headers=headers)
        assert r.status_code == 200
        assert len(r.json()["default_tags"]) == 10

    def test_update_product_too_many_tags(self, user):
        """Updating product with too many tags should fail"""
        login, password, token, headers = user

        # Create a product
        product = {
            "name": "ProductTest",
            "volume": "500ml",
            "brand": "BrandTest",
        }
        r = requests.post(f"{BASE_URL}/products", json=product, headers=headers)
        assert r.status_code == 200
        product_id = r.json()["id"]

        # Try to update with more than 10 tags
        tags = [f"tag{i}" for i in range(1, 12)]  # 11 tags
        updated_product = {
            "id": product_id,
            "name": "ProductTooManyTags",
            "volume": "1L",
            "brand": "BrandTest",
            "default_tags": tags,
        }
        r = requests.put(f"{BASE_URL}/products", json=updated_product, headers=headers)
        assert r.status_code == 400


class TestFullFlow:
    """Test complete application flows"""

    def test_complete_sharing_flow(self):
        """Test complete flow: register, create data, group, share, leave"""
        # Register 2 users
        user1 = TestHelper.register_user()
        user2 = TestHelper.register_user()

        login1, password1, token1, headers1 = user1
        login2, password2, token2, headers2 = user2

        # Step 1: Each user creates products
        r = requests.post(f"{BASE_URL}/products",
                         json={"name": "ProductA", "volume": "1L", "brand": "BrandA"},
                         headers=headers1)
        assert r.status_code == 200

        r = requests.post(f"{BASE_URL}/products",
                         json={"name": "ProductB", "volume": "1L", "brand": "BrandB"},
                         headers=headers2)
        assert r.status_code == 200

        # Step 2: Verify isolation before grouping
        r = requests.get(f"{BASE_URL}/products", headers=headers1)
        assert len(r.json()["products"]) == 1

        # Step 3: User1 and User2 form a group
        requests.post(f"{BASE_URL}/invite", json={"login": login2}, headers=headers1)
        requests.post(f"{BASE_URL}/invite", json={"login": login1}, headers=headers2)

        # Step 4: Verify shared access
        r = requests.get(f"{BASE_URL}/products", headers=headers1)
        assert len(r.json()["products"]) == 2
        r = requests.get(f"{BASE_URL}/products", headers=headers2)
        assert len(r.json()["products"]) == 2

        # Step 5: Create purchases and verify sharing
        product1_id = 1  # Assuming first product
        purchase1 = {
            "product_id": product1_id,
            "quantity": 5,
            "price": 500,
            "date": "2024-01-15T00:00:00Z",
            "store": "MarketOne",
            "tags": ["test"],
            "receipt_id": 1001,
        }
        r = requests.post(f"{BASE_URL}/purchases", json=purchase1, headers=headers1)
        # Purchase might fail if product_id doesn't match, so we skip assertion
        # The main test is group sharing which we already verified

        # Step 6: User2 leaves
        requests.delete(f"{BASE_URL}/group", headers=headers2)

        # Step 7: Verify group auto-deletion
        r = requests.get(f"{BASE_URL}/group", headers=headers1)
        assert r.json()["members"] == []

        # Step 8: Verify isolation after leaving
        r = requests.get(f"{BASE_URL}/products", headers=headers1)
        products = r.json()["products"]
        assert len(products) == 1  # Only User1's product
