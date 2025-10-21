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

### Groups

Groups allow users to share access to purchases and products. Users in a group can view each other's purchases and products.

#### Group Mechanics

**Key Rules:**
- Maximum group size: **5 members**
- Each user can only be in **one group** at a time
- Groups are created automatically when two users send mutual invites to each other
- When a group has only 1 member remaining (after others leave), the group is **automatically deleted**

**Group Creation Flow:**
1. User A sends an invite to User B
2. User B sends an invite to User A (mutual invite)
3. System detects mutual invites and automatically creates a group with both users
4. Both invites are deleted from the system

**Group Expansion:**
- Users already in a group can invite **free users** (users not in any group)
- Cannot invite users who are already in another group
- Group size cannot exceed 5 members

#### GET /group
Get all members of the authenticated user's group.

**Headers:**
- `Authorization: Bearer <token>` (required)

**Response:**
- **200 OK**: Returns list of group members
```json
{
  "members": [
    {
      "group_id": 1,
      "user_id": 123,
      "login": "user1"
    },
    {
      "group_id": 1,
      "user_id": 456,
      "login": "user2"
    }
  ]
}
```
- **200 OK**: Empty list if user is not in a group
```json
{
  "members": []
}
```
- **401 Unauthorized**: Invalid or missing token
- **500 Internal Server Error**: Server error

#### DELETE /group
Leave the current group. If only 1 member remains after leaving, the group is automatically deleted.

**Headers:**
- `Authorization: Bearer <token>` (required)

**Response:**
- **200 OK**: Successfully left the group
```json
{
  "message": "left group successfully"
}
```
- **400 Bad Request**: User is not in a group
- **401 Unauthorized**: Invalid or missing token
- **500 Internal Server Error**: Server error

### Invites

Invites allow users to form groups by sending and accepting invitations.

#### Invite Mechanics

**Key Rules:**
- Cannot invite yourself
- Cannot invite a user who is already in a group (unless you're inviting them to join your group)
- Free users (not in a group) can only send invites to other free users
- Users in a group can send invites to free users to expand their group
- Duplicate invites are prevented by database constraints
- **Mutual invites** automatically create a group and delete both invites

**Invitation Flow Example 1 (New Group):**
1. Free User A sends invite to Free User B
2. Free User B sends invite to Free User A
3. System detects mutual invites
4. Creates new group with both users
5. Deletes both invites

**Invitation Flow Example 2 (Expanding Existing Group):**
1. User A (in group with 3 members) sends invite to Free User B
2. Free User B sends invite to User A
3. System detects mutual invites
4. Adds User B to User A's existing group (now 4 members)
5. Deletes both invites

#### GET /invite
Get all incoming invites for the authenticated user.

**Headers:**
- `Authorization: Bearer <token>` (required)

**Response:**
- **200 OK**: Returns list of incoming invites
```json
{
  "invites": [
    {
      "id": 1,
      "from_user_id": 456,
      "to_user_id": 123,
      "from_login": "sender_username",
      "to_login": "your_username",
      "created_at": "2023-10-15T12:34:56Z"
    }
  ]
}
```
- **401 Unauthorized**: Invalid or missing token
- **500 Internal Server Error**: Server error

#### POST /invite
Send an invite to another user. If mutual invites are detected, a group is created automatically.

**Headers:**
- `Authorization: Bearer <token>` (required)

**Request Body:**
```json
{
  "login": "target_username"
}
```

**Response:**
- **200 OK**: Invite sent successfully
```json
{
  "message": "invite sent",
  "invite_id": 1
}
```
- **200 OK**: Mutual invite detected, group created
```json
{
  "message": "group created",
  "mutual_invite": true
}
```
- **400 Bad Request**: Various validation errors
  - Target user is already in a group
  - Current user's group has reached maximum size (5 members)
  - Invite already exists
- **404 Not Found**: Target user not found
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

### GroupMember
```json
{
  "group_id": 1,
  "user_id": 123,
  "login": "username"
}
```

### Invite
```json
{
  "id": 1,
  "from_user_id": 456,
  "to_user_id": 123,
  "from_login": "sender_username",
  "to_login": "receiver_username",
  "created_at": "2023-10-15T12:34:56Z"
}
```

## Error Responses

All endpoints may return the following error responses:

- **400 Bad Request**: Invalid request data or validation error
- **401 Unauthorized**: Missing, invalid, or expired authentication token
- **405 Method Not Allowed**: HTTP method not supported for this endpoint
- **500 Internal Server Error**: Server-side error

Error responses include a plain text error message in the response body.