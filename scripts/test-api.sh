#!/bin/bash

# E-Wallet API Test Script
# Tests all API endpoints with HMAC authentication

set -e

API_URL="http://localhost:8080/api/v1"
USER_ID="alif_partner"
SECRET_KEY="alif_secret_2025"

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to compute HMAC-SHA1
compute_hmac() {
    local data="$1"
    echo -n "$data" | openssl dgst -sha1 -hmac "$SECRET_KEY" | cut -d' ' -f2
}

# Function to make API request (expects success)
api_request() {
    local endpoint="$1"
    local data="$2"
    # shellcheck disable=SC2155
    local digest=$(compute_hmac "$data")
    
    echo -e "${YELLOW}Testing: $endpoint${NC}"
    echo "Request: $data"
    
    response=$(curl -s -X POST "$API_URL$endpoint" \
        -H "Content-Type: application/json" \
        -H "X-UserId: $USER_ID" \
        -H "X-Digest: $digest" \
        -d "$data")
    
    echo "Response: $response"
    echo ""
    
    if echo "$response" | grep -q "error"; then
        echo -e "${RED}✗ Test failed${NC}"
        return 1
    else
        echo -e "${GREEN}✓ Test passed${NC}"
        return 0
    fi
}

# Function to make API request (expects error)
api_request_error() {
    local endpoint="$1"
    local data="$2"
    local expected_error="$3"
    # shellcheck disable=SC2155
    local digest=$(compute_hmac "$data")
    
    echo -e "${YELLOW}Testing: $endpoint${NC}"
    echo "Request: $data"
    
    response=$(curl -s -X POST "$API_URL$endpoint" \
        -H "Content-Type: application/json" \
        -H "X-UserId: $USER_ID" \
        -H "X-Digest: $digest" \
        -d "$data")
    
    echo "Response: $response"
    echo ""
    
    if echo "$response" | grep -q "$expected_error"; then
        echo -e "${GREEN}✓ Test passed (error expected)${NC}"
        return 0
    else
        echo -e "${RED}✗ Test failed (expected error: $expected_error)${NC}"
        return 1
    fi
}

echo "========================================="
echo "E-Wallet API Test Suite"
echo "========================================="
echo ""

# Test 1: Check if wallet exists
echo -e "${YELLOW}Test 1: Check Wallet Existence${NC}"
api_request "/wallet/check" '{"account_id":"992900123456"}'
echo "========================================="
echo ""

# Test 2: Check non-existent wallet
echo -e "${YELLOW}Test 2: Check Non-Existent Wallet${NC}"
api_request "/wallet/check" '{"account_id":"999999999999"}'
echo "========================================="
echo ""

# Test 3: Get wallet balance
echo -e "${YELLOW}Test 3: Get Wallet Balance${NC}"
api_request "/wallet/balance" '{"account_id":"992900123456"}'
echo "========================================="
echo ""

# Test 4: Deposit to wallet
echo -e "${YELLOW}Test 4: Deposit to Wallet${NC}"
api_request "/wallet/deposit" '{"account_id":"992900123456","amount":10000}'
echo "========================================="
echo ""

# Test 5: Get monthly stats
echo -e "${YELLOW}Test 5: Get Monthly Statistics${NC}"
api_request "/wallet/monthly-stats" '{"account_id":"992900123456"}'
echo "========================================="
echo ""

# Test 6: Deposit exceeding limit (unidentified)
echo -e "${YELLOW}Test 6: Deposit Exceeding Limit (should fail)${NC}"
api_request_error "/wallet/deposit" '{"account_id":"992901234567","amount":2000000}' "BALANCE_LIMIT_EXCEEDED"
echo "========================================="
echo ""

# Test 7: Invalid amount
echo -e "${YELLOW}Test 7: Invalid Amount (should fail)${NC}"
api_request_error "/wallet/deposit" '{"account_id":"992900123456","amount":-1000}' "VALIDATION_FAILED"
echo "========================================="
echo ""

# Test 8: Missing authentication
echo -e "${YELLOW}Test 8: Missing Authentication (should fail)${NC}"
response=$(curl -s -X POST "$API_URL/wallet/balance" \
    -H "Content-Type: application/json" \
    -d '{"account_id":"992900123456"}')
echo "Response: $response"
if echo "$response" | grep -q "MISSING_AUTH_DATA"; then
    echo -e "${GREEN}✓ Test passed (error expected)${NC}"
else
    echo -e "${RED}✗ Test failed${NC}"
fi
echo "========================================="
echo ""

# Test 9: Invalid HMAC signature
echo -e "${YELLOW}Test 9: Invalid HMAC Signature (should fail)${NC}"
response=$(curl -s -X POST "$API_URL/wallet/balance" \
    -H "Content-Type: application/json" \
    -H "X-UserId: $USER_ID" \
    -H "X-Digest: invalid_signature_here" \
    -d '{"account_id":"992900123456"}')
echo "Response: $response"
if echo "$response" | grep -q "INVALID_SIGNATURE"; then
    echo -e "${GREEN}✓ Test passed (error expected)${NC}"
else
    echo -e "${RED}✗ Test failed${NC}"
fi
echo "========================================="

echo ""
echo -e "${GREEN}All tests completed!${NC}"
