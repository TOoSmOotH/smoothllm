#!/bin/bash

# Test script for media upload functionality

# First, register and login to get a token
echo "Registering test user..."
REGISTER_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "testuser@example.com",
    "password": "testpassword123",
    "username": "testuser"
  }')

echo "Register response: $REGISTER_RESPONSE"

# Login to get token
echo "Logging in..."
LOGIN_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "testuser@example.com",
    "password": "testpassword123"
  }')

echo "Login response: $LOGIN_RESPONSE"

# Extract token
TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"access_token":"[^"]*' | sed 's/"access_token":"//')

if [ -z "$TOKEN" ]; then
  echo "Failed to get token"
  exit 1
fi

echo "Token: $TOKEN"

# Create a simple test image (1x1 pixel PNG)
echo "Creating test image..."
echo "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNkYPhfDwAChwGA60e6kgAAAABJRU5ErkJggg==" | base64 -d > test-avatar.png

# Test avatar upload
echo "Testing avatar upload..."
UPLOAD_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/media/avatar \
  -H "Authorization: Bearer $TOKEN" \
  -F "avatar=@test-avatar.png")

echo "Upload response: $UPLOAD_RESPONSE"

# Extract media ID from response
MEDIA_ID=$(echo $UPLOAD_RESPONSE | grep -o '"media_id":[0-9]*' | sed 's/"media_id"://')

if [ -n "$MEDIA_ID" ]; then
  echo "Successfully uploaded avatar with media ID: $MEDIA_ID"
  
  # Test getting media info
  echo "Testing get media info..."
  curl -s -X GET http://localhost:8080/api/v1/media/$MEDIA_ID \
    -H "Authorization: Bearer $TOKEN"
  
  # Test deleting media
  echo "Testing delete media..."
  DELETE_RESPONSE=$(curl -s -X DELETE http://localhost:8080/api/v1/media/$MEDIA_ID \
    -H "Authorization: Bearer $TOKEN")
  
  echo "Delete response: $DELETE_RESPONSE"
else
  echo "Failed to upload avatar"
fi

# Clean up
rm -f test-avatar.png

echo "Test completed!"