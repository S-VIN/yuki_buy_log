# Manual Testing Guide for Groups and Invites

This guide provides step-by-step instructions for manually testing the new group and invite functionality.

## Prerequisites

1. Start the server and database:
```bash
docker-compose up --build
```

2. The server should be running at `http://localhost:8080`

## Test Scenario 1: Create Users and Form a Group

### Step 1: Register three users

```bash
# Register user1
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"login": "user1", "password": "password123"}'
# Save the token from response as USER1_TOKEN

# Register user2
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"login": "user2", "password": "password123"}'
# Save the token from response as USER2_TOKEN

# Register user3
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"login": "user3", "password": "password123"}'
# Save the token from response as USER3_TOKEN
```

### Step 2: Test initial group status (should be empty)

```bash
# User1 checks their group (should be empty)
curl -X GET http://localhost:8080/group \
  -H "Authorization: Bearer $USER1_TOKEN"
# Expected: {"members":[]}
```

### Step 3: Send invites (mutual invite scenario)

```bash
# User1 sends invite to User2
curl -X POST http://localhost:8080/invite \
  -H "Authorization: Bearer $USER1_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"login": "user2"}'
# Expected: {"message":"invite sent","invite_id":1}

# User2 checks incoming invites
curl -X GET http://localhost:8080/invite \
  -H "Authorization: Bearer $USER2_TOKEN"
# Expected: Should show invite from user1

# User2 sends invite to User1 (mutual invite - group should be created!)
curl -X POST http://localhost:8080/invite \
  -H "Authorization: Bearer $USER2_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"login": "user1"}'
# Expected: {"message":"group created","mutual_invite":true}
```

### Step 4: Verify group was created

```bash
# User1 checks their group
curl -X GET http://localhost:8080/group \
  -H "Authorization: Bearer $USER1_TOKEN"
# Expected: {"members":[{"group_id":1,"user_id":...,"login":"user1"},{"group_id":1,"user_id":...,"login":"user2"}]}

# User2 checks their group (should be same)
curl -X GET http://localhost:8080/group \
  -H "Authorization: Bearer $USER2_TOKEN"
# Expected: Same as above
```

### Step 5: Verify invites were deleted

```bash
# User1 checks invites (should be empty)
curl -X GET http://localhost:8080/invite \
  -H "Authorization: Bearer $USER1_TOKEN"
# Expected: {"invites":[]}

# User2 checks invites (should be empty)
curl -X GET http://localhost:8080/invite \
  -H "Authorization: Bearer $USER2_TOKEN"
# Expected: {"invites":[]}
```

## Test Scenario 2: Expand Group

### Step 6: Add User3 to the group

```bash
# User1 sends invite to User3
curl -X POST http://localhost:8080/invite \
  -H "Authorization: Bearer $USER1_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"login": "user3"}'
# Expected: {"message":"invite sent","invite_id":2}

# User3 sends mutual invite to User1
curl -X POST http://localhost:8080/invite \
  -H "Authorization: Bearer $USER3_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"login": "user1"}'
# Expected: {"message":"group created","mutual_invite":true}

# Verify User3 is now in the group
curl -X GET http://localhost:8080/group \
  -H "Authorization: Bearer $USER1_TOKEN"
# Expected: Should show 3 members (user1, user2, user3)
```

## Test Scenario 3: Test Group Size Limit

### Step 7: Create 5-member group and test limit

```bash
# Register user4 and user5
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"login": "user4", "password": "password123"}'
# Save token as USER4_TOKEN

curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"login": "user5", "password": "password123"}'
# Save token as USER5_TOKEN

# Add user4 to group (via mutual invites with user1)
curl -X POST http://localhost:8080/invite \
  -H "Authorization: Bearer $USER1_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"login": "user4"}'

curl -X POST http://localhost:8080/invite \
  -H "Authorization: Bearer $USER4_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"login": "user1"}'

# Add user5 to group (via mutual invites with user1)
curl -X POST http://localhost:8080/invite \
  -H "Authorization: Bearer $USER1_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"login": "user5"}'

curl -X POST http://localhost:8080/invite \
  -H "Authorization: Bearer $USER5_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"login": "user1"}'

# Verify group now has 5 members
curl -X GET http://localhost:8080/group \
  -H "Authorization: Bearer $USER1_TOKEN"
# Expected: 5 members

# Register user6 and try to add (should fail - group full)
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"login": "user6", "password": "password123"}'
# Save token as USER6_TOKEN

curl -X POST http://localhost:8080/invite \
  -H "Authorization: Bearer $USER1_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"login": "user6"}'
# Expected: {"error":"group has reached maximum size of 5 members"} or similar
```

