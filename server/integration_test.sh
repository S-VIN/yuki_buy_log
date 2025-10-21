#!/bin/bash

# Integration Test Script for Yuki Buy Log Server
# This script tests the complete flow of group sharing functionality

set -e  # Exit on error

BASE_URL="${BASE_URL:-http://localhost:8080}"
VERBOSE="${VERBOSE:-1}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

log_info() {
    if [ "$VERBOSE" = "1" ]; then
        echo -e "${GREEN}[INFO]${NC} $1"
    fi
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

log_test() {
    echo -e "${YELLOW}[TEST]${NC} $1"
}

# Test counter
TESTS_RUN=0
TESTS_PASSED=0
TESTS_FAILED=0

assert_status() {
    local expected=$1
    local actual=$2
    local test_name=$3

    TESTS_RUN=$((TESTS_RUN + 1))

    if [ "$expected" -eq "$actual" ]; then
        log_info "✓ $test_name (HTTP $actual)"
        TESTS_PASSED=$((TESTS_PASSED + 1))
        return 0
    else
        log_error "✗ $test_name (Expected HTTP $expected, got $actual)"
        TESTS_FAILED=$((TESTS_FAILED + 1))
        return 1
    fi
}

# Clean up function
cleanup() {
    log_info "Cleaning up..."
}

trap cleanup EXIT

echo "========================================="
echo "  Yuki Buy Log Integration Test Suite"
echo "  Testing Server at: $BASE_URL"
echo "========================================="
echo ""

# Test 1: Server Health Check
log_test "Test 1: Server Health Check"
STATUS=$(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/products" || echo "000")
if [ "$STATUS" = "401" ]; then
    log_info "✓ Server is responding (unauthorized as expected)"
    TESTS_RUN=$((TESTS_RUN + 1))
    TESTS_PASSED=$((TESTS_PASSED + 1))
else
    log_error "✗ Server health check failed (got HTTP $STATUS)"
    TESTS_RUN=$((TESTS_RUN + 1))
    TESTS_FAILED=$((TESTS_FAILED + 1))
    exit 1
fi

# Test 2: Register User 1
log_test "Test 2: Register User 1"
RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/register" \
    -H "Content-Type: application/json" \
    -d '{"login":"testuser1","password":"password123"}')
STATUS=$(echo "$RESPONSE" | tail -n1)
BODY=$(echo "$RESPONSE" | head -n-1)
TOKEN1=$(echo "$BODY" | grep -o '"token":"[^"]*"' | cut -d'"' -f4 || echo "")

if assert_status 200 "$STATUS" "Register User 1"; then
    log_info "Token1: ${TOKEN1:0:20}..."
fi

# Test 3: Register User 2
log_test "Test 3: Register User 2"
RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/register" \
    -H "Content-Type: application/json" \
    -d '{"login":"testuser2","password":"password123"}')
STATUS=$(echo "$RESPONSE" | tail -n1)
BODY=$(echo "$RESPONSE" | head -n-1)
TOKEN2=$(echo "$BODY" | grep -o '"token":"[^"]*"' | cut -d'"' -f4 || echo "")

if assert_status 200 "$STATUS" "Register User 2"; then
    log_info "Token2: ${TOKEN2:0:20}..."
fi

# Test 4: Register User 3
log_test "Test 4: Register User 3"
RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/register" \
    -H "Content-Type: application/json" \
    -d '{"login":"testuser3","password":"password123"}')
STATUS=$(echo "$RESPONSE" | tail -n1)
BODY=$(echo "$RESPONSE" | head -n-1)
TOKEN3=$(echo "$BODY" | grep -o '"token":"[^"]*"' | cut -d'"' -f4 || echo "")

if assert_status 200 "$STATUS" "Register User 3"; then
    log_info "Token3: ${TOKEN3:0:20}..."
fi

# Test 5: User1 creates a product
log_test "Test 5: User1 creates a product"
STATUS=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE_URL/products" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN1" \
    -d '{"name":"Product1","volume":"500ml","brand":"Brand1","default_tags":["tag1"]}')
assert_status 200 "$STATUS" "User1 creates product"

# Test 6: User2 creates a product
log_test "Test 6: User2 creates a product"
STATUS=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE_URL/products" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN2" \
    -d '{"name":"Product2","volume":"1L","brand":"Brand2","default_tags":["tag2"]}')
assert_status 200 "$STATUS" "User2 creates product"

# Test 7: User1 gets products (should only see their own)
log_test "Test 7: User1 gets products (no group yet)"
RESPONSE=$(curl -s -X GET "$BASE_URL/products" \
    -H "Authorization: Bearer $TOKEN1")
