# YukiBuyLog API

All requests (except registration and login) require an `Authorization` header with a bearer token.

## POST /register
Register a new user.

**Request body**
```json
{
  "login": "alice",
  "password": "secret"
}
```

**Response**
```json
{
  "token": "..."
}
```

## POST /login
Authenticate an existing user.

**Request body**
```json
{
  "login": "alice",
  "password": "secret"
}
```

**Response**
```json
{
  "token": "..."
}
```

## GET /products
Returns list of products.

**Response**
```json
{
  "products": [
    {
      "id": 1,
      "name": "Tea",
      "volume": "500ml",
      "brand": "Brand1",
      "category": "Drink",
      "description": "Green tea",
      "creation_date": "2023-01-01"
    }
  ]
}
```

## POST /products
Create a new product.

**Request body**
```json
{
  "name": "Tea",
  "volume": "500ml",
  "brand": "Brand1",
  "category": "Drink",
  "description": "Green tea"
}
```
Server assigns the creation date automatically.

**Response**
```json
{
  "id": 1,
  "name": "Tea",
  "volume": "500ml",
  "brand": "Brand1",
  "category": "Drink",
  "description": "Green tea"
}
```

## GET /purchases
Returns list of purchases.

**Response**
```json
{
  "purchases": [
    {
      "id": 1,
      "product_id": 1,
      "quantity": 2,
      "price": 100,
      "date": "2023-03-01",
      "store": "Store",
      "receipt_id": 1,
      "login": "alice"
    }
  ]
}
```

## POST /purchases
Create a new purchase.

**Request body**
```json
{
  "product_id": 1,
  "quantity": 2,
  "price": 100,
  "date": "2023-03-01",
  "store": "Store",
  "receipt_id": 1
}
```

## GET /family/members
Returns list of logins in the current user's family. If the user has no family, an empty list is returned.

**Response**
```json
{
  "members": ["alice", "bob"]
}
```

## GET /family/invitations
Returns list of pending invitations for the current user.

**Response**
```json
{
  "invitations": ["alice"]
}
```

## POST /family/invite
Invite another user to your family by login. A family is created automatically if the inviter has none.

**Request body**
```json
{
  "login": "bob"
}
```

**Response**
```json
{
  "status": "invited"
}
```

## POST /family/respond
Accept or decline a family invitation.

**Request body**
```json
{
  "login": "alice",
  "accept": true
}
```

**Response**
```json
{
  "status": "accepted"
}
```

## DELETE /family/leave
Leave the current family. If only one member remains afterwards, the family is disbanded.

**Response**
```json
{
  "status": "ok"
}
```
