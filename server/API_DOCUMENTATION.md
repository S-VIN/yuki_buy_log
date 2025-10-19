# API Documentation

## Authentication

All API endpoints (except registration and login) require a Bearer token in the Authorization header:
```
Authorization: Bearer <token>
```

## Endpoints

### Authentication

#### POST /register
Register a new user.

**Request Body:**
```json
{
  "login": "username",
  "password": "password123"
}
```

**Response:**
- **200 OK**: Returns authentication token
```json
{
  "token": "jwt_token_here"
}
```
- **400 Bad Request**: Invalid request data
- **500 Internal Server Error**: Server error

#### POST /login
Authenticate an existing user.

**Request Body:**
```json
{
  "login": "username",
  "password": "password123"
}
```

**Response:**
- **200 OK**: Returns authentication token
```json
{
  "token": "jwt_token_here"
}
```
- **401 Unauthorized**: Invalid credentials
- **400 Bad Request**: Invalid request data
- **500 Internal Server Error**: Server error

### Products

#### GET /products
Get all products for the authenticated user.

**Headers:**
- `Authorization: Bearer <token>` (required)

**Response:**
- **200 OK**: Returns list of products
```json
{
  "products": [
    {
      "id": 1,
      "name": "ProductName",
      "volume": "500ml",
      "brand": "BrandName",
      "default_tags": ["tag1", "tag2"],
      "user_id": 123
    }
  ]
}
```
- **401 Unauthorized**: Invalid or missing token
- **500 Internal Server Error**: Server error

#### POST /products
Create a new product.

**Headers:**
- `Authorization: Bearer <token>` (required)

**Request Body:**
```json
{
  "name": "ProductName",
  "volume": "500ml",
  "brand": "BrandName",
  "default_tags": ["tag1", "tag2"]
}
```

**Validation Rules:**
- `name`: 1-30 characters, letters only
- `volume`: 1-10 characters
- `brand`: 1-30 characters, letters and digits only
- `default_tags`: max 10 tags, each tag 1-20 characters

**Response:**
- **200 OK**: Returns created product
```json
{
  "id": 1,
  "name": "ProductName",
  "volume": "500ml",
  "brand": "BrandName",
  "default_tags": ["tag1", "tag2"],
  "user_id": 123
}
```
- **400 Bad Request**: Validation error
- **401 Unauthorized**: Invalid or missing token
- **500 Internal Server Error**: Server error

### Purchases

#### GET /purchases
Get all purchases for the authenticated user.

**Headers:**
- `Authorization: Bearer <token>` (required)

**Response:**
- **200 OK**: Returns list of purchases
```json
{
  "purchases": [
    {
      "id": 1,
      "product_id": 1,
      "quantity": 2,
      "price": 1500,
      "date": "2023-10-15T00:00:00Z",
      "store": "StoreName",
      "tags": ["tag1", "tag2"],
      "receipt_id": 12345,
      "user_id": 123
    }
  ]
}
```
- **401 Unauthorized**: Invalid or missing token
- **500 Internal Server Error**: Server error

#### POST /purchases
Create a new purchase.

**Headers:**
- `Authorization: Bearer <token>` (required)

**Request Body:**
```json
{
  "product_id": 1,
  "quantity": 2,
  "price": 1500,
  "date": "2023-10-15T00:00:00Z",
  "store": "StoreName",
  "tags": ["tag1", "tag2"],
  "receipt_id": 12345
}
```

**Validation Rules:**
- `product_id`: positive integer
- `quantity`: 1-100000
- `price`: 1-100000000 (in kopecks/cents)
- `date`: valid date/time
- `store`: 1-30 characters, letters only
- `tags`: max 10 tags, each tag 1-20 characters
- `receipt_id`: positive integer

**Response:**
- **200 OK**: Returns created purchase
```json
{
  "id": 1,
  "product_id": 1,
  "quantity": 2,
  "price": 1500,
  "date": "2023-10-15T00:00:00Z",
  "store": "StoreName",
  "tags": ["tag1", "tag2"],
  "receipt_id": 12345,
  "user_id": 123
}
```
- **400 Bad Request**: Validation error
- **401 Unauthorized**: Invalid or missing token
- **500 Internal Server Error**: Server error

## Data Models

### User
```json
{
  "id": 123,
  "login": "username",
  "password": "password123"
}
```

### Product
```json
{
  "id": 1,
  "name": "ProductName",
  "volume": "500ml",
  "brand": "BrandName",
  "default_tags": ["tag1", "tag2"],
  "user_id": 123
}
```

### Purchase
```json
{
  "id": 1,
  "product_id": 1,
  "quantity": 2,
  "price": 1500,
  "date": "2023-10-15T00:00:00Z",
  "store": "StoreName",
  "tags": ["tag1", "tag2"],
  "receipt_id": 12345,
  "user_id": 123
}
```

## Error Responses

All endpoints may return the following error responses:

- **400 Bad Request**: Invalid request data or validation error
- **401 Unauthorized**: Missing, invalid, or expired authentication token
- **405 Method Not Allowed**: HTTP method not supported for this endpoint
- **500 Internal Server Error**: Server-side error

Error responses include a plain text error message in the response body.