COUNT=$(echo "$RESPONSE" | grep -o '"id"' | wc -l)
if [ "$COUNT" -eq 1 ]; then
    log_info "✓ User1 sees only their product ($COUNT product)"
    TESTS_RUN=$((TESTS_RUN + 1))
    TESTS_PASSED=$((TESTS_PASSED + 1))
else
    log_error "✗ Expected 1 product, got $COUNT"
    TESTS_RUN=$((TESTS_RUN + 1))
    TESTS_FAILED=$((TESTS_FAILED + 1))
fi

# Test 8: User1 sends invite to User2
log_test "Test 8: User1 sends invite to User2"
STATUS=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE_URL/invite" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN1" \
    -d '{"login":"testuser2"}')
assert_status 200 "$STATUS" "User1 sends invite"

# Test 9: User2 checks incoming invites
log_test "Test 9: User2 checks incoming invites"
RESPONSE=$(curl -s -X GET "$BASE_URL/invite" \
    -H "Authorization: Bearer $TOKEN2")
COUNT=$(echo "$RESPONSE" | grep -o '"from_login":"testuser1"' | wc -l)
if [ "$COUNT" -ge 1 ]; then
    log_info "✓ User2 sees invite from User1"
    TESTS_RUN=$((TESTS_RUN + 1))
    TESTS_PASSED=$((TESTS_PASSED + 1))
else
    log_error "✗ User2 should see invite from User1"
    TESTS_RUN=$((TESTS_RUN + 1))
    TESTS_FAILED=$((TESTS_FAILED + 1))
fi

# Test 10: User2 sends mutual invite to User1 (group created!)
log_test "Test 10: User2 sends mutual invite (creates group)"
RESPONSE=$(curl -s -X POST "$BASE_URL/invite" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN2" \
    -d '{"login":"testuser1"}')
if echo "$RESPONSE" | grep -q "group created\|mutual_invite"; then
    log_info "✓ Group created via mutual invites"
    TESTS_RUN=$((TESTS_RUN + 1))
    TESTS_PASSED=$((TESTS_PASSED + 1))
else
    log_error "✗ Group creation failed"
    log_info "Response: $RESPONSE"
    TESTS_RUN=$((TESTS_RUN + 1))
    TESTS_FAILED=$((TESTS_FAILED + 1))
fi

# Test 11: User1 checks group members
log_test "Test 11: User1 checks group members"
RESPONSE=$(curl -s -X GET "$BASE_URL/group" \
    -H "Authorization: Bearer $TOKEN1")
COUNT=$(echo "$RESPONSE" | grep -o '"user_id"' | wc -l)
if [ "$COUNT" -eq 2 ]; then
    log_info "✓ Group has 2 members"
    TESTS_RUN=$((TESTS_RUN + 1))
    TESTS_PASSED=$((TESTS_PASSED + 1))
else
    log_error "✗ Expected 2 group members, got $COUNT"
    log_info "Response: $RESPONSE"
    TESTS_RUN=$((TESTS_RUN + 1))
    TESTS_FAILED=$((TESTS_FAILED + 1))
fi

# Test 12: User1 gets products (should see both users' products now!)
log_test "Test 12: User1 gets products (with group sharing)"
RESPONSE=$(curl -s -X GET "$BASE_URL/products" \
    -H "Authorization: Bearer $TOKEN1")
COUNT=$(echo "$RESPONSE" | grep -o '"id"' | wc -l)
if [ "$COUNT" -eq 2 ]; then
    log_info "✓ User1 now sees products from group (2 products)"
    TESTS_RUN=$((TESTS_RUN + 1))
    TESTS_PASSED=$((TESTS_PASSED + 1))
else
    log_error "✗ Expected 2 products from group, got $COUNT"
    log_info "Response: $RESPONSE"
    TESTS_RUN=$((TESTS_RUN + 1))
    TESTS_FAILED=$((TESTS_FAILED + 1))
fi

# Test 13: Verify user_id fields are correct
log_test "Test 13: Verify user_id fields in shared products"
RESPONSE=$(curl -s -X GET "$BASE_URL/products" \
    -H "Authorization: Bearer $TOKEN1")
if echo "$RESPONSE" | grep -q '"user_id":' && [ "$(echo "$RESPONSE" | grep -o '"user_id":' | wc -l)" -eq 2 ]; then
    log_info "✓ All products have user_id field"
    TESTS_RUN=$((TESTS_RUN + 1))
    TESTS_PASSED=$((TESTS_PASSED + 1))
else
    log_error "✗ Products missing user_id field"
    TESTS_RUN=$((TESTS_RUN + 1))
    TESTS_FAILED=$((TESTS_FAILED + 1))
fi

# Test 14: User3 (not in group) should not see other users' products
log_test "Test 14: User3 (not in group) isolation check"
RESPONSE=$(curl -s -X GET "$BASE_URL/products" \
    -H "Authorization: Bearer $TOKEN3")
COUNT=$(echo "$RESPONSE" | grep -o '"id"' | wc -l)
if [ "$COUNT" -eq 0 ]; then
    log_info "✓ User3 correctly sees no products (not in group)"
    TESTS_RUN=$((TESTS_RUN + 1))
    TESTS_PASSED=$((TESTS_PASSED + 1))
else
    log_error "✗ User3 should see 0 products, saw $COUNT"
    TESTS_RUN=$((TESTS_RUN + 1))
    TESTS_FAILED=$((TESTS_FAILED + 1))
fi

# Test 15: Test purchase sharing
log_test "Test 15: User1 creates a purchase"
STATUS=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE_URL/purchases" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN1" \
    -d '{"product_id":1,"quantity":2,"price":1000,"date":"2024-01-01T00:00:00Z","store":"Store1","tags":["buy"],"receipt_id":100}')