## Test Scenario 4: Test Invite Restrictions

### Step 8: Test that users in a group can't invite users already in another group

```bash
# Register user7 and user8
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"login": "user7", "password": "password123"}'
# Save token as USER7_TOKEN

curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"login": "user8", "password": "password123"}'
# Save token as USER8_TOKEN

# User7 and User8 form their own group
curl -X POST http://localhost:8080/invite \
  -H "Authorization: Bearer $USER7_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"login": "user8"}'

curl -X POST http://localhost:8080/invite \
  -H "Authorization: Bearer $USER8_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"login": "user7"}'

# User6 (free user) tries to invite User7 (already in group)
curl -X POST http://localhost:8080/invite \
  -H "Authorization: Bearer $USER6_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"login": "user7"}'
# Expected: {"error":"target user is already in a group"} or similar
```

## Test Scenario 5: Leave Group

### Step 9: Test leaving a group

```bash
# User3 leaves the group
curl -X DELETE http://localhost:8080/group \
  -H "Authorization: Bearer $USER3_TOKEN"
# Expected: {"message":"left group successfully"}

# Verify User3 is no longer in the group
curl -X GET http://localhost:8080/group \
  -H "Authorization: Bearer $USER3_TOKEN"
# Expected: {"members":[]}

# Verify group still has 4 members
curl -X GET http://localhost:8080/group \
  -H "Authorization: Bearer $USER1_TOKEN"
# Expected: 4 members (user1, user2, user4, user5)

# User3 (now free) can be invited again
curl -X POST http://localhost:8080/invite \
  -H "Authorization: Bearer $USER1_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"login": "user3"}'
# Expected: {"message":"invite sent","invite_id":...}
```

## Test Scenario 6: Auto-delete Group

### Step 10: Test automatic group deletion when only 1 member remains

```bash
# User7 and User8 are in a 2-member group
# User8 leaves the group
curl -X DELETE http://localhost:8080/group \
  -H "Authorization: Bearer $USER8_TOKEN"
# Expected: {"message":"left group successfully"}

# Verify User7's group was automatically deleted
curl -X GET http://localhost:8080/group \
  -H "Authorization: Bearer $USER7_TOKEN"
# Expected: {"members":[]} (group should be auto-deleted)
```

## Test Scenario 7: Error Cases

### Step 11: Test various error cases

```bash
# Try to leave a group when not in one
curl -X DELETE http://localhost:8080/group \
  -H "Authorization: Bearer $USER6_TOKEN"
# Expected: 400 Bad Request - "you are not in a group"

# Try to send invite to non-existent user
curl -X POST http://localhost:8080/invite \
  -H "Authorization: Bearer $USER6_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"login": "nonexistent"}'
# Expected: 404 Not Found - "user not found"

# Try to access endpoints without authentication
curl -X GET http://localhost:8080/group
# Expected: 401 Unauthorized

curl -X GET http://localhost:8080/invite
# Expected: 401 Unauthorized
```

## Summary of Expected Behaviors

✅ **Mutual invites** automatically create a group and delete both invites
✅ **Groups** can have a maximum of 5 members
✅ **Users in groups** can invite free users (not in any group)
✅ **Free users** can only invite other free users
✅ **Cannot invite** users who are already in a group
✅ **Leaving a group** removes you from the group
✅ **Auto-deletion** occurs when a group has only 1 member remaining
✅ **All endpoints** require authentication

## Notes

- Replace `$USER1_TOKEN`, `$USER2_TOKEN`, etc. with actual tokens from registration responses
- Use `-v` flag with curl for verbose output to see response headers and status codes
- Use `| jq` at the end of curl commands for formatted JSON output (if jq is installed)