assert_status 200 "$STATUS" "User1 creates purchase"

# Test 16: User2 creates a purchase
log_test "Test 16: User2 creates a purchase"
STATUS=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE_URL/purchases" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN2" \
    -d '{"product_id":2,"quantity":3,"price":2000,"date":"2024-01-02T00:00:00Z","store":"Store2","tags":["sale"],"receipt_id":200}')
assert_status 200 "$STATUS" "User2 creates purchase"

# Test 17: User1 gets purchases (should see both users' purchases)
log_test "Test 17: User1 gets shared purchases from group"
RESPONSE=$(curl -s -X GET "$BASE_URL/purchases" \
    -H "Authorization: Bearer $TOKEN1")
COUNT=$(echo "$RESPONSE" | grep -o '"id"' | wc -l)
if [ "$COUNT" -eq 2 ]; then
    log_info "✓ User1 sees shared purchases from group (2 purchases)"
    TESTS_RUN=$((TESTS_RUN + 1))
    TESTS_PASSED=$((TESTS_PASSED + 1))
else
    log_error "✗ Expected 2 shared purchases, got $COUNT"
    log_info "Response: $RESPONSE"
    TESTS_RUN=$((TESTS_RUN + 1))
    TESTS_FAILED=$((TESTS_FAILED + 1))
fi

# Test 18: User2 leaves the group
log_test "Test 18: User2 leaves the group"
STATUS=$(curl -s -o /dev/null -w "%{http_code}" -X DELETE "$BASE_URL/group" \
    -H "Authorization: Bearer $TOKEN2")
assert_status 200 "$STATUS" "User2 leaves group"

# Test 19: Verify group was auto-deleted (only 1 member remained)
log_test "Test 19: Verify group auto-deletion"
RESPONSE=$(curl -s -X GET "$BASE_URL/group" \
    -H "Authorization: Bearer $TOKEN1")
if echo "$RESPONSE" | grep -q '"members":\[\]'; then
    log_info "✓ Group was auto-deleted (only 1 member remained)"
    TESTS_RUN=$((TESTS_RUN + 1))
    TESTS_PASSED=$((TESTS_PASSED + 1))
else
    log_error "✗ Group should be auto-deleted"
    log_info "Response: $RESPONSE"
    TESTS_RUN=$((TESTS_RUN + 1))
    TESTS_FAILED=$((TESTS_FAILED + 1))
fi

# Test 20: User1 should only see their own products again
log_test "Test 20: User1 products after group dissolution"
RESPONSE=$(curl -s -X GET "$BASE_URL/products" \
    -H "Authorization: Bearer $TOKEN1")
COUNT=$(echo "$RESPONSE" | grep -o '"id"' | wc -l)
if [ "$COUNT" -eq 1 ]; then
    log_info "✓ User1 back to seeing only their products (1 product)"
    TESTS_RUN=$((TESTS_RUN + 1))
    TESTS_PASSED=$((TESTS_PASSED + 1))
else
    log_error "✗ Expected 1 product after group dissolved, got $COUNT"
    TESTS_RUN=$((TESTS_RUN + 1))
    TESTS_FAILED=$((TESTS_FAILED + 1))
fi

echo ""
echo "========================================="
echo "           Test Summary"
echo "========================================="
echo "Total Tests:  $TESTS_RUN"
echo "Passed:       $TESTS_PASSED"
echo "Failed:       $TESTS_FAILED"
echo "========================================="

if [ "$TESTS_FAILED" -eq 0 ]; then
    echo -e "${GREEN}All tests passed!${NC}"
    exit 0
else
    echo -e "${RED}Some tests failed!${NC}"
    exit 1
fi